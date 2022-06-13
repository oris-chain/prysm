package validator

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	ct "github.com/prysmaticlabs/prysm/beacon-chain/blockchain/testing"
	bt "github.com/prysmaticlabs/prysm/beacon-chain/builder/testing"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/altair"
	dbTest "github.com/prysmaticlabs/prysm/beacon-chain/db/testing"
	"github.com/prysmaticlabs/prysm/beacon-chain/state/stategen"
	fieldparams "github.com/prysmaticlabs/prysm/config/fieldparams"
	"github.com/prysmaticlabs/prysm/config/params"
	"github.com/prysmaticlabs/prysm/consensus-types/interfaces"
	"github.com/prysmaticlabs/prysm/consensus-types/wrapper"
	"github.com/prysmaticlabs/prysm/encoding/bytesutil"
	v1 "github.com/prysmaticlabs/prysm/proto/engine/v1"
	ethpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/testing/require"
	"github.com/prysmaticlabs/prysm/testing/util"
)

func TestServer_buildHeaderBlock(t *testing.T) {
	db := dbTest.SetupDB(t)
	ctx := context.Background()

	params.SetupTestConfigCleanup(t)
	params.OverrideBeaconConfig(params.MainnetConfig())
	beaconState, keys := util.DeterministicGenesisStateAltair(t, 16384)
	sCom, err := altair.NextSyncCommittee(context.Background(), beaconState)
	require.NoError(t, err)
	require.NoError(t, beaconState.SetCurrentSyncCommittee(sCom))
	copiedState := beaconState.Copy()

	proposerServer := &Server{
		BeaconDB: db,
		StateGen: stategen.New(db),
	}
	b, err := util.GenerateFullBlockAltair(copiedState, keys, util.DefaultBlockGenConfig(), 1)
	require.NoError(t, err)
	r := bytesutil.ToBytes32(b.Block.ParentRoot)
	wb, err := wrapper.WrappedSignedBeaconBlock(b)
	require.NoError(t, err)
	require.NoError(t, proposerServer.BeaconDB.SaveBlock(ctx, wb))
	require.NoError(t, proposerServer.BeaconDB.SaveState(ctx, beaconState, r))

	b1, err := util.GenerateFullBlockAltair(copiedState, keys, util.DefaultBlockGenConfig(), 2)
	require.NoError(t, err)

	vs := &Server{StateGen: stategen.New(db), BeaconDB: db}
	h := &ethpb.ExecutionPayloadHeader{
		BlockNumber:      123,
		GasLimit:         456,
		GasUsed:          789,
		ParentHash:       make([]byte, fieldparams.RootLength),
		FeeRecipient:     make([]byte, fieldparams.FeeRecipientLength),
		StateRoot:        make([]byte, fieldparams.RootLength),
		ReceiptsRoot:     make([]byte, fieldparams.RootLength),
		LogsBloom:        make([]byte, fieldparams.LogsBloomLength),
		PrevRandao:       make([]byte, fieldparams.RootLength),
		BaseFeePerGas:    make([]byte, fieldparams.RootLength),
		BlockHash:        make([]byte, fieldparams.RootLength),
		TransactionsRoot: make([]byte, fieldparams.RootLength),
		ExtraData:        make([]byte, 0),
	}
	got, err := vs.buildHeaderBlock(ctx, b1.Block, h)
	require.NoError(t, err)
	require.DeepEqual(t, h, got.GetBlindedBellatrix().Body.ExecutionPayloadHeader)
}

func TestServer_readyForBuilder(t *testing.T) {
	vs := &Server{BeaconDB: dbTest.SetupDB(t)}
	b := util.NewBeaconBlockBellatrix()
	wb, err := wrapper.WrappedSignedBeaconBlock(b)
	require.NoError(t, err)
	wbr, err := wb.Block().HashTreeRoot()
	require.NoError(t, err)
	b1 := util.NewBeaconBlockBellatrix()
	b1.Block.Body.ExecutionPayload.BlockNumber = 1
	wb1, err := wrapper.WrappedSignedBeaconBlock(b1)
	require.NoError(t, err)
	wbr1, err := wb1.Block().HashTreeRoot()
	require.NoError(t, err)
	ctx := context.Background()
	require.NoError(t, vs.BeaconDB.SaveBlock(ctx, wb))
	require.NoError(t, vs.BeaconDB.SaveBlock(ctx, wb1))
	cs := &ct.ChainService{FinalizedCheckPoint: &ethpb.Checkpoint{Root: wbr[:]}}
	vs.FinalizationFetcher = cs
	ready, err := vs.readyForBuilder(ctx)
	require.NoError(t, err)
	require.Equal(t, false, ready)

	cs = &ct.ChainService{FinalizedCheckPoint: &ethpb.Checkpoint{Root: wbr1[:]}}
	vs.FinalizationFetcher = cs
	ready, err = vs.readyForBuilder(ctx)
	require.NoError(t, err)
	require.Equal(t, true, ready)
}

