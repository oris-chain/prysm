package stategen

import (
	"context"
	"testing"

	testDB "github.com/prysmaticlabs/prysm/v3/beacon-chain/db/testing"
	doublylinkedtree "github.com/prysmaticlabs/prysm/v3/beacon-chain/forkchoice/doubly-linked-tree"
	forkchoicetypes "github.com/prysmaticlabs/prysm/v3/beacon-chain/forkchoice/types"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/state"
	v1 "github.com/prysmaticlabs/prysm/v3/beacon-chain/state/v1"
	"github.com/prysmaticlabs/prysm/v3/config/params"
	types "github.com/prysmaticlabs/prysm/v3/consensus-types/primitives"
	ethpb "github.com/prysmaticlabs/prysm/v3/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v3/testing/assert"
	"github.com/prysmaticlabs/prysm/v3/testing/require"
	"github.com/prysmaticlabs/prysm/v3/testing/util"
	"github.com/prysmaticlabs/prysm/v3/time/slots"
)

func TestHotStateCache_RoundTrip(t *testing.T) {
	c := newHotStateCache()
	root := [32]byte{'A'}
	s := c.get(root)
	assert.Equal(t, state.BeaconState(nil), s)
	assert.Equal(t, false, c.has(root), "Empty cache has an object")

	s, err := v1.InitializeFromProto(&ethpb.BeaconState{
		Slot: 10,
	})
	require.NoError(t, err)

	c.put(root, s)
	assert.Equal(t, true, c.has(root), "Empty cache does not have an object")

	res := c.get(root)
	assert.NotNil(t, s)
	assert.DeepEqual(t, res.CloneInnerState(), s.CloneInnerState(), "Expected equal protos to return from cache")

	c.delete(root)
	assert.Equal(t, false, c.has(root), "Cache not supposed to have the object")
}

func TestHotStateSaving_Enabled(t *testing.T) {
	h := &hotStateSaver{
		db: testDB.SetupDB(t),
	}
	h.enableSnapshots()
	require.Equal(t, PersistenceModeSnapshot, h.mode())
}

func TestHotStateSaving_AlreadyEnabled(t *testing.T) {
	h := &hotStateSaver{
		db: testDB.SetupDB(t),
		m:  PersistenceModeSnapshot,
	}
	h.enableSnapshots()
	require.Equal(t, PersistenceModeSnapshot, h.mode())
}

func TestHotStateSaving_Disabled(t *testing.T) {
	ctx := context.Background()
	h := &hotStateSaver{
		db: testDB.SetupDB(t),
		m:  PersistenceModeSnapshot,
	}
	b := util.NewBeaconBlock()
	util.SaveBlock(t, ctx, h.db, b)
	r, err := b.Block.HashTreeRoot()
	require.NoError(t, err)
	h.savedRoots = [][32]byte{r}
	require.NoError(t, h.disableSnapshots(ctx))
	require.Equal(t, PersistenceModeMemory, h.mode())
	require.Equal(t, 0, len(h.savedRoots))
}

func TestHotStateSaving_AlreadyDisabled(t *testing.T) {
	h := &hotStateSaver{}
	require.NoError(t, h.disableSnapshots(context.Background()))
	require.Equal(t, PersistenceModeMemory, h.mode())
}

func TestHotStateSaving_DisabledByDefault(t *testing.T) {
	h := &hotStateSaver{
		db: testDB.SetupDB(t),
		fc: doublylinkedtree.New(),
	}
	fin := h.fc.FinalizedCheckpoint()
	finslot, err := slots.EpochStart(fin.Epoch)
	require.NoError(t, err)
	h.cs = &mockCurrentSlotter{Slot: finslot}
	require.Equal(t, PersistenceModeMemory, h.mode())
	mode, err := h.refreshMode(context.Background())
	require.NoError(t, err)
	require.Equal(t, PersistenceModeMemory, mode)
}

func TestHotStateSaving_Enabling(t *testing.T) {
	h := &hotStateSaver{
		db: testDB.SetupDB(t),
		fc: doublylinkedtree.New(),
		cs: &mockCurrentSlotter{Slot: types.Slot(uint64(params.BeaconConfig().SlotsPerEpoch) * uint64(hotStateSaveThreshold))},
	}
	mode, err := h.refreshMode(context.Background())
	require.NoError(t, err)
	require.Equal(t, PersistenceModeSnapshot, mode)
}

func TestHotStateSaving_DisableAfterFinality(t *testing.T) {
	h := &hotStateSaver{
		db: testDB.SetupDB(t),
		fc: doublylinkedtree.New(),
		cs: &mockCurrentSlotter{Slot: types.Slot(uint64(params.BeaconConfig().SlotsPerEpoch) * uint64(hotStateSaveThreshold))},
	}
	mode, err := h.refreshMode(context.Background())
	require.NoError(t, err)
	require.Equal(t, PersistenceModeSnapshot, mode)

	// set current slot equal to finalized and ask for an update, should be disabled
	fin := h.fc.FinalizedCheckpoint()
	finslot, err := slots.EpochStart(fin.Epoch)
	require.NoError(t, err)
	h.cs = &mockCurrentSlotter{Slot: finslot}
	mode, err = h.refreshMode(context.Background())
	require.NoError(t, err)
	require.Equal(t, PersistenceModeMemory, mode)
}

type mockFinalizedCheckpointer struct {
	c *forkchoicetypes.Checkpoint
}

func (m *mockFinalizedCheckpointer) FinalizedCheckpoint() *forkchoicetypes.Checkpoint {
	return m.c
}

var _ FinalizedCheckpointer = &mockFinalizedCheckpointer{}

func TestUpdateHotStateMode_CurrentSlotBeforeFinalized(t *testing.T) {
	h := &hotStateSaver{
		db: testDB.SetupDB(t),
		fc: &mockFinalizedCheckpointer{c: &forkchoicetypes.Checkpoint{Epoch: 1}},
		cs: &mockCurrentSlotter{Slot: 0},
	}
	mode, err := h.refreshMode(context.Background())
	require.ErrorIs(t, err, errCurrentEpochBehindFinalized)
	require.Equal(t, PersistenceModeMemory, mode)
}

func TestUpdateHotStateMode_NilFinalized(t *testing.T) {
	h := &hotStateSaver{
		db: testDB.SetupDB(t),
		fc: &mockFinalizedCheckpointer{c: nil},
		cs: &mockCurrentSlotter{Slot: 0},
	}
	mode, err := h.refreshMode(context.Background())
	require.ErrorIs(t, err, errForkchoiceFinalizedNil)
	require.Equal(t, PersistenceModeMemory, mode)
}
