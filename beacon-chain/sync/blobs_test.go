package sync

import (
	"context"
	"encoding/binary"
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
	mock "github.com/prysmaticlabs/prysm/v4/beacon-chain/blockchain/testing"
	db "github.com/prysmaticlabs/prysm/v4/beacon-chain/db/testing"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/p2p"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/p2p/encoder"
	p2ptest "github.com/prysmaticlabs/prysm/v4/beacon-chain/p2p/testing"
	fieldparams "github.com/prysmaticlabs/prysm/v4/config/fieldparams"
	"github.com/prysmaticlabs/prysm/v4/config/params"
	types "github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	leakybucket "github.com/prysmaticlabs/prysm/v4/container/leaky-bucket"
	"github.com/prysmaticlabs/prysm/v4/encoding/bytesutil"
	"github.com/prysmaticlabs/prysm/v4/network/forks"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
	ethpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v4/testing/require"
	"github.com/prysmaticlabs/prysm/v4/testing/util"
	"github.com/prysmaticlabs/prysm/v4/time/slots"
)

type blobsTestCase struct {
	name                string
	nblocks             int                  // how many blocks to loop through in setting up test fixtures & requests
	missing             map[int]map[int]bool // skip this blob index, so that we can test different custody scenarios
	expired             map[int]bool         // mark block expired to test scenarios where requests are outside retention window
	chain               *mock.ChainService   // allow tests to control retention window via current slot and finalized checkpoint
	total               *int                 // allow a test to specify the total number of responses received
	err                 error
	serverHandle        testHandler
	defineExpected      expectedDefiner
	requestFromSidecars requestFromSidecars
	topic               protocol.ID
	oldestSlot          oldestSlotCallback
}

type testHandler func(s *Service) rpcHandler
type expectedDefiner func(t *testing.T, scs []*ethpb.BlobSidecar, req interface{}) []*expectedBlobChunk
type requestFromSidecars func([]*ethpb.BlobSidecar) interface{}
type oldestSlotCallback func(t *testing.T) types.Slot

func generateTestBlockWithSidecars(t *testing.T, parent [32]byte, slot types.Slot, nblobs int) (*ethpb.SignedBeaconBlockDeneb, []*ethpb.BlobSidecar) {
	// Start service with 160 as allowed blocks capacity (and almost zero capacity recovery).
	stateRoot := bytesutil.PadTo([]byte("stateRoot"), fieldparams.RootLength)
	receiptsRoot := bytesutil.PadTo([]byte("receiptsRoot"), fieldparams.RootLength)
	logsBloom := bytesutil.PadTo([]byte("logs"), fieldparams.LogsBloomLength)
	parentHash := bytesutil.PadTo([]byte("parentHash"), fieldparams.RootLength)
	tx := gethTypes.NewTransaction(
		0,
		common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87"),
		big.NewInt(0), 0, big.NewInt(0),
		nil,
	)
	txs := []*gethTypes.Transaction{tx}
	encodedBinaryTxs := make([][]byte, 1)
	var err error
	encodedBinaryTxs[0], err = txs[0].MarshalBinary()
	require.NoError(t, err)
	blockHash := bytesutil.ToBytes32([]byte("foo"))
	payload := &enginev1.ExecutionPayloadDeneb{
		ParentHash:    parentHash,
		FeeRecipient:  make([]byte, fieldparams.FeeRecipientLength),
		StateRoot:     stateRoot,
		ReceiptsRoot:  receiptsRoot,
		LogsBloom:     logsBloom,
		PrevRandao:    blockHash[:],
		BlockNumber:   0,
		GasLimit:      0,
		GasUsed:       0,
		Timestamp:     0,
		ExtraData:     make([]byte, 0),
		BaseFeePerGas: bytesutil.PadTo([]byte("baseFeePerGas"), fieldparams.RootLength),
		ExcessDataGas: bytesutil.PadTo([]byte("excessDataGas"), fieldparams.RootLength),
		BlockHash:     blockHash[:],
		Transactions:  encodedBinaryTxs,
	}
	block := util.NewBeaconBlockDeneb()
	block.Block.Body.ExecutionPayload = payload
	block.Block.Slot = slot
	block.Block.ParentRoot = parent[:]
	commitments := make([][48]byte, nblobs)
	block.Block.Body.BlobKzgCommitments = make([][]byte, nblobs)
	for i := range commitments {
		binary.LittleEndian.PutUint64(commitments[i][:], uint64(i))
		block.Block.Body.BlobKzgCommitments[i] = commitments[i][:]
	}

	root, err := block.Block.HashTreeRoot()
	require.NoError(t, err)

	sidecars := make([]*ethpb.BlobSidecar, len(commitments))
	for i, c := range block.Block.Body.BlobKzgCommitments {
		sidecars[i] = generateTestSidecar(root, block, i, c)
	}
	return block, sidecars
}

