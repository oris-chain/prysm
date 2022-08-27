package stategen

import (
	"context"
	"testing"

	testDB "github.com/prysmaticlabs/prysm/v3/beacon-chain/db/testing"
	forkchoicetypes "github.com/prysmaticlabs/prysm/v3/beacon-chain/forkchoice/types"
	"github.com/prysmaticlabs/prysm/v3/config/params"
	types "github.com/prysmaticlabs/prysm/v3/consensus-types/primitives"
	"github.com/prysmaticlabs/prysm/v3/testing/assert"
	"github.com/prysmaticlabs/prysm/v3/testing/require"
	"github.com/prysmaticlabs/prysm/v3/testing/util"
	"github.com/prysmaticlabs/prysm/v3/time/slots"
	logTest "github.com/sirupsen/logrus/hooks/test"
)

func TestSaveState_HotStateCanBeSaved(t *testing.T) {
	ctx := context.Background()
	beaconDB := testDB.SetupDB(t)
	genroot := [32]byte{}
	cp := &forkchoicetypes.Checkpoint{Epoch: 0, Root: genroot}
	h := newTestSaver(beaconDB, withFinalizedCheckpointer(&mockFinalizedCheckpointer{c: cp}))
	stateSlot := firstSaveableSlotAfter(t, h)
	h.cs = &mockCurrentSlotter{Slot: stateSlot + h.snapshotInterval}
	mode, err := h.refreshMode(ctx)
	require.NoError(t, err)
	require.Equal(t, PersistenceModeSnapshot, mode)
	service := New(beaconDB, h)

	service.slotsPerArchivedPoint = 1
	beaconState, _ := util.DeterministicGenesisState(t, 32)
	// This goes to hot section, verify it can save on epoch boundary.
	require.NoError(t, beaconState.SetSlot(params.BeaconConfig().SlotsPerEpoch))

	r := [32]byte{'a'}
	require.NoError(t, service.SaveState(ctx, r, beaconState))

	// Should save both state and state summary.
	_, ok, err := service.epochBoundaryStateCache.getByBlockRoot(r)
	require.NoError(t, err)
	assert.Equal(t, true, ok, "Should have saved the state")
	assert.Equal(t, true, service.beaconDB.HasStateSummary(ctx, r), "Should have saved the state summary")
}

func TestSaveState_HotStateCached(t *testing.T) {
	hook := logTest.NewGlobal()
	ctx := context.Background()
	beaconDB := testDB.SetupDB(t)

	service := New(beaconDB, newTestSaver(beaconDB))
	service.slotsPerArchivedPoint = 1
	beaconState, _ := util.DeterministicGenesisState(t, 32)
	require.NoError(t, beaconState.SetSlot(params.BeaconConfig().SlotsPerEpoch))

	// Cache the state prior.
	r := [32]byte{'a'}
	service.hotStateCache.put(r, beaconState)
	require.NoError(t, service.SaveState(ctx, r, beaconState))

	// Should not save the state and state summary.
	assert.Equal(t, false, service.beaconDB.HasState(ctx, r), "Should not have saved the state")
	assert.Equal(t, false, service.beaconDB.HasStateSummary(ctx, r), "Should have saved the state summary")
	require.LogsDoNotContain(t, hook, "Saved full state on epoch boundary")
}

func TestState_ForceCheckpoint_SavesStateToDatabase(t *testing.T) {
	ctx := context.Background()
	beaconDB := testDB.SetupDB(t)

	svc := New(beaconDB, newTestSaver(beaconDB))
	beaconState, _ := util.DeterministicGenesisState(t, 32)
	require.NoError(t, beaconState.SetSlot(params.BeaconConfig().SlotsPerEpoch))

	r := [32]byte{'a'}
	svc.hotStateCache.put(r, beaconState)

	require.Equal(t, false, beaconDB.HasState(ctx, r), "Database has state stored already")
	assert.NoError(t, svc.ForceCheckpoint(ctx, r[:]))
	assert.Equal(t, true, beaconDB.HasState(ctx, r), "Did not save checkpoint to database")

	// Should not panic with genesis finalized root.
	assert.NoError(t, svc.ForceCheckpoint(ctx, params.BeaconConfig().ZeroHash[:]))
}

func TestSaveState_Alreadyhas(t *testing.T) {
	hook := logTest.NewGlobal()
	ctx := context.Background()
	beaconDB := testDB.SetupDB(t)
	service := New(beaconDB, newTestSaver(beaconDB))

	beaconState, _ := util.DeterministicGenesisState(t, 32)
	require.NoError(t, beaconState.SetSlot(params.BeaconConfig().SlotsPerEpoch))
	r := [32]byte{'A'}

	// Pre cache the hot state.
	service.hotStateCache.put(r, beaconState)
	require.NoError(t, service.saveStateByRoot(ctx, r, beaconState))

	// Should not save the state and state summary.
	assert.Equal(t, false, service.beaconDB.HasState(ctx, r), "Should not have saved the state")
	assert.Equal(t, false, service.beaconDB.HasStateSummary(ctx, r), "Should have saved the state summary")
	require.LogsDoNotContain(t, hook, "Saved full state on epoch boundary")
}

