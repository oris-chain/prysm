package validator

import (
	"sync"

	"github.com/prysmaticlabs/prysm/v4/consensus-types/blocks"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/interfaces"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
	ethpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v4/runtime/version"
)

var bundleCache = &blobsBundleCache{}

// BlobsBundleCache holds the KZG commitments and other relevant sidecar data for a local beacon block.
type blobsBundleCache struct {
	sync.Mutex
	slot   primitives.Slot
	bundle *enginev1.BlobsBundle
}

// add adds a blobs bundle to the cache.
// same slot overwrites the previous bundle.
func (c *blobsBundleCache) add(slot primitives.Slot, bundle *enginev1.BlobsBundle) {
	c.Lock()
	defer c.Unlock()

	if slot >= c.slot {
		c.bundle = bundle
		c.slot = slot
	}
}

// get gets a blobs bundle from the cache.
func (c *blobsBundleCache) get(slot primitives.Slot) *enginev1.BlobsBundle {
	c.Lock()
	defer c.Unlock()

	if c.slot == slot {
		return c.bundle
	}

	return nil
}

// prune acquires the lock before pruning.
func (c *blobsBundleCache) prune(minSlot primitives.Slot) {
	c.Lock()
	defer c.Unlock()

	if minSlot > c.slot {
		c.slot = 0
		c.bundle = nil
	}
}

// buildBlobSidecars given a block, builds the blob sidecars for the block.
func buildBlobSidecars(blk interfaces.SignedBeaconBlock) ([]*ethpb.BlobSidecar, error) {
	if blk.Version() < version.Deneb {
		return nil, nil // No blobs before deneb.
	}
	bundle := bundleCache.get(blk.Block().Slot())
	if bundle == nil {
		return nil, nil
	}
	blobSidecars := make([]*ethpb.BlobSidecar, len(bundle.KzgCommitments))
	header, err := blk.Header()
	if err != nil {
		return nil, err
	}
	body := blk.Block().Body()
	for i := range blobSidecars {
		proof, err := blocks.MerkleProofKZGCommitment(body, i)
		if err != nil {
			return nil, err
		}
		blobSidecars[i] = &ethpb.BlobSidecar{
			Index:                    uint64(i),
			Blob:                     bundle.Blobs[i],
			KzgCommitment:            bundle.KzgCommitments[i],
			KzgProof:                 bundle.Proofs[i],
			SignedBlockHeader:        header,
			CommitmentInclusionProof: proof,
		}
	}
	return blobSidecars, nil
}

// coverts a blinds blobs bundle to a sidecar format.
func blindBlobsBundleToSidecars(bundle *enginev1.BlindedBlobsBundle, blk interfaces.ReadOnlyBeaconBlock) ([]*ethpb.BlindedBlobSidecar, error) {
	if blk.Version() < version.Deneb {
		return nil, nil
	}
	if bundle == nil || len(bundle.KzgCommitments) == 0 {
		return nil, nil
	}
	r, err := blk.HashTreeRoot()
	if err != nil {
		return nil, err
	}
	pr := blk.ParentRoot()

	sidecars := make([]*ethpb.BlindedBlobSidecar, len(bundle.BlobRoots))
	for i := 0; i < len(bundle.BlobRoots); i++ {
		sidecars[i] = &ethpb.BlindedBlobSidecar{
			BlockRoot:       r[:],
			Index:           uint64(i),
			Slot:            blk.Slot(),
			BlockParentRoot: pr[:],
			ProposerIndex:   blk.ProposerIndex(),
			BlobRoot:        bundle.BlobRoots[i],
			KzgCommitment:   bundle.KzgCommitments[i],
			KzgProof:        bundle.Proofs[i],
		}
	}

	return sidecars, nil
}
