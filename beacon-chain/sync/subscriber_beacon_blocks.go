package sync

import (
	"context"
	"fmt"
	"time"

	"github.com/prysmaticlabs/prysm/v4/beacon-chain/blockchain"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/p2p/types"
	"github.com/prysmaticlabs/prysm/v4/config/params"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/blocks"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/interfaces"
	"github.com/prysmaticlabs/prysm/v4/encoding/bytesutil"
	eth "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v4/runtime/version"
	"github.com/prysmaticlabs/prysm/v4/time/slots"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

func (s *Service) beaconBlockSubscriber(ctx context.Context, msg proto.Message) error {
	signed, err := blocks.NewSignedBeaconBlock(msg)
	if err != nil {
		return err
	}
	if err := blocks.BeaconBlockIsNil(signed); err != nil {
		return err
	}

	s.setSeenBlockIndexSlot(signed.Block().Slot(), signed.Block().ProposerIndex())

	block := signed.Block()

	root, err := block.HashTreeRoot()
	if err != nil {
		return err
	}

	if block.Version() >= version.Deneb {
		if err := s.blockAndBlobs.addBlock(signed); err != nil {
			return err
		}
		return s.importBlockAndBlobs(ctx, root)
	}

	return s.receiveBlock(ctx, signed, root)
}

func (s *Service) importBlockAndBlobs(ctx context.Context, root [32]byte) error {
	canImport, err := s.blockAndBlobs.canImport(root)
	if err != nil {
		return err
	}
	if !canImport {
		return nil
	}
	return s.receiveBlockAndBlobs(ctx, root)
}

func (s *Service) receiveBlockAndBlobs(ctx context.Context, root [32]byte) error {
	signed, err := s.blockAndBlobs.getBlock(root)
	if err != nil {
		return err
	}

	if err = s.receiveBlock(ctx, signed, root); err != nil {
		return err
	}
	kzgs, err := signed.Block().Body().BlobKzgCommitments()
	if err != nil {
		return err
	}

	scs := make([]*eth.BlobSidecar, len(kzgs))
	for i := 0; i < len(kzgs); i++ {
		index := uint64(i)
		scs[i], err = s.blockAndBlobs.getBlob(root, index)
		if err != nil {
			return err
		}
	}

	if len(scs) > 0 {
		log.WithFields(logrus.Fields{
			"slot":      scs[0].Slot,
			"root":      fmt.Sprintf("%#x", scs[0].BlockRoot),
			"blobCount": len(scs),
		}).Info("Saving blobs")
		if err := s.cfg.beaconDB.SaveBlobSidecar(ctx, scs); err != nil {
			return err
		}
	}

	s.blockAndBlobs.delete(root)

	return nil
}

func (s *Service) receiveBlock(ctx context.Context, signed interfaces.ReadOnlySignedBeaconBlock, root [32]byte) error {
	if err := s.cfg.chain.ReceiveBlock(ctx, signed, root); err != nil {
		if blockchain.IsInvalidBlock(err) {
			r := blockchain.InvalidBlockRoot(err)
			if r != [32]byte{} {
				s.setBadBlock(ctx, r) // Setting head block as bad.
			} else {
				s.setBadBlock(ctx, root)
			}
		}
		// Set the returned invalid ancestors as bad.
		for _, root := range blockchain.InvalidAncestorRoots(err) {
			s.setBadBlock(ctx, root)
		}
		return err
	}
	return nil
}

func (s *Service) requestMissingBlobsRoutine(ctx context.Context) {

	go func() {
		ticker := slots.NewSlotTickerWithOffset(s.cfg.chain.GenesisTime(), time.Second, params.BeaconConfig().SecondsPerSlot)
		for {
			select {
			case <-ticker.C():
				log.Info("Tick ", ticker.C())
				cp := s.cfg.chain.FinalizedCheckpt()
				_, bestPeers := s.cfg.p2p.Peers().BestFinalized(maxPeerRequest, cp.Epoch)
				if len(bestPeers) == 0 {
					log.Warn("No peers to request missing blobs")
					continue
				}
				var reqs []*eth.BlobIdentifier

				m, err := s.blockAndBlobs.missingRootAndIndex(ctx)
				if err != nil {
					log.WithError(err).Error("Failed to get missing root and index")
					continue
				}
				if len(m) == 0 {
					continue
				}
				for r, indices := range m {
					for _, i := range indices {
						reqs = append(reqs, &eth.BlobIdentifier{
							BlockRoot: r[:],
							Index:     i,
						})
					}
				}
				req := types.BlobSidecarsByRootReq(reqs)
				scs, err := SendBlobSidecarByRoot(ctx, s.cfg.chain, s.cfg.p2p, bestPeers[0], &req)
				if err != nil {
					log.WithError(err).Error("Failed to send blob sidecar by root")
					continue
				}
				for _, sc := range scs {
					if err := s.blockAndBlobs.addBlob(sc); err != nil {
						log.WithError(err).Error("Failed to add blob")
						continue
					}
					if err := s.importBlockAndBlobs(ctx, bytesutil.ToBytes32(sc.BlockRoot)); err != nil {
						log.WithError(err).Error("Failed to import block and blobs")
						continue
					}
				}

			case <-ctx.Done():
				log.Debug("Context closed, exiting routine")
				return
			}
		}
	}()
}
