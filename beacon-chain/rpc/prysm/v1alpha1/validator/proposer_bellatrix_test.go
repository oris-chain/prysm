package validator

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	blockchainTest "github.com/prysmaticlabs/prysm/v3/beacon-chain/blockchain/testing"
	builderTest "github.com/prysmaticlabs/prysm/v3/beacon-chain/builder/testing"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/cache"
	consensusblocks "github.com/prysmaticlabs/prysm/v3/beacon-chain/core/blocks"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/core/helpers"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/core/signing"
	prysmtime "github.com/prysmaticlabs/prysm/v3/beacon-chain/core/time"
	dbTest "github.com/prysmaticlabs/prysm/v3/beacon-chain/db/testing"
	mockExecution "github.com/prysmaticlabs/prysm/v3/beacon-chain/execution/testing"
	doublylinkedtree "github.com/prysmaticlabs/prysm/v3/beacon-chain/forkchoice/doubly-linked-tree"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/operations/attestations"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/operations/slashings"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/operations/synccommittee"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/operations/voluntaryexits"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/state/stategen"
	mockSync "github.com/prysmaticlabs/prysm/v3/beacon-chain/sync/initial-sync/testing"
	fieldparams "github.com/prysmaticlabs/prysm/v3/config/fieldparams"
	"github.com/prysmaticlabs/prysm/v3/config/params"
	"github.com/prysmaticlabs/prysm/v3/consensus-types/blocks"
	"github.com/prysmaticlabs/prysm/v3/consensus-types/interfaces"
	types "github.com/prysmaticlabs/prysm/v3/consensus-types/primitives"
	"github.com/prysmaticlabs/prysm/v3/crypto/bls"
	"github.com/prysmaticlabs/prysm/v3/encoding/bytesutil"
	"github.com/prysmaticlabs/prysm/v3/encoding/ssz"
	v1 "github.com/prysmaticlabs/prysm/v3/proto/engine/v1"
	ethpb "github.com/prysmaticlabs/prysm/v3/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v3/testing/require"
	"github.com/prysmaticlabs/prysm/v3/testing/util"
	"github.com/prysmaticlabs/prysm/v3/time/slots"
	logTest "github.com/sirupsen/logrus/hooks/test"
)

