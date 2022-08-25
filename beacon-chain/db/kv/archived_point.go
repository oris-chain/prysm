package kv

import (
	"context"

	types "github.com/prysmaticlabs/prysm/v3/consensus-types/primitives"
	"github.com/prysmaticlabs/prysm/v3/encoding/bytesutil"
	bolt "go.etcd.io/bbolt"
	"go.opencensus.io/trace"
)

// LastArchivedSlot from the db.
func (s *Store) LastArchivedSlot(ctx context.Context) (types.Slot, error) {
	ctx, span := trace.StartSpan(ctx, "BeaconDB.LastArchivedSlot")
	defer span.End()
	var index types.Slot
	err := s.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(stateSlotIndicesBucket)
		b, _ := bkt.Cursor().Last()
		index = bytesutil.BytesToSlotBigEndian(b)
		return nil
	})

	return index, err
}

// ArchivedPointRoot returns the block root of an archived point from the DB.
// This is essential for cold state management and to restore a cold state.
func (s *Store) ArchivedPointRoot(ctx context.Context, slot types.Slot) [32]byte {
	ctx, span := trace.StartSpan(ctx, "BeaconDB.ArchivedPointRoot")
	defer span.End()

	var blockRoot []byte
	if err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(stateSlotIndicesBucket)
		blockRoot = bucket.Get(bytesutil.SlotToBytesBigEndian(slot))
		return nil
	}); err != nil { // This view never returns an error, but we'll handle anyway for sanity.
		panic(err)
	}

	return bytesutil.ToBytes32(blockRoot)
}
