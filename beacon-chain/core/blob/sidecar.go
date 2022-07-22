package blob

import (
	types2 "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/blocks"
	"github.com/prysmaticlabs/prysm/consensus-types/interfaces"
	types "github.com/prysmaticlabs/prysm/consensus-types/primitives"
	"github.com/prysmaticlabs/prysm/encoding/bytesutil"
	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
)

var (
	ErrInvalidBlobSlot            = errors.New("invalid blob slot")
	ErrInvalidBlobBeaconBlockRoot = errors.New("invalid blob beacon block root")
	ErrInvalidBlobsLength         = errors.New("invalid blobs length")
	ErrCouldNotComputeCommitment  = errors.New("could not compute commitment")
	ErrMissmatchKzgs              = errors.New("missmatch kzgs")
)

// VerifyBlobsSidecar verifies the integrity of a sidecar.
// def verify_blobs_sidecar(slot: Slot, beacon_block_root: Root,
//                         expected_kzgs: Sequence[KZGCommitment], blobs_sidecar: BlobsSidecar):
//    assert slot == blobs_sidecar.beacon_block_slot
//    assert beacon_block_root == blobs_sidecar.beacon_block_root
//    blobs = blobs_sidecar.blobs
//    assert len(expected_kzgs) == len(blobs)
//    for kzg, blob in zip(expected_kzgs, blobs):
//        assert blob_to_kzg(blob) == kzg
func VerifyBlobsSidecar(slot types.Slot, beaconBlockRoot [32]byte, expectedKZGs [][48]byte, blobsSidecar *eth.BlobsSidecar) error {
	// TODO(EIP-4844): Apply optimization - https://github.com/ethereum/consensus-specs/blob/0ba5b3b5c5bb58fbe0f094dcd02dedc4ff1c6f7c/specs/eip4844/validator.md#verify_blobs_sidecar

	if slot != blobsSidecar.BeaconBlockSlot {
		return ErrInvalidBlobSlot
	}
	if beaconBlockRoot != bytesutil.ToBytes32(blobsSidecar.BeaconBlockRoot) {
		return ErrInvalidBlobBeaconBlockRoot
	}
	if len(expectedKZGs) != len(blobsSidecar.Blobs) {
		return ErrInvalidBlobsLength
	}
	for i, expectedKzg := range expectedKZGs {
		var blob types2.Blob
		for i, b := range blobsSidecar.Blobs[i].Blob {
			var f types2.BLSFieldElement
			copy(f[:], b)
			blob[i] = f
		}
		kzg, ok := blob.ComputeCommitment()
		if !ok {
			return ErrCouldNotComputeCommitment
		}
		if kzg != expectedKzg {
			return ErrMissmatchKzgs
		}
	}
	return nil
}

func BlockContainsKZGs(b interfaces.BeaconBlock) bool {
	if blocks.IsPreEIP4844Version(b.Version()) {
		return false
	}
	blobKzgs, err := b.Body().BlobKzgs()
	if err != nil {
		// cannot happen!
		return false
	}
	return len(blobKzgs) != 0
}
