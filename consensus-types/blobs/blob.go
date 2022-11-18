package blobs

import (
	"github.com/ethereum/go-ethereum/crypto/kzg"
	types "github.com/prysmaticlabs/prysm/v3/consensus-types/primitives"
	"github.com/prysmaticlabs/prysm/v3/encoding/bytesutil"
	v1 "github.com/prysmaticlabs/prysm/v3/proto/engine/v1"
	eth "github.com/prysmaticlabs/prysm/v3/proto/prysm/v1alpha1"
)

type commitmentSequenceImpl [][]byte

func (s commitmentSequenceImpl) At(i int) kzg.KZGCommitment {
	var out kzg.KZGCommitment
	copy(out[:], s[i])
	return out
}

func (s commitmentSequenceImpl) Len() int {
	return len(s)
}

type BlobImpl []byte

func (b BlobImpl) At(i int) [32]byte {
	var out [32]byte
	copy(out[:], b[i*32:(i+1)*32-1])
	return out
}

func (b BlobImpl) Len() int {
	return len(b)
}

type BlobsSequenceImpl []*v1.Blob

func (s BlobsSequenceImpl) At(i int) kzg.Blob {
	return BlobImpl(s[i].Data)
}

func (s BlobsSequenceImpl) Len() int {
	return len(s)
}

// ValidateBlobsSidecar verifies the integrity of a sidecar, returning nil if the blob is valid.
func ValidateBlobsSidecar(slot types.Slot, root [32]byte, commitments [][]byte, sidecar *eth.BlobsSidecar) error {
	kzgSidecar := kzg.BlobsSidecar{
		BeaconBlockRoot:    kzg.Root(bytesutil.ToBytes32(sidecar.BeaconBlockRoot)),
		BeaconBlockSlot:    kzg.Slot(sidecar.BeaconBlockSlot),
		Blobs:              BlobsSequenceImpl(sidecar.Blobs),
		KZGAggregatedProof: kzg.KZGProof(bytesutil.ToBytes48(sidecar.AggregatedProof)),
	}
	return kzg.ValidateBlobsSidecar(kzg.Slot(slot), kzg.Root(root), commitmentSequenceImpl(commitments), kzgSidecar)
}