func generateTestSidecar(root [32]byte, block *ethpb.SignedBeaconBlockDeneb, index int, commitment []byte) *ethpb.BlobSidecar {
	blob := &enginev1.Blob{
		Data: make([]byte, fieldparams.BlobSize),
	}
	binary.LittleEndian.PutUint64(blob.Data, uint64(index))
	sc := &ethpb.BlobSidecar{
		BlockRoot:       root[:],
		Index:           uint64(index),
		Slot:            block.Block.Slot,
		BlockParentRoot: block.Block.ParentRoot,
		ProposerIndex:   block.Block.ProposerIndex,
		Blob:            blob,
		KzgCommitment:   commitment,
		KzgProof:        commitment,
	}
	return sc
}

type expectedBlobChunk struct {
	code    uint8
	sidecar *ethpb.BlobSidecar
	message string
}

func (r *expectedBlobChunk) requireExpected(t *testing.T, s *Service, stream network.Stream) {
	d := s.cfg.p2p.Encoding().DecodeWithMaxLength

	code, _, err := ReadStatusCode(stream, &encoder.SszNetworkEncoder{})
	require.NoError(t, err)
	require.Equal(t, r.code, code, "unexpected response code")
	//require.Equal(t, r.message, msg, "unexpected error message")
	if code != responseCodeSuccess {
		return
	}

	c, err := readContextFromStream(stream, s.cfg.chain)
	require.NoError(t, err)

	valRoot := s.cfg.chain.GenesisValidatorsRoot()
	ctxBytes, err := forks.ForkDigestFromEpoch(slots.ToEpoch(r.sidecar.GetSlot()), valRoot[:])
	require.NoError(t, err)
	require.Equal(t, ctxBytes, bytesutil.ToBytes4(c))

	sc := &ethpb.BlobSidecar{}
	require.NoError(t, d(stream, sc))
	require.Equal(t, bytesutil.ToBytes32(sc.BlockRoot), bytesutil.ToBytes32(r.sidecar.BlockRoot))
	require.Equal(t, sc.Index, r.sidecar.Index)
}

func (c *blobsTestCase) setup(t *testing.T) (*Service, []*ethpb.BlobSidecar, func()) {
	cfg := params.BeaconConfig()
	repositionFutureEpochs(cfg)
	undo, err := params.SetActiveWithUndo(cfg)
	require.NoError(t, err)
	cleanup := func() {
		require.NoError(t, undo())
	}
	maxBlobs := int(params.BeaconConfig().MaxBlobsPerBlock)
	if c.chain == nil {
		c.chain = defaultMockChain(t)
	}
	d := db.SetupDB(t)

	sidecars := make([]*ethpb.BlobSidecar, 0)
	oldest := c.oldestSlot(t)
	var parentRoot [32]byte
	for i := 0; i < c.nblocks; i++ {
		// check if there is a slot override for this index
		// ie to create a block outside the minimum_request_epoch
		var bs types.Slot
		if c.expired[i] {
			// the lowest possible bound of the retention period is the deneb epoch, so make sure
			// the slot of an expired block is at least one slot less than the deneb epoch.
			bs = oldest - 1 - types.Slot(i)
		} else {
			bs = oldest + types.Slot(i)
		}
		block, bsc := generateTestBlockWithSidecars(t, parentRoot, bs, maxBlobs)
		root, err := block.Block.HashTreeRoot()
		require.NoError(t, err)
		for _, sc := range bsc {
			sidecars = append(sidecars, sc)
		}
		util.SaveBlock(t, context.Background(), d, block)
		parentRoot = root
	}

	client := p2ptest.NewTestP2P(t)
	s := &Service{
		cfg:         &config{p2p: client, chain: c.chain, beaconDB: d},
		rateLimiter: newRateLimiter(client),
	}

	byRootRate := params.BeaconNetworkConfig().MaxRequestBlobsSidecars * params.BeaconConfig().MaxBlobsPerBlock
	byRangeRate := params.BeaconNetworkConfig().MaxRequestBlobsSidecars * params.BeaconConfig().MaxBlobsPerBlock
	s.setRateCollector(p2p.RPCBlobSidecarsByRootTopicV1, leakybucket.NewCollector(0.000001, int64(byRootRate), time.Second, false))
	s.setRateCollector(p2p.RPCBlobSidecarsByRangeTopicV1, leakybucket.NewCollector(0.000001, int64(byRangeRate), time.Second, false))

	return s, sidecars, cleanup
}