func TestServer_getPayloadHeader(t *testing.T) {
	emptyRoot, err := ssz.TransactionsRoot([][]byte{})
	require.NoError(t, err)
	ti, err := slots.ToTime(uint64(time.Now().Unix()), 0)
	require.NoError(t, err)

	sk, err := bls.RandKey()
	require.NoError(t, err)
	bid := &ethpb.BuilderBid{
		Header: &v1.ExecutionPayloadHeader{
			FeeRecipient:     make([]byte, fieldparams.FeeRecipientLength),
			StateRoot:        make([]byte, fieldparams.RootLength),
			ReceiptsRoot:     make([]byte, fieldparams.RootLength),
			LogsBloom:        make([]byte, fieldparams.LogsBloomLength),
			PrevRandao:       make([]byte, fieldparams.RootLength),
			BaseFeePerGas:    make([]byte, fieldparams.RootLength),
			BlockHash:        make([]byte, fieldparams.RootLength),
			TransactionsRoot: bytesutil.PadTo([]byte{1}, fieldparams.RootLength),
			ParentHash:       params.BeaconConfig().ZeroHash[:],
			Timestamp:        uint64(ti.Unix()),
		},
		Pubkey: sk.PublicKey().Marshal(),
		Value:  bytesutil.PadTo([]byte{1, 2, 3}, 32),
	}
	d := params.BeaconConfig().DomainApplicationBuilder
	domain, err := signing.ComputeDomain(d, nil, nil)
	require.NoError(t, err)
	sr, err := signing.ComputeSigningRoot(bid, domain)
	require.NoError(t, err)
	sBid := &ethpb.SignedBuilderBid{
		Message:   bid,
		Signature: sk.Sign(sr[:]).Marshal(),
	}

	require.NoError(t, err)
	tests := []struct {
		name           string
		head           interfaces.SignedBeaconBlock
		mock           *builderTest.MockBuilderService
		fetcher        *blockchainTest.ChainService
		err            string
		returnedHeader *v1.ExecutionPayloadHeader
	}{
		{
			name: "head is not bellatrix ready",
			mock: &builderTest.MockBuilderService{},
			fetcher: &blockchainTest.ChainService{
				Block: func() interfaces.SignedBeaconBlock {
					wb, err := blocks.NewSignedBeaconBlock(util.NewBeaconBlock())
					require.NoError(t, err)
					return wb
				}(),
			},
		},
		{
			name: "get header failed",
			mock: &builderTest.MockBuilderService{
				ErrGetHeader: errors.New("can't get header"),
			},
			fetcher: &blockchainTest.ChainService{
				Block: func() interfaces.SignedBeaconBlock {
					wb, err := blocks.NewSignedBeaconBlock(util.NewBeaconBlockBellatrix())
					require.NoError(t, err)
					return wb
				}(),
			},
			err: "can't get header",
		},
		{
			name: "0 bid",
			mock: &builderTest.MockBuilderService{
				Bid: &ethpb.SignedBuilderBid{
					Message: &ethpb.BuilderBid{
						Header: &v1.ExecutionPayloadHeader{
							BlockNumber: 123,
						},
					},
				},
			},
			fetcher: &blockchainTest.ChainService{
				Block: func() interfaces.SignedBeaconBlock {
					wb, err := blocks.NewSignedBeaconBlock(util.NewBeaconBlockBellatrix())
					require.NoError(t, err)
					return wb
				}(),
			},
			err: "builder returned header with 0 bid amount",
		},
		{
			name: "invalid tx root",
			mock: &builderTest.MockBuilderService{
				Bid: &ethpb.SignedBuilderBid{
					Message: &ethpb.BuilderBid{
						Value: []byte{1},
						Header: &v1.ExecutionPayloadHeader{
							BlockNumber:      123,
							TransactionsRoot: emptyRoot[:],
						},
					},
				},
			},
			fetcher: &blockchainTest.ChainService{
				Block: func() interfaces.SignedBeaconBlock {
					wb, err := blocks.NewSignedBeaconBlock(util.NewBeaconBlockBellatrix())
					require.NoError(t, err)
					return wb
				}(),
			},
			err: "builder returned header with an empty tx root",
		},
		{
			name: "can get header",
			mock: &builderTest.MockBuilderService{
				Bid: sBid,
			},
			fetcher: &blockchainTest.ChainService{
				Block: func() interfaces.SignedBeaconBlock {
					wb, err := blocks.NewSignedBeaconBlock(util.NewBeaconBlockBellatrix())
					require.NoError(t, err)
					return wb
				}(),
			},
			returnedHeader: bid.Header,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			vs := &Server{BlockBuilder: tc.mock, HeadFetcher: tc.fetcher, TimeFetcher: &blockchainTest.ChainService{
				Genesis: time.Now(),
			}}
			h, err := vs.getPayloadHeaderFromBuilder(context.Background(), 0, 0)
			if tc.err != "" {
				require.ErrorContains(t, tc.err, err)
			} else {
				require.NoError(t, err)
				require.DeepEqual(t, tc.returnedHeader, h)
			}
		})
	}
}