func TestServer_getPayloadHeader(t *testing.T) {
	tests := []struct {
		name           string
		head           interfaces.SignedBeaconBlock
		mock           *bt.MockBuilderService
		fetcher        *ct.ChainService
		err            string
		returnedHeader *ethpb.ExecutionPayloadHeader
	}{
		{
			name: "builder is not ready",
			mock: &bt.MockBuilderService{
				ErrStatus: errors.New("builder is not ready"),
			},
			err: "builder is not ready",
		},
		{
			name: "head is not bellatrix ready",
			mock: &bt.MockBuilderService{},
			fetcher: &ct.ChainService{
				Block: func() interfaces.SignedBeaconBlock {
					wb, err := wrapper.WrappedSignedBeaconBlock(util.NewBeaconBlock())
					require.NoError(t, err)
					return wb
				}(),
			},
		},
		{
			name: "get header failed",
			mock: &bt.MockBuilderService{
				ErrGetHeader: errors.New("can't get header"),
			},
			fetcher: &ct.ChainService{
				Block: func() interfaces.SignedBeaconBlock {
					wb, err := wrapper.WrappedSignedBeaconBlock(util.NewBeaconBlockBellatrix())
					require.NoError(t, err)
					return wb
				}(),
			},
			err: "can't get header",
		},
		{
			name: "get header correct",
			mock: &bt.MockBuilderService{
				Bid: &ethpb.SignedBuilderBid{
					Message: &ethpb.BuilderBid{
						Header: &ethpb.ExecutionPayloadHeader{
							BlockNumber: 123,
						},
					},
				},
				ErrGetHeader: errors.New("can't get header"),
			},
			fetcher: &ct.ChainService{
				Block: func() interfaces.SignedBeaconBlock {
					wb, err := wrapper.WrappedSignedBeaconBlock(util.NewBeaconBlockBellatrix())
					require.NoError(t, err)
					return wb
				}(),
			},
			returnedHeader: &ethpb.ExecutionPayloadHeader{
				BlockNumber: 123,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			vs := &Server{BlockBuilder: tc.mock, HeadFetcher: tc.fetcher}
			h, err := vs.getPayloadHeader(context.Background(), 0, 0)
			if err != nil {
				require.ErrorContains(t, tc.err, err)
			} else {
				require.DeepEqual(t, tc.returnedHeader, h)
			}
		})
	}
}

func TestServer_getBuilderBlock(t *testing.T) {
	tests := []struct {
		name        string
		blk         interfaces.SignedBeaconBlock
		mock        *bt.MockBuilderService
		err         string
		returnedBlk interfaces.SignedBeaconBlock
	}{
		{
			name: "old block version",
			blk: func() interfaces.SignedBeaconBlock {
				wb, err := wrapper.WrappedSignedBeaconBlock(util.NewBeaconBlock())
				require.NoError(t, err)
				return wb
			}(),
			returnedBlk: func() interfaces.SignedBeaconBlock {
				wb, err := wrapper.WrappedSignedBeaconBlock(util.NewBeaconBlock())
				require.NoError(t, err)
				return wb
			}(),
		},
		{
			name: "not configured",
			blk: func() interfaces.SignedBeaconBlock {
				wb, err := wrapper.WrappedSignedBeaconBlock(util.NewBlindedBeaconBlockBellatrix())
				require.NoError(t, err)
				return wb
			}(),
			mock: &bt.MockBuilderService{
				HasConfigured: false,
			},
			returnedBlk: func() interfaces.SignedBeaconBlock {
				wb, err := wrapper.WrappedSignedBeaconBlock(util.NewBlindedBeaconBlockBellatrix())
				require.NoError(t, err)
				return wb
			}(),
		},
		{
			name: "builder is not ready",
			blk: func() interfaces.SignedBeaconBlock {
				wb, err := wrapper.WrappedSignedBeaconBlock(util.NewBlindedBeaconBlockBellatrix())
				require.NoError(t, err)
				return wb
			}(),
			mock: &bt.MockBuilderService{
				HasConfigured: true,
				ErrStatus:     errors.New("builder is not ready"),
			},
			err: "builder is not ready",
		},
		{
			name: "submit blind block error",
			blk: func() interfaces.SignedBeaconBlock {
				b := util.NewBlindedBeaconBlockBellatrix()
				b.Block.Slot = 1
				b.Block.ProposerIndex = 2
				wb, err := wrapper.WrappedSignedBeaconBlock(b)
				require.NoError(t, err)
				return wb
			}(),
			mock: &bt.MockBuilderService{
				HasConfigured:         true,
				ErrSubmitBlindedBlock: errors.New("can't submit"),
			},
			err: "can't submit",
		},
		{
			name: "can submit block",
			blk: func() interfaces.SignedBeaconBlock {
				b := util.NewBlindedBeaconBlockBellatrix()
				b.Block.Slot = 1
				b.Block.ProposerIndex = 2
				wb, err := wrapper.WrappedSignedBeaconBlock(b)
				require.NoError(t, err)
				return wb
			}(),
			mock: &bt.MockBuilderService{
				HasConfigured: true,
				Payload:       &v1.ExecutionPayload{GasLimit: 123},
			},
			returnedBlk: func() interfaces.SignedBeaconBlock {
				b := util.NewBeaconBlockBellatrix()
				b.Block.Slot = 1
				b.Block.ProposerIndex = 2
				b.Block.Body.ExecutionPayload = &v1.ExecutionPayload{GasLimit: 123}
				wb, err := wrapper.WrappedSignedBeaconBlock(b)
				require.NoError(t, err)
				return wb
			}(),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			vs := &Server{BlockBuilder: tc.mock}
			gotBlk, err := vs.getBuilderBlock(context.Background(), tc.blk)
			if err != nil {
				require.ErrorContains(t, tc.err, err)
			} else {
				require.DeepEqual(t, tc.returnedBlk, gotBlk)
			}
		})
	}
}