func (c *blobsTestCase) run(t *testing.T) {
	s, sidecars, cleanup := c.setup(t)
	defer cleanup()
	req := c.requestFromSidecars(sidecars)
	expect := c.defineExpected(t, sidecars, req)
	m := map[types.Slot][]*ethpb.BlobSidecar{}
	for _, sc := range expect {
		m[sc.sidecar.Slot] = append(m[sc.sidecar.Slot], sc.sidecar)
	}
	for _, blobSidecars := range m {
		require.NoError(t, s.cfg.beaconDB.SaveBlobSidecar(context.Background(), blobSidecars))
	}
	if c.total != nil {
		require.Equal(t, *c.total, len(expect))
	}
	nh := func(stream network.Stream) {
		for _, ex := range expect {
			ex.requireExpected(t, s, stream)
		}
	}
	rht := &rpcHandlerTest{
		t:       t,
		topic:   c.topic,
		timeout: time.Second * 10,
		err:     c.err,
		s:       s,
	}
	rht.testHandler(nh, c.serverHandle(s), req)
}

// we use max uints for future forks, but this causes overflows when computing slots
// so it is helpful in tests to temporarily reposition the epochs to give room for some math.
func repositionFutureEpochs(cfg *params.BeaconChainConfig) {
	if cfg.CapellaForkEpoch == math.MaxUint64 {
		cfg.CapellaForkEpoch = cfg.BellatrixForkEpoch + 100
	}
	if cfg.DenebForkEpoch == math.MaxUint64 {
		cfg.DenebForkEpoch = cfg.CapellaForkEpoch + 100
	}
}

func defaultMockChain(t *testing.T) *mock.ChainService {
	df, err := forks.Fork(params.BeaconConfig().DenebForkEpoch)
	require.NoError(t, err)
	ce := params.BeaconConfig().DenebForkEpoch + params.BeaconNetworkConfig().MinEpochsForBlobsSidecarsRequest + 1000
	fe := ce - 2
	cs, err := slots.EpochStart(ce)
	require.NoError(t, err)

	return &mock.ChainService{
		ValidatorsRoot:      [32]byte{},
		Slot:                &cs,
		FinalizedCheckPoint: &ethpb.Checkpoint{Epoch: fe},
		Fork:                df}
}

func TestTestcaseSetup_BlocksAndBlobs(t *testing.T) {
	ctx := context.Background()
	c := &blobsTestCase{nblocks: 10}
	c.oldestSlot = c.defaultOldestSlotByRoot
	s, sidecars, cleanup := c.setup(t)
	req := blobRootRequestFromSidecars(sidecars)
	expect := c.filterExpectedByRoot(t, sidecars, req)
	defer cleanup()
	require.Equal(t, 40, len(sidecars))
	require.Equal(t, 40, len(expect))
	for _, sc := range sidecars {
		blk, err := s.cfg.beaconDB.Block(ctx, bytesutil.ToBytes32(sc.BlockRoot))
		require.NoError(t, err)
		var found *int
		comms, err := blk.Block().Body().BlobKzgCommitments()
		require.NoError(t, err)
		for i, cm := range comms {
			if bytesutil.ToBytes48(sc.KzgCommitment) == bytesutil.ToBytes48(cm) {
				found = &i
			}
		}
		require.Equal(t, true, found != nil)
	}
}

func TestRoundTripDenebSave(t *testing.T) {
	ctx := context.Background()
	cfg := params.BeaconConfig()
	repositionFutureEpochs(cfg)
	undo, err := params.SetActiveWithUndo(cfg)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, undo())
	}()
	parentRoot := [32]byte{}
	c := blobsTestCase{nblocks: 10}
	c.chain = defaultMockChain(t)
	oldest, err := slots.EpochStart(blobMinReqEpoch(c.chain.FinalizedCheckPoint.Epoch, slots.ToEpoch(c.chain.CurrentSlot())))
	require.NoError(t, err)
	maxBlobs := int(params.BeaconConfig().MaxBlobsPerBlock)
	block, bsc := generateTestBlockWithSidecars(t, parentRoot, oldest, maxBlobs)
	require.Equal(t, len(block.Block.Body.BlobKzgCommitments), len(bsc))
	require.Equal(t, maxBlobs, len(bsc))
	for i := range bsc {
		require.DeepEqual(t, block.Block.Body.BlobKzgCommitments[i], bsc[i].KzgCommitment)
	}
	d := db.SetupDB(t)
	util.SaveBlock(t, ctx, d, block)
	root, err := block.Block.HashTreeRoot()
	require.NoError(t, err)
	dbBlock, err := d.Block(ctx, root)
	require.NoError(t, err)
	comms, err := dbBlock.Block().Body().BlobKzgCommitments()
	require.NoError(t, err)
	require.Equal(t, maxBlobs, len(comms))
	for i := range bsc {
		require.DeepEqual(t, comms[i], bsc[i].KzgCommitment)
	}
}
