package state_native

import (
	fieldparams "github.com/prysmaticlabs/prysm/v4/config/fieldparams"
	multi_value_slice "github.com/prysmaticlabs/prysm/v4/container/multi-value-slice"
	"github.com/prysmaticlabs/prysm/v4/encoding/bytesutil"
	ethpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
)

// MultiValueRandaoMixes is a multi-value slice of randao mixes.
type MultiValueRandaoMixes = multi_value_slice.Slice[[32]byte, *BeaconState]

// NewMultiValueRandaoMixes creates a new slice whose shared items will be populated with copies of input values.
func NewMultiValueRandaoMixes(mixes [][]byte) *MultiValueRandaoMixes {
	items := make([][32]byte, fieldparams.RandaoMixesLength)
	for i, v := range mixes {
		items[i] = [32]byte(bytesutil.PadTo(v, 32))
	}
	mv := &MultiValueRandaoMixes{}
	mv.Init(items)
	return mv
}

// MultiValueBlockRoots is a multi-value slice of block roots.
type MultiValueBlockRoots = multi_value_slice.Slice[[32]byte, *BeaconState]

// NewMultiValueBlockRoots creates a new slice whose shared items will be populated with copies of input values.
func NewMultiValueBlockRoots(roots [][]byte) *MultiValueBlockRoots {
	items := make([][32]byte, fieldparams.BlockRootsLength)
	for i, v := range roots {
		items[i] = [32]byte(bytesutil.PadTo(v, 32))
	}
	mv := &MultiValueBlockRoots{}
	mv.Init(items)
	return mv
}

// MultiValueStateRoots is a multi-value slice of state roots.
type MultiValueStateRoots = multi_value_slice.Slice[[32]byte, *BeaconState]

// NewMultiValueStateRoots creates a new slice whose shared items will be populated with copies of input values.
func NewMultiValueStateRoots(roots [][]byte) *MultiValueStateRoots {
	items := make([][32]byte, fieldparams.StateRootsLength)
	for i, v := range roots {
		items[i] = [32]byte(bytesutil.PadTo(v, 32))
	}
	mv := &MultiValueStateRoots{}
	mv.Init(items)
	return mv
}

// MultiValueBalances is a multi-value slice of balances.
type MultiValueBalances = multi_value_slice.Slice[uint64, *BeaconState]

// NewMultiValueBalances creates a new slice whose shared items will be populated with copies of input values.
func NewMultiValueBalances(balances []uint64) *MultiValueBalances {
	items := make([]uint64, len(balances))
	copy(items, balances)
	mv := &MultiValueBalances{}
	mv.Init(items)
	return mv
}

// MultiValueInactivityScores is a multi-value slice of inactivity scores.
type MultiValueInactivityScores = multi_value_slice.Slice[uint64, *BeaconState]

// NewMultiValueInactivityScores creates a new slice whose shared items will be populated with copies of input values.
func NewMultiValueInactivityScores(scores []uint64) *MultiValueInactivityScores {
	items := make([]uint64, len(scores))
	copy(items, scores)
	mv := &MultiValueInactivityScores{}
	mv.Init(items)
	return mv
}

// MultiValueValidators is a multi-value slice of validator references.
type MultiValueValidators = multi_value_slice.Slice[*ethpb.Validator, *BeaconState]

// NewMultiValueValidators creates a new slice whose shared items will be populated with input values.
func NewMultiValueValidators(vals []*ethpb.Validator) *MultiValueValidators {
	mv := &MultiValueValidators{}
	mv.Init(vals)
	return mv
}