func TestServer_getBuilderBlock(t *testing.T) {
	p := emptyPayload()
	p.GasLimit = 123

	tests := []struct {
		name        string
		blk         interfaces.SignedBeaconBlock
		mock        *builderTest.MockBuilderService
		err         string
		returnedBlk interfaces.SignedBeaconBlock
	}{
		{
			name: "nil block",
			blk:  nil,
			err:  "signed beacon block can't be nil",
		},
		{
			name: "old block version",
			blk: func() interfaces.SignedBeaconBlock {
				wb, err := blocks.NewSignedBeaconBlock(util.NewBeaconBlock())
				require.NoError(t, err)
				return wb
			}(),
			returnedBlk: func() interfaces.SignedBeaconBlock {
				wb, err := blocks.NewSignedBeaconBlock(util.NewBeaconBlock())
				require.NoError(t, err)
				return wb
			}(),
		},
		{
			name: "not configured",
			blk: func() interfaces.SignedBeaconBlock {
				wb, err := blocks.NewSignedBeaconBlock(util.NewBlindedBeaconBlockBellatrix())
				require.NoError(t, err)
				return wb
			}(),
			mock: &builderTest.MockBuilderService{
				HasConfigured: false,
			},
			returnedBlk: func() interfaces.SignedBeaconBlock {
				wb, err := blocks.NewSignedBeaconBlock(util.NewBlindedBeaconBlockBellatrix())
				require.NoError(t, err)
				return wb
			}(),
		},
		{
			name: "submit blind block error",
			blk: func() interfaces.SignedBeaconBlock {
				b := util.NewBlindedBeaconBlockBellatrix()
				b.Block.Slot = 1
				b.Block.ProposerIndex = 2
				wb, err := blocks.NewSignedBeaconBlock(b)
				require.NoError(t, err)
				return wb
			}(),
			mock: &builderTest.MockBuilderService{
				HasConfigured:         true,
				ErrSubmitBlindedBlock: errors.New("can't submit"),
			},
			err: "can't submit",
		},
		{
			name: "head and payload root mismatch",
			blk: func() interfaces.SignedBeaconBlock {
				b := util.NewBlindedBeaconBlockBellatrix()
				b.Block.Slot = 1
				b.Block.ProposerIndex = 2
				wb, err := blocks.NewSignedBeaconBlock(b)
				require.NoError(t, err)
				return wb
			}(),
			mock: &builderTest.MockBuilderService{
				HasConfigured: true,
				Payload:       p,
			},
			returnedBlk: func() interfaces.SignedBeaconBlock {
				b := util.NewBeaconBlockBellatrix()
				b.Block.Slot = 1
				b.Block.ProposerIndex = 2
				b.Block.Body.ExecutionPayload = p
				wb, err := blocks.NewSignedBeaconBlock(b)
				require.NoError(t, err)
				return wb
			}(),
			err: "header and payload root do not match",
		},
		{
			name: "can get payload",
			blk: func() interfaces.SignedBeaconBlock {
				b := util.NewBlindedBeaconBlockBellatrix()
				b.Block.Slot = 1
				b.Block.ProposerIndex = 2
				txRoot, err := ssz.TransactionsRoot([][]byte{})
				require.NoError(t, err)
				b.Block.Body.ExecutionPayloadHeader = &v1.ExecutionPayloadHeader{
					ParentHash:       make([]byte, fieldparams.RootLength),
					FeeRecipient:     make([]byte, fieldparams.FeeRecipientLength),
					StateRoot:        make([]byte, fieldparams.RootLength),
					ReceiptsRoot:     make([]byte, fieldparams.RootLength),
					LogsBloom:        make([]byte, fieldparams.LogsBloomLength),
					PrevRandao:       make([]byte, fieldparams.RootLength),
					BaseFeePerGas:    make([]byte, fieldparams.RootLength),
					BlockHash:        make([]byte, fieldparams.RootLength),
					TransactionsRoot: txRoot[:],
					GasLimit:         123,
				}
				wb, err := blocks.NewSignedBeaconBlock(b)
				require.NoError(t, err)
				return wb
			}(),
			mock: &builderTest.MockBuilderService{
				HasConfigured: true,
				Payload:       p,
			},
			returnedBlk: func() interfaces.SignedBeaconBlock {
				b := util.NewBeaconBlockBellatrix()
				b.Block.Slot = 1
				b.Block.ProposerIndex = 2
				b.Block.Body.ExecutionPayload = p
				wb, err := blocks.NewSignedBeaconBlock(b)
				require.NoError(t, err)
				return wb
			}(),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			vs := &Server{BlockBuilder: tc.mock}
			gotBlk, err := vs.unblindBuilderBlock(context.Background(), tc.blk)
			if tc.err != "" {
				require.ErrorContains(t, tc.err, err)
			} else {
				require.NoError(t, err)
				require.DeepEqual(t, tc.returnedBlk, gotBlk)
			}
		})
	}
}

func TestServer_readyForBuilder(t *testing.T) {
	ctx := context.Background()
	vs := &Server{BeaconDB: dbTest.SetupDB(t)}
	cs := &blockchainTest.ChainService{FinalizedCheckPoint: &ethpb.Checkpoint{}} // Checkpoint root is zeros.
	vs.FinalizationFetcher = cs
	ready, err := vs.readyForBuilder(ctx)
	require.NoError(t, err)
	require.Equal(t, false, ready)

	b := util.NewBeaconBlockBellatrix()
	wb, err := blocks.NewSignedBeaconBlock(b)
	require.NoError(t, err)
	wbr, err := wb.Block().HashTreeRoot()
	require.NoError(t, err)

	b1 := util.NewBeaconBlockBellatrix()
	b1.Block.Body.ExecutionPayload.BlockNumber = 1 // Execution enabled.
	wb1, err := blocks.NewSignedBeaconBlock(b1)
	require.NoError(t, err)
	wbr1, err := wb1.Block().HashTreeRoot()
	require.NoError(t, err)
	require.NoError(t, vs.BeaconDB.SaveBlock(ctx, wb))
	require.NoError(t, vs.BeaconDB.SaveBlock(ctx, wb1))

	// Ready is false given finalized block does not have execution.
	cs = &blockchainTest.ChainService{FinalizedCheckPoint: &ethpb.Checkpoint{Root: wbr[:]}}
	vs.FinalizationFetcher = cs
	ready, err = vs.readyForBuilder(ctx)
	require.NoError(t, err)
	require.Equal(t, false, ready)

	// Ready is true given finalized block has execution.
	cs = &blockchainTest.ChainService{FinalizedCheckPoint: &ethpb.Checkpoint{Root: wbr1[:]}}
	vs.FinalizationFetcher = cs
	ready, err = vs.readyForBuilder(ctx)
	require.NoError(t, err)
	require.Equal(t, true, ready)
}

