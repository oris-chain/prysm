package light

import (
	"context"
	"sync"
	"time"

	"github.com/prysmaticlabs/prysm/beacon-chain/blockchain"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/feed"
	statefeed "github.com/prysmaticlabs/prysm/beacon-chain/core/feed/state"
	"github.com/prysmaticlabs/prysm/beacon-chain/db/iface"
	syncSrv "github.com/prysmaticlabs/prysm/beacon-chain/sync"
	"github.com/prysmaticlabs/prysm/config/params"
	"github.com/prysmaticlabs/prysm/encoding/bytesutil"
	v1 "github.com/prysmaticlabs/prysm/proto/eth/v1"
	ethpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/time/slots"
)

type Config struct {
	Database            iface.NoHeadAccessDatabase
	HeadFetcher         blockchain.HeadFetcher
	FinalizationFetcher blockchain.FinalizationFetcher
	StateNotifier       statefeed.Notifier
	TimeFetcher         blockchain.TimeFetcher
	SyncChecker         syncSrv.Checker
}

type Service struct {
	cfg          *Config
	cancelFunc   context.CancelFunc
	prevHeadData map[[32]byte]*ethpb.SyncAttestedData
	lock         sync.RWMutex
}

// New --
func New(ctx context.Context, cfg *Config) *Service {
	return &Service{
		cfg:          cfg,
		prevHeadData: make(map[[32]byte]*ethpb.SyncAttestedData),
	}
}

func (s *Service) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancelFunc = cancel
	s.waitForSync(ctx, s.cfg.TimeFetcher.GenesisTime())
	checkpointRoot := bytesutil.ToBytes32(s.cfg.FinalizationFetcher.FinalizedCheckpt().Root)
	block, err := s.cfg.Database.Block(ctx, checkpointRoot)
	if err != nil {
		log.Error(err)
		return
	}
	st, err := s.cfg.Database.State(ctx, checkpointRoot)
	if err != nil {
		log.Error(err)
		return
	}
	s.cfg.FinalizationFetcher.FinalizedCheckpt()
	// Call with finalized checkpoint data.
	if err := s.onFinalized(ctx, st, block.Block()); err != nil {
		log.Fatal(err)
	}
	go s.listenForNewHead(ctx)
}

func (s *Service) Stop() error {
	s.cancelFunc()
	return nil
}

func (s *Service) Status() error {
	return nil
}

func (s *Service) listenForNewHead(ctx context.Context) {
	stateChan := make(chan *feed.Event, 1)
	sub := s.cfg.StateNotifier.StateFeed().Subscribe(stateChan)
	defer sub.Unsubscribe()
	for {
		select {
		case ev := <-stateChan:
			if ev.Type == statefeed.NewHead {
				head, err := s.cfg.HeadFetcher.HeadBlock(ctx)
				if err != nil {
					log.Error(err)
					continue
				}
				st, err := s.cfg.HeadFetcher.HeadState(ctx)
				if err != nil {
					log.Error(err)
					continue
				}
				if err := s.onHead(ctx, st, head.Block()); err != nil {
					log.Error(err)
					continue
				}
			} else if ev.Type == statefeed.FinalizedCheckpoint {
				finalizedCheckpoint, ok := ev.Data.(*v1.EventFinalizedCheckpoint)
				if !ok {
					continue
				}
				checkpointRoot := bytesutil.ToBytes32(finalizedCheckpoint.Block)
				block, err := s.cfg.Database.Block(ctx, checkpointRoot)
				if err != nil {
					log.Error(err)
					continue
				}
				st, err := s.cfg.Database.State(ctx, checkpointRoot)
				if err != nil {
					log.Error(err)
					continue
				}
				if err := s.onFinalized(ctx, st, block.Block()); err != nil {
					log.Error(err)
					continue
				}
			}
		case <-sub.Err():
			return
		case <-ctx.Done():
			return
		}
	}
}

func (s *Service) waitForSync(ctx context.Context, genesisTime time.Time) {
	if slots.SinceGenesis(genesisTime) == 0 || !s.cfg.SyncChecker.Syncing() {
		return
	}
	slotTicker := slots.NewSlotTicker(genesisTime, params.BeaconConfig().SecondsPerSlot)
	defer slotTicker.Done()
	for {
		select {
		case <-slotTicker.C():
			// If node is still syncing, do not operate.
			if s.cfg.SyncChecker.Syncing() {
				continue
			}
			return
		case <-ctx.Done():
			return
		}
	}
}
