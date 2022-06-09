package state_native

import (
	ethpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
)

// LatestExecutionPayloadHeader of the beacon state.
func (b *BeaconState) LatestExecutionPayloadHeader() (*ethpb.ExecutionPayloadHeader, error) {
	if !b.version.IsExecutionPayloadCompatible() {
		return nil, errNotSupported("LatestExecutionPayloadHeader", b.version)
	}

	if b.latestExecutionPayloadHeader == nil {
		return nil, nil
	}

	b.lock.RLock()
	defer b.lock.RUnlock()

	return b.latestExecutionPayloadHeaderVal(), nil
}

// latestExecutionPayloadHeaderVal of the beacon state.
// This assumes that a lock is already held on BeaconState.
func (b *BeaconState) latestExecutionPayloadHeaderVal() *ethpb.ExecutionPayloadHeader {
	return ethpb.CopyExecutionPayloadHeader(b.latestExecutionPayloadHeader)
}