func TestSaveState_CanSaveOnEpochBoundary(t *testing.T) {
	ctx := context.Background()
	beaconDB := testDB.SetupDB(t)
	genroot := [32]byte{}
	cp := &forkchoicetypes.Checkpoint{Epoch: 0, Root: genroot}
	h := newTestSaver(beaconDB, withFinalizedCheckpointer(&mockFinalizedCheckpointer{c: cp}))
	stateSlot := firstSaveableSlotAfter(t, h)
	h.cs = &mockCurrentSlotter{Slot: stateSlot + h.snapshotInterval}
	mode, err := h.refreshMode(ctx)
	require.NoError(t, err)
	require.Equal(t, PersistenceModeSnapshot, mode)
	service := New(beaconDB, h)

	beaconState, _ := util.DeterministicGenesisState(t, 32)
	require.NoError(t, beaconState.SetSlot(params.BeaconConfig().SlotsPerEpoch))
	r := [32]byte{'A'}

	require.NoError(t, service.saveStateByRoot(ctx, r, beaconState))

	// Should save both state and state summary.
	_, ok, err := service.epochBoundaryStateCache.getByBlockRoot(r)
	require.NoError(t, err)
	require.Equal(t, true, ok, "Did not save epoch boundary state")
	assert.Equal(t, true, service.beaconDB.HasStateSummary(ctx, r), "Should have saved the state summary")
	// Should have not been saved in DB.
	require.Equal(t, false, beaconDB.HasState(ctx, r))
}

func TestSaveState_NoSaveNotEpochBoundary(t *testing.T) {
	hook := logTest.NewGlobal()
	ctx := context.Background()
	beaconDB := testDB.SetupDB(t)
	genroot := [32]byte{}
	cp := &forkchoicetypes.Checkpoint{Epoch: 0, Root: genroot}
	h := newTestSaver(beaconDB, withFinalizedCheckpointer(&mockFinalizedCheckpointer{c: cp}))
	stateSlot := firstSaveableSlotAfter(t, h)
	h.cs = &mockCurrentSlotter{Slot: stateSlot + h.snapshotInterval}
	mode, err := h.refreshMode(ctx)
	require.NoError(t, err)
	require.Equal(t, PersistenceModeSnapshot, mode)
	service := New(beaconDB, h)

	beaconState, _ := util.DeterministicGenesisState(t, 32)
	require.NoError(t, beaconState.SetSlot(params.BeaconConfig().SlotsPerEpoch-1))
	r := [32]byte{'A'}
	b := util.NewBeaconBlock()
	util.SaveBlock(t, ctx, beaconDB, b)
	gRoot, err := b.Block.HashTreeRoot()
	require.NoError(t, err)
	require.NoError(t, beaconDB.SaveGenesisBlockRoot(ctx, gRoot))
	require.NoError(t, service.SaveState(ctx, r, beaconState))

	// Should only save state summary.
	assert.Equal(t, false, service.beaconDB.HasState(ctx, r), "Should not have saved the state")
	assert.Equal(t, true, service.beaconDB.HasStateSummary(ctx, r), "Should have saved the state summary")
	require.LogsDoNotContain(t, hook, "Saved full state on epoch boundary")
	// Should have not been saved in DB.
	require.Equal(t, false, beaconDB.HasState(ctx, r))

	_, ok, err := service.epochBoundaryStateCache.getByBlockRoot(r)
	require.NoError(t, err)
	require.Equal(t, false, ok, "saved to epoch boundary cache in error")
}

func firstSaveableSlotAfter(t *testing.T, h *hotStateSaver) types.Slot {
	min, err := slots.EpochStart(hotStateSaveThreshold)
	require.NoError(t, err)
	f := h.fc.FinalizedCheckpoint()
	require.NotNil(t, f)
	fslot, err := slots.EpochStart(f.Epoch)
	require.NoError(t, err)
	min += fslot
	diff := h.snapshotInterval - (min % h.snapshotInterval)
	aligned := min + diff
	require.Equal(t, types.Slot(0), aligned%h.snapshotInterval)
	return min + diff
}

func TestSaveState_CanSaveHotStateToDB(t *testing.T) {
	ctx := context.Background()
	beaconDB := testDB.SetupDB(t)
	genroot := [32]byte{}
	cp := &forkchoicetypes.Checkpoint{Epoch: 0, Root: genroot}
	h := newTestSaver(beaconDB, withFinalizedCheckpointer(&mockFinalizedCheckpointer{c: cp}))
	stateSlot := firstSaveableSlotAfter(t, h)
	h.cs = &mockCurrentSlotter{Slot: stateSlot + h.snapshotInterval}
	mode, err := h.refreshMode(ctx)
	require.NoError(t, err)
	require.Equal(t, PersistenceModeSnapshot, mode)
	service := New(beaconDB, h)
	beaconState, _ := util.DeterministicGenesisState(t, 32)
	require.NoError(t, beaconState.SetSlot(stateSlot))

	r := [32]byte{'A'}
	require.NoError(t, service.saveStateByRoot(ctx, r, beaconState))

	// Should have saved in DB.
	require.Equal(t, true, h.db.HasState(ctx, r))
}
