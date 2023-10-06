package sync

import (
	"context"

	libp2pcore "github.com/libp2p/go-libp2p/core"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/db"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/execution"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/p2p/types"
	"github.com/prysmaticlabs/prysm/v4/config/params"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/blocks"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/interfaces"
	"github.com/prysmaticlabs/prysm/v4/encoding/bytesutil"
	eth "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v4/runtime/version"
)

// sendRecentBeaconBlocksRequest sends a recent beacon blocks request to a peer to get
// those corresponding blocks from that peer.
func (s *Service) sendRecentBeaconBlocksRequest(ctx context.Context, blockRoots *types.BeaconBlockByRootsReq, id peer.ID) error {
	ctx, cancel := context.WithTimeout(ctx, respTimeout)
	defer cancel()

	blks, err := SendBeaconBlocksByRootRequest(ctx, s.cfg.clock, s.cfg.p2p, id, blockRoots, func(blk interfaces.ReadOnlySignedBeaconBlock) error {
		blkRoot, err := blk.Block().HashTreeRoot()
		if err != nil {
			return err
		}
		s.pendingQueueLock.Lock()
		defer s.pendingQueueLock.Unlock()
		log.Errorf("Inserting block to pending queue: %s %d", blkRoot, blk.Block().Slot())
		if err := s.insertBlockToPendingQueue(blk.Block().Slot(), blk, blkRoot); err != nil {
			return err
		}
		return nil
	})

	for _, blk := range blks {
		// Skip blocks before deneb because they have no blob.
		if blk.Version() < version.Deneb {
			continue
		}
		blkRoot, err := blk.Block().HashTreeRoot()
		if err != nil {
			return err
		}
		if err := s.requestPendingBlobs(ctx, blk.Block(), blkRoot[:], id); err != nil {
			return err
		}
	}

	return err
}

// beaconBlocksRootRPCHandler looks up the request blocks from the database from the given block roots.
func (s *Service) beaconBlocksRootRPCHandler(ctx context.Context, msg interface{}, stream libp2pcore.Stream) error {
	ctx, cancel := context.WithTimeout(ctx, ttfbTimeout)
	defer cancel()
	SetRPCStreamDeadlines(stream)
	log := log.WithField("handler", "beacon_blocks_by_root")

	rawMsg, ok := msg.(*types.BeaconBlockByRootsReq)
	if !ok {
		return errors.New("message is not type BeaconBlockByRootsReq")
	}
	blockRoots := *rawMsg
	if err := s.rateLimiter.validateRequest(stream, uint64(len(blockRoots))); err != nil {
		return err
	}
	if len(blockRoots) == 0 {
		// Add to rate limiter in the event no
		// roots are requested.
		s.rateLimiter.add(stream, 1)
		s.writeErrorResponseToStream(responseCodeInvalidRequest, "no block roots provided in request", stream)
		return errors.New("no block roots provided")
	}

	if uint64(len(blockRoots)) > params.BeaconNetworkConfig().MaxRequestBlocks {
		s.cfg.p2p.Peers().Scorers().BadResponsesScorer().Increment(stream.Conn().RemotePeer())
		s.writeErrorResponseToStream(responseCodeInvalidRequest, "requested more than the max block limit", stream)
		return errors.New("requested more than the max block limit")
	}
	s.rateLimiter.add(stream, int64(len(blockRoots)))

	for _, root := range blockRoots {
		blk, err := s.cfg.beaconDB.Block(ctx, root)
		if err != nil {
			log.WithError(err).Debug("Could not fetch block")
			s.writeErrorResponseToStream(responseCodeServerError, types.ErrGeneric.Error(), stream)
			return err
		}
		if err := blocks.BeaconBlockIsNil(blk); err != nil {
			continue
		}

		if blk.Block().IsBlinded() {
			blk, err = s.cfg.executionPayloadReconstructor.ReconstructFullBlock(ctx, blk)
			if err != nil {
				if errors.Is(err, execution.EmptyBlockHash) {
					log.WithError(err).Warn("Could not reconstruct block from header with syncing execution client. Waiting to complete syncing")
				} else {
					log.WithError(err).Error("Could not get reconstruct full block from blinded body")
				}
				s.writeErrorResponseToStream(responseCodeServerError, types.ErrGeneric.Error(), stream)
				return err
			}
		}

		if err := s.chunkBlockWriter(stream, blk); err != nil {
			return err
		}
	}

	closeStream(stream, log)
	return nil
}

func (s *Service) requestPendingBlobs(ctx context.Context, b interfaces.ReadOnlyBeaconBlock, br []byte, id peer.ID) error {
	// Block before deneb has no blob.
	if b.Version() < version.Deneb {
		return nil
	}
	c, err := b.Body().BlobKzgCommitments()
	if err != nil {
		return err
	}
	// No op if the block has no blob commitments.
	if len(c) == 0 {
		return nil
	}

	savedBlobs, err := s.cfg.beaconDB.BlobSidecarsByRoot(ctx, bytesutil.ToBytes32(br))
	switch {
	case errors.Is(err, db.ErrNotFound):
	case err != nil:
		return err
	}
	savedIndices := make(map[uint64]struct{})
	for _, blob := range savedBlobs {
		savedIndices[blob.Index] = struct{}{}
	}

	// Build request for blob sidecars.
	blobId := make([]*eth.BlobIdentifier, 0, len(c))
	for i := range c {
		if _, ok := savedIndices[uint64(i)]; ok {
			continue
		}
		blobId = append(blobId, &eth.BlobIdentifier{Index: uint64(i), BlockRoot: br})
	}
	if len(blobId) == 0 {
		return nil
	}

	ctxByte, err := ContextByteVersionsForValRoot(s.cfg.chain.GenesisValidatorsRoot())
	if err != nil {
		return err
	}
	req := types.BlobSidecarsByRootReq(blobId)

	// Send request to a random peer.
	blobSidecars, err := SendBlobSidecarByRoot(ctx, s.cfg.clock, s.cfg.p2p, id, ctxByte, &req)
	if err != nil {
		return err
	}

	for _, sidecar := range blobSidecars {
		log.WithFields(blobFields(sidecar)).Debug("Received blob sidecar gossip RPC")
	}

	return s.cfg.beaconDB.SaveBlobSidecar(ctx, blobSidecars)
}