func TestServer_GetBellatrixBeaconBlock_HappyCase(t *testing.T) {
	db := dbTest.SetupDB(t)
	ctx := context.Background()
	hook := logTest.NewGlobal()

	terminalBlockHash := bytesutil.PadTo([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, 32)
	params.SetupTestConfigCleanup(t)
	cfg := params.BeaconConfig().Copy()
	cfg.BellatrixForkEpoch = 2
	cfg.AltairForkEpoch = 1
	cfg.TerminalBlockHash = common.BytesToHash(terminalBlockHash)
	cfg.TerminalBlockHashActivationEpoch = 2
	params.OverrideBeaconConfig(cfg)

	beaconState, privKeys := util.DeterministicGenesisState(t, 64)
	stateRoot, err := beaconState.HashTreeRoot(ctx)
	require.NoError(t, err, "Could not hash genesis state")

	genesis := consensusblocks.NewGenesisBlock(stateRoot[:])
	wsb, err := blocks.NewSignedBeaconBlock(genesis)
	require.NoError(t, err)
	require.NoError(t, db.SaveBlock(ctx, wsb), "Could not save genesis block")

	parentRoot, err := genesis.Block.HashTreeRoot()
	require.NoError(t, err, "Could not get signing root")
	require.NoError(t, db.SaveState(ctx, beaconState, parentRoot), "Could not save genesis state")
	require.NoError(t, db.SaveHeadBlockRoot(ctx, parentRoot), "Could not save genesis state")

	bellatrixSlot, err := slots.EpochStart(params.BeaconConfig().BellatrixForkEpoch)
	require.NoError(t, err)

	emptyPayload := &v1.ExecutionPayload{
		ParentHash:    make([]byte, fieldparams.RootLength),
		FeeRecipient:  make([]byte, fieldparams.FeeRecipientLength),
		StateRoot:     make([]byte, fieldparams.RootLength),
		ReceiptsRoot:  make([]byte, fieldparams.RootLength),
		LogsBloom:     make([]byte, fieldparams.LogsBloomLength),
		PrevRandao:    make([]byte, fieldparams.RootLength),
		BaseFeePerGas: make([]byte, fieldparams.RootLength),
		BlockHash:     make([]byte, fieldparams.RootLength),
	}
	var scBits [fieldparams.SyncAggregateSyncCommitteeBytesLength]byte
	blk := &ethpb.SignedBeaconBlockBellatrix{
		Block: &ethpb.BeaconBlockBellatrix{
			Slot:       bellatrixSlot + 1,
			ParentRoot: parentRoot[:],
			StateRoot:  genesis.Block.StateRoot,
			Body: &ethpb.BeaconBlockBodyBellatrix{
				RandaoReveal:     genesis.Block.Body.RandaoReveal,
				Graffiti:         genesis.Block.Body.Graffiti,
				Eth1Data:         genesis.Block.Body.Eth1Data,
				SyncAggregate:    &ethpb.SyncAggregate{SyncCommitteeBits: scBits[:], SyncCommitteeSignature: make([]byte, 96)},
				ExecutionPayload: emptyPayload,
			},
		},
		Signature: genesis.Signature,
	}

	blkRoot, err := blk.Block.HashTreeRoot()
	require.NoError(t, err, "Could not get signing root")
	require.NoError(t, db.SaveState(ctx, beaconState, blkRoot), "Could not save genesis state")
	require.NoError(t, db.SaveHeadBlockRoot(ctx, blkRoot), "Could not save genesis state")

	proposerServer := &Server{
		HeadFetcher:       &blockchainTest.ChainService{State: beaconState, Root: parentRoot[:], Optimistic: false},
		TimeFetcher:       &blockchainTest.ChainService{Genesis: time.Now()},
		SyncChecker:       &mockSync.Sync{IsSyncing: false},
		BlockReceiver:     &blockchainTest.ChainService{},
		HeadUpdater:       &blockchainTest.ChainService{},
		ChainStartFetcher: &mockExecution.Chain{},
		Eth1InfoFetcher:   &mockExecution.Chain{},
		MockEth1Votes:     true,
		AttPool:           attestations.NewPool(),
		SlashingsPool:     slashings.NewPool(),
		ExitPool:          voluntaryexits.NewPool(),
		StateGen:          stategen.New(db, doublylinkedtree.New()),
		SyncCommitteePool: synccommittee.NewStore(),
		ExecutionEngineCaller: &mockExecution.EngineClient{
			PayloadIDBytes:   &v1.PayloadIDBytes{1},
			ExecutionPayload: emptyPayload,
		},
		BeaconDB:               db,
		ProposerSlotIndexCache: cache.NewProposerPayloadIDsCache(),
		BlockBuilder:           &builderTest.MockBuilderService{},
	}
	proposerServer.ProposerSlotIndexCache.SetProposerAndPayloadIDs(17, 11, [8]byte{'a'}, parentRoot)

	randaoReveal, err := util.RandaoReveal(beaconState, 0, privKeys)
	require.NoError(t, err)

	block, err := proposerServer.getBellatrixBeaconBlock(ctx, &ethpb.BlockRequest{
		Slot:         bellatrixSlot + 1,
		RandaoReveal: randaoReveal,
	})
	require.NoError(t, err)
	bellatrixBlk, ok := block.GetBlock().(*ethpb.GenericBeaconBlock_Bellatrix)
	require.Equal(t, true, ok)
	require.LogsContain(t, hook, "Computed state root")
	require.DeepEqual(t, emptyPayload, bellatrixBlk.Bellatrix.Body.ExecutionPayload) // Payload should equal.
}

func TestServer_GetBellatrixBeaconBlock_LocalProgressingWithBuilderSkipped(t *testing.T) {
	db := dbTest.SetupDB(t)
	ctx := context.Background()
	hook := logTest.NewGlobal()

	terminalBlockHash := bytesutil.PadTo([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, 32)
	params.SetupTestConfigCleanup(t)
	cfg := params.BeaconConfig().Copy()
	cfg.BellatrixForkEpoch = 2
	cfg.AltairForkEpoch = 1
	cfg.TerminalBlockHash = common.BytesToHash(terminalBlockHash)
	cfg.TerminalBlockHashActivationEpoch = 2
	params.OverrideBeaconConfig(cfg)

	beaconState, privKeys := util.DeterministicGenesisState(t, 64)
	stateRoot, err := beaconState.HashTreeRoot(ctx)
	require.NoError(t, err, "Could not hash genesis state")

	genesis := consensusblocks.NewGenesisBlock(stateRoot[:])
	wsb, err := blocks.NewSignedBeaconBlock(genesis)
	require.NoError(t, err)
	require.NoError(t, db.SaveBlock(ctx, wsb), "Could not save genesis block")

	parentRoot, err := genesis.Block.HashTreeRoot()
	require.NoError(t, err, "Could not get signing root")
	require.NoError(t, db.SaveState(ctx, beaconState, parentRoot), "Could not save genesis state")
	require.NoError(t, db.SaveHeadBlockRoot(ctx, parentRoot), "Could not save genesis state")

	bellatrixSlot, err := slots.EpochStart(params.BeaconConfig().BellatrixForkEpoch)
	require.NoError(t, err)

	emptyPayload := &v1.ExecutionPayload{
		ParentHash:    make([]byte, fieldparams.RootLength),
		FeeRecipient:  make([]byte, fieldparams.FeeRecipientLength),
		StateRoot:     make([]byte, fieldparams.RootLength),
		ReceiptsRoot:  make([]byte, fieldparams.RootLength),
		LogsBloom:     make([]byte, fieldparams.LogsBloomLength),
		PrevRandao:    make([]byte, fieldparams.RootLength),
		BaseFeePerGas: make([]byte, fieldparams.RootLength),
		BlockHash:     make([]byte, fieldparams.RootLength),
	}
	var scBits [fieldparams.SyncAggregateSyncCommitteeBytesLength]byte
	blk := &ethpb.SignedBeaconBlockBellatrix{
		Block: &ethpb.BeaconBlockBellatrix{
			Slot:       bellatrixSlot + 1,
			ParentRoot: parentRoot[:],
			StateRoot:  genesis.Block.StateRoot,
			Body: &ethpb.BeaconBlockBodyBellatrix{
				RandaoReveal:     genesis.Block.Body.RandaoReveal,
				Graffiti:         genesis.Block.Body.Graffiti,
				Eth1Data:         genesis.Block.Body.Eth1Data,
				SyncAggregate:    &ethpb.SyncAggregate{SyncCommitteeBits: scBits[:], SyncCommitteeSignature: make([]byte, 96)},
				ExecutionPayload: emptyPayload,
			},
		},
		Signature: genesis.Signature,
	}

	blkRoot, err := blk.Block.HashTreeRoot()
	require.NoError(t, err, "Could not get signing root")
	require.NoError(t, db.SaveState(ctx, beaconState, blkRoot), "Could not save genesis state")
	require.NoError(t, db.SaveHeadBlockRoot(ctx, blkRoot), "Could not save genesis state")

	proposerServer := &Server{
		HeadFetcher:       &blockchainTest.ChainService{State: beaconState, Root: parentRoot[:], Optimistic: false},
		TimeFetcher:       &blockchainTest.ChainService{Genesis: time.Now()},
		SyncChecker:       &mockSync.Sync{IsSyncing: false},
		BlockReceiver:     &blockchainTest.ChainService{},
		HeadUpdater:       &blockchainTest.ChainService{},
		ChainStartFetcher: &mockExecution.Chain{},
		Eth1InfoFetcher:   &mockExecution.Chain{},
		MockEth1Votes:     true,
		AttPool:           attestations.NewPool(),
		SlashingsPool:     slashings.NewPool(),
		ExitPool:          voluntaryexits.NewPool(),
		StateGen:          stategen.New(db, doublylinkedtree.New()),
		SyncCommitteePool: synccommittee.NewStore(),
		ExecutionEngineCaller: &mockExecution.EngineClient{
			PayloadIDBytes:   &v1.PayloadIDBytes{1},
			ExecutionPayload: emptyPayload,
		},
		BeaconDB:               db,
		ProposerSlotIndexCache: cache.NewProposerPayloadIDsCache(),
		BlockBuilder:           &builderTest.MockBuilderService{},
	}
	proposerServer.ProposerSlotIndexCache.SetProposerAndPayloadIDs(17, 11, [8]byte{'a'}, parentRoot)

	randaoReveal, err := util.RandaoReveal(beaconState, 0, privKeys)
	require.NoError(t, err)

	// Configure builder, this should fail if it's not local engine processing
	proposerServer.BlockBuilder = &builderTest.MockBuilderService{HasConfigured: true, ErrGetHeader: errors.New("bad)")}
	block, err := proposerServer.getBellatrixBeaconBlock(ctx, &ethpb.BlockRequest{
		Slot:         bellatrixSlot + 1,
		RandaoReveal: randaoReveal,
		SkipMevBoost: true,
	})
	require.NoError(t, err)
	bellatrixBlk, ok := block.GetBlock().(*ethpb.GenericBeaconBlock_Bellatrix)
	require.Equal(t, true, ok)
	require.LogsContain(t, hook, "Computed state root")
	require.DeepEqual(t, emptyPayload, bellatrixBlk.Bellatrix.Body.ExecutionPayload) // Payload should equal.
}

func TestServer_GetBellatrixBeaconBlock_BuilderCase(t *testing.T) {
	db := dbTest.SetupDB(t)
	ctx := context.Background()
	hook := logTest.NewGlobal()

	params.SetupTestConfigCleanup(t)
	cfg := params.BeaconConfig().Copy()
	cfg.BellatrixForkEpoch = 2
	cfg.AltairForkEpoch = 1
	params.OverrideBeaconConfig(cfg)

	beaconState, privKeys := util.DeterministicGenesisState(t, 64)
	stateRoot, err := beaconState.HashTreeRoot(ctx)
	require.NoError(t, err, "Could not hash genesis state")

	genesis := consensusblocks.NewGenesisBlock(stateRoot[:])
	wsb, err := blocks.NewSignedBeaconBlock(genesis)
	require.NoError(t, err)
	require.NoError(t, db.SaveBlock(ctx, wsb), "Could not save genesis block")

	parentRoot, err := genesis.Block.HashTreeRoot()
	require.NoError(t, err, "Could not get signing root")
	require.NoError(t, db.SaveState(ctx, beaconState, parentRoot), "Could not save genesis state")
	require.NoError(t, db.SaveHeadBlockRoot(ctx, parentRoot), "Could not save genesis state")

	bellatrixSlot, err := slots.EpochStart(params.BeaconConfig().BellatrixForkEpoch)
	require.NoError(t, err)

	emptyPayload := &v1.ExecutionPayload{
		ParentHash:    make([]byte, fieldparams.RootLength),
		FeeRecipient:  make([]byte, fieldparams.FeeRecipientLength),
		StateRoot:     make([]byte, fieldparams.RootLength),
		ReceiptsRoot:  make([]byte, fieldparams.RootLength),
		LogsBloom:     make([]byte, fieldparams.LogsBloomLength),
		PrevRandao:    make([]byte, fieldparams.RootLength),
		BaseFeePerGas: make([]byte, fieldparams.RootLength),
		BlockHash:     make([]byte, fieldparams.RootLength),
	}
	var scBits [fieldparams.SyncAggregateSyncCommitteeBytesLength]byte
	blk := &ethpb.SignedBeaconBlockBellatrix{
		Block: &ethpb.BeaconBlockBellatrix{
			Slot:       bellatrixSlot + 1,
			ParentRoot: parentRoot[:],
			StateRoot:  genesis.Block.StateRoot,
			Body: &ethpb.BeaconBlockBodyBellatrix{
				RandaoReveal:     genesis.Block.Body.RandaoReveal,
				Graffiti:         genesis.Block.Body.Graffiti,
				Eth1Data:         genesis.Block.Body.Eth1Data,
				SyncAggregate:    &ethpb.SyncAggregate{SyncCommitteeBits: scBits[:], SyncCommitteeSignature: make([]byte, 96)},
				ExecutionPayload: emptyPayload,
			},
		},
		Signature: genesis.Signature,
	}

	blkRoot, err := blk.Block.HashTreeRoot()
	require.NoError(t, err)
	require.NoError(t, err, "Could not get signing root")
	require.NoError(t, db.SaveState(ctx, beaconState, blkRoot), "Could not save genesis state")
	require.NoError(t, db.SaveHeadBlockRoot(ctx, blkRoot), "Could not save genesis state")

	b1 := util.NewBeaconBlockBellatrix()
	b1.Block.Body.ExecutionPayload.BlockNumber = 1 // Execution enabled.
	wb1, err := blocks.NewSignedBeaconBlock(b1)
	require.NoError(t, err)
	wbr1, err := wb1.Block().HashTreeRoot()
	require.NoError(t, err)
	require.NoError(t, db.SaveBlock(ctx, wb1))

	random, err := helpers.RandaoMix(beaconState, prysmtime.CurrentEpoch(beaconState))
	require.NoError(t, err)

	tstamp, err := slots.ToTime(beaconState.GenesisTime(), bellatrixSlot+1)
	require.NoError(t, err)
	h := &v1.ExecutionPayloadHeader{
		BlockNumber:      123,
		GasLimit:         456,
		GasUsed:          789,
		ParentHash:       make([]byte, fieldparams.RootLength),
		FeeRecipient:     make([]byte, fieldparams.FeeRecipientLength),
		StateRoot:        make([]byte, fieldparams.RootLength),
		ReceiptsRoot:     make([]byte, fieldparams.RootLength),
		LogsBloom:        make([]byte, fieldparams.LogsBloomLength),
		PrevRandao:       random,
		BaseFeePerGas:    make([]byte, fieldparams.RootLength),
		BlockHash:        make([]byte, fieldparams.RootLength),
		TransactionsRoot: make([]byte, fieldparams.RootLength),
		ExtraData:        make([]byte, 0),
		Timestamp:        uint64(tstamp.Unix()),
	}

	proposerServer := &Server{
		FinalizationFetcher: &blockchainTest.ChainService{FinalizedCheckPoint: &ethpb.Checkpoint{Root: wbr1[:]}},
		HeadFetcher:         &blockchainTest.ChainService{State: beaconState, Root: parentRoot[:], Optimistic: false, Block: wb1},
		TimeFetcher:         &blockchainTest.ChainService{Genesis: time.Unix(int64(beaconState.GenesisTime()), 0)},
		SyncChecker:         &mockSync.Sync{IsSyncing: false},
		BlockReceiver:       &blockchainTest.ChainService{},
		HeadUpdater:         &blockchainTest.ChainService{},
		ForkFetcher:         &blockchainTest.ChainService{Fork: &ethpb.Fork{}},
		GenesisFetcher:      &blockchainTest.ChainService{},
		ChainStartFetcher:   &mockExecution.Chain{},
		Eth1InfoFetcher:     &mockExecution.Chain{},
		MockEth1Votes:       true,
		AttPool:             attestations.NewPool(),
		SlashingsPool:       slashings.NewPool(),
		ExitPool:            voluntaryexits.NewPool(),
		StateGen:            stategen.New(db, doublylinkedtree.New()),
		SyncCommitteePool:   synccommittee.NewStore(),
		ExecutionEngineCaller: &mockExecution.EngineClient{
			PayloadIDBytes:   &v1.PayloadIDBytes{1},
			ExecutionPayload: emptyPayload,
		},
		BeaconDB: db,
	}
	sk, err := bls.RandKey()
	require.NoError(t, err)
	bid := &ethpb.BuilderBid{
		Header: h,
		Pubkey: sk.PublicKey().Marshal(),
		Value:  bytesutil.PadTo([]byte{1, 2, 3}, 32),
	}
	d := params.BeaconConfig().DomainApplicationBuilder
	domain, err := signing.ComputeDomain(d, nil, nil)
	require.NoError(t, err)
	sr, err := signing.ComputeSigningRoot(bid, domain)
	require.NoError(t, err)
	sBid := &ethpb.SignedBuilderBid{
		Message:   bid,
		Signature: sk.Sign(sr[:]).Marshal(),
	}
	proposerServer.BlockBuilder = &builderTest.MockBuilderService{HasConfigured: true, Bid: sBid}
	proposerServer.ForkFetcher = &blockchainTest.ChainService{ForkChoiceStore: doublylinkedtree.New()}
	randaoReveal, err := util.RandaoReveal(beaconState, 0, privKeys)
	require.NoError(t, err)

	require.NoError(t, proposerServer.BeaconDB.SaveRegistrationsByValidatorIDs(ctx, []types.ValidatorIndex{11},
		[]*ethpb.ValidatorRegistrationV1{{FeeRecipient: bytesutil.PadTo([]byte{}, fieldparams.FeeRecipientLength), Pubkey: bytesutil.PadTo([]byte{}, fieldparams.BLSPubkeyLength)}}))

	params.SetupTestConfigCleanup(t)
	cfg.MaxBuilderConsecutiveMissedSlots = bellatrixSlot + 1
	cfg.MaxBuilderEpochMissedSlots = 32
	params.OverrideBeaconConfig(cfg)

	block, err := proposerServer.getBellatrixBeaconBlock(ctx, &ethpb.BlockRequest{
		Slot:         bellatrixSlot + 1,
		RandaoReveal: randaoReveal,
	})
	require.NoError(t, err)
	bellatrixBlk, ok := block.GetBlock().(*ethpb.GenericBeaconBlock_BlindedBellatrix)
	require.Equal(t, true, ok)
	require.LogsContain(t, hook, "Computed state root")
	require.DeepEqual(t, h, bellatrixBlk.BlindedBellatrix.Body.ExecutionPayloadHeader) // Payload header should equal.
}

func TestServer_validateBuilderSignature(t *testing.T) {
	sk, err := bls.RandKey()
	require.NoError(t, err)
	bid := &ethpb.BuilderBid{
		Header: &v1.ExecutionPayloadHeader{
			ParentHash:       make([]byte, fieldparams.RootLength),
			FeeRecipient:     make([]byte, fieldparams.FeeRecipientLength),
			StateRoot:        make([]byte, fieldparams.RootLength),
			ReceiptsRoot:     make([]byte, fieldparams.RootLength),
			LogsBloom:        make([]byte, fieldparams.LogsBloomLength),
			PrevRandao:       make([]byte, fieldparams.RootLength),
			BaseFeePerGas:    make([]byte, fieldparams.RootLength),
			BlockHash:        make([]byte, fieldparams.RootLength),
			TransactionsRoot: make([]byte, fieldparams.RootLength),
			BlockNumber:      1,
		},
		Pubkey: sk.PublicKey().Marshal(),
		Value:  bytesutil.PadTo([]byte{1, 2, 3}, 32),
	}
	d := params.BeaconConfig().DomainApplicationBuilder
	domain, err := signing.ComputeDomain(d, nil, nil)
	require.NoError(t, err)
	sr, err := signing.ComputeSigningRoot(bid, domain)
	require.NoError(t, err)
	sBid := &ethpb.SignedBuilderBid{
		Message:   bid,
		Signature: sk.Sign(sr[:]).Marshal(),
	}
	require.NoError(t, validateBuilderSignature(sBid))

	sBid.Message.Value = make([]byte, 32)
	require.ErrorIs(t, validateBuilderSignature(sBid), signing.ErrSigFailedToVerify)
}
