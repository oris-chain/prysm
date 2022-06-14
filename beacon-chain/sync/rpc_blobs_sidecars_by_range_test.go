package sync

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/kevinms/leakybucket-go"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/protocol"
	chainMock "github.com/prysmaticlabs/prysm/beacon-chain/blockchain/testing"
	db "github.com/prysmaticlabs/prysm/beacon-chain/db/testing"
	"github.com/prysmaticlabs/prysm/beacon-chain/p2p"
	p2ptest "github.com/prysmaticlabs/prysm/beacon-chain/p2p/testing"
	fieldparams "github.com/prysmaticlabs/prysm/config/fieldparams"
	types "github.com/prysmaticlabs/prysm/consensus-types/primitives"
	"github.com/prysmaticlabs/prysm/consensus-types/wrapper"
	ethpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/testing/assert"
	"github.com/prysmaticlabs/prysm/testing/require"
	"github.com/prysmaticlabs/prysm/testing/util"
)

func newBlobsSidecar() *ethpb.SignedBlobsSidecar {
	return &ethpb.SignedBlobsSidecar{
		Message: &ethpb.BlobsSidecar{
			BeaconBlockRoot: make([]byte, fieldparams.RootLength),
		},
		Signature: make([]byte, fieldparams.BLSSignatureLength),
	}
}

func TestRPCBlobsSidecarsByRange_RPCHandlerReturnsBlobsSidecars(t *testing.T) {
	p1 := p2ptest.NewTestP2P(t)
	p2 := p2ptest.NewTestP2P(t)
	p1.Connect(p2)
	assert.Equal(t, 1, len(p1.BHost.Network().Peers()), "Expected peers to be connected")
	d := db.SetupDB(t)

	req := &ethpb.BlobsSidecarsByRangeRequest{
		StartSlot: 100,
		Count:     4,
	}

	for i := req.StartSlot; i < req.StartSlot.Add(req.Count); i += types.Slot(1) {
		// save the BeaconBlock to index the slots used to retrieve sidecars
		blk := util.NewBeaconBlock()
		blk.Block.Slot = i
		wsb, err := wrapper.WrappedSignedBeaconBlock(blk)
		require.NoError(t, err)
		require.NoError(t, d.SaveBlock(context.Background(), wsb))

		sidecar := newBlobsSidecar()
		root, err := blk.Block.HashTreeRoot()
		require.NoError(t, err)
		sidecar.Message.BeaconBlockRoot = root[:]
		sidecar.Message.BeaconBlockSlot = blk.Block.Slot
		require.NoError(t, d.SaveBlobsSidecar(context.Background(), sidecar.Message))
	}

	// Start service with 160 as allowed blocks capacity (and almost zero capacity recovery).
	r := &Service{cfg: &config{p2p: p1, beaconDB: d, chain: &chainMock.ChainService{}}, rateLimiter: newRateLimiter(p1)}
	pcl := protocol.ID(p2p.RPCBlobsSidecarsByRangeTopicV1)
	topic := string(pcl)
	r.rateLimiter.limiterMap[topic] = leakybucket.NewCollector(1000, 10000, false)
	var wg sync.WaitGroup
	wg.Add(1)
	p2.BHost.SetStreamHandler(pcl, func(stream network.Stream) {
		defer wg.Done()
		for i := req.StartSlot; i < req.StartSlot.Add(req.Count); i += types.Slot(1) {
			expectSuccess(t, stream)
			sidecar := new(ethpb.BlobsSidecar)
			assert.NoError(t, r.cfg.p2p.Encoding().DecodeWithMaxLength(stream, sidecar))
		}
	})

	stream1, err := p1.BHost.NewStream(context.Background(), p2.BHost.ID(), pcl)
	require.NoError(t, err)

	err = r.blobsSidecarsByRangeRPCHandler(context.Background(), req, stream1)
	require.NoError(t, err)

	// Make sure that rate limiter doesn't limit capacity exceedingly.
	remainingCapacity := r.rateLimiter.limiterMap[topic].Remaining(p2.PeerID().String())
	expectedCapacity := int64(10000 - 40*req.Count) // an empty sidecar is 40 bytes
	require.Equal(t, expectedCapacity, remainingCapacity, "Unexpected rate limiting capacity")

	if util.WaitTimeout(&wg, 1*time.Second) {
		t.Fatal("Did not receive stream within 1 sec")
	}
}
