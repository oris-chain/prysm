package state_native

import (
	"fmt"

	"github.com/pkg/errors"
	customtypes "github.com/prysmaticlabs/prysm/v4/beacon-chain/state/state-native/custom-types"
	"github.com/prysmaticlabs/prysm/v4/config/features"
	ethpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v4/runtime/version"
)

// ToProtoUnsafe returns the pointer value of the underlying
// beacon state proto object, bypassing immutability. Use with care.
func (b *BeaconState) ToProtoUnsafe() interface{} {
	if b == nil {
		return nil
	}

	gvrCopy := b.genesisValidatorsRoot
	br := customtypes.BlockRoots(b.blockRootsVal())
	sr := customtypes.StateRoots(b.stateRootsVal())
	rm := customtypes.RandaoMixes(b.randaoMixesVal())

	switch b.version {
	case version.Phase0:
		return &ethpb.BeaconState{
			GenesisTime:                 b.genesisTime,
			GenesisValidatorsRoot:       gvrCopy[:],
			Slot:                        b.slot,
			Fork:                        b.fork,
			LatestBlockHeader:           b.latestBlockHeader,
			BlockRoots:                  br.Slice(),
			StateRoots:                  sr.Slice(),
			HistoricalRoots:             b.historicalRoots.Slice(),
			Eth1Data:                    b.eth1Data,
			Eth1DataVotes:               b.eth1DataVotes,
			Eth1DepositIndex:            b.eth1DepositIndex,
			Validators:                  b.validatorsVal(),
			Balances:                    b.balancesVal(),
			RandaoMixes:                 rm.Slice(),
			Slashings:                   b.slashings,
			PreviousEpochAttestations:   b.previousEpochAttestations,
			CurrentEpochAttestations:    b.currentEpochAttestations,
			JustificationBits:           b.justificationBits,
			PreviousJustifiedCheckpoint: b.previousJustifiedCheckpoint,
			CurrentJustifiedCheckpoint:  b.currentJustifiedCheckpoint,
			FinalizedCheckpoint:         b.finalizedCheckpoint,
		}
	case version.Altair:
		return &ethpb.BeaconStateAltair{
			GenesisTime:                 b.genesisTime,
			GenesisValidatorsRoot:       gvrCopy[:],
			Slot:                        b.slot,
			Fork:                        b.fork,
			LatestBlockHeader:           b.latestBlockHeader,
			BlockRoots:                  br.Slice(),
			StateRoots:                  sr.Slice(),
			HistoricalRoots:             b.historicalRoots.Slice(),
			Eth1Data:                    b.eth1Data,
			Eth1DataVotes:               b.eth1DataVotes,
			Eth1DepositIndex:            b.eth1DepositIndex,
			Validators:                  b.validatorsVal(),
			Balances:                    b.balancesVal(),
			RandaoMixes:                 rm.Slice(),
			Slashings:                   b.slashings,
			PreviousEpochParticipation:  b.previousEpochParticipation,
			CurrentEpochParticipation:   b.currentEpochParticipation,
			JustificationBits:           b.justificationBits,
			PreviousJustifiedCheckpoint: b.previousJustifiedCheckpoint,
			CurrentJustifiedCheckpoint:  b.currentJustifiedCheckpoint,
			FinalizedCheckpoint:         b.finalizedCheckpoint,
			InactivityScores:            b.inactivityScoresVal(),
			CurrentSyncCommittee:        b.currentSyncCommittee,
			NextSyncCommittee:           b.nextSyncCommittee,
		}
	case version.Bellatrix:
		return &ethpb.BeaconStateBellatrix{
			GenesisTime:                  b.genesisTime,
			GenesisValidatorsRoot:        gvrCopy[:],
			Slot:                         b.slot,
			Fork:                         b.fork,
			LatestBlockHeader:            b.latestBlockHeader,
			BlockRoots:                   br.Slice(),
			StateRoots:                   sr.Slice(),
			HistoricalRoots:              b.historicalRoots.Slice(),
			Eth1Data:                     b.eth1Data,
			Eth1DataVotes:                b.eth1DataVotes,
			Eth1DepositIndex:             b.eth1DepositIndex,
			Validators:                   b.validatorsVal(),
			Balances:                     b.balancesVal(),
			RandaoMixes:                  rm.Slice(),
			Slashings:                    b.slashings,
			PreviousEpochParticipation:   b.previousEpochParticipation,
			CurrentEpochParticipation:    b.currentEpochParticipation,
			JustificationBits:            b.justificationBits,
			PreviousJustifiedCheckpoint:  b.previousJustifiedCheckpoint,
			CurrentJustifiedCheckpoint:   b.currentJustifiedCheckpoint,
			FinalizedCheckpoint:          b.finalizedCheckpoint,
			InactivityScores:             b.inactivityScoresVal(),
			CurrentSyncCommittee:         b.currentSyncCommittee,
			NextSyncCommittee:            b.nextSyncCommittee,
			LatestExecutionPayloadHeader: b.latestExecutionPayloadHeader,
		}
	case version.Capella:
		return &ethpb.BeaconStateCapella{
			GenesisTime:                  b.genesisTime,
			GenesisValidatorsRoot:        gvrCopy[:],
			Slot:                         b.slot,
			Fork:                         b.fork,
			LatestBlockHeader:            b.latestBlockHeader,
			BlockRoots:                   br.Slice(),
			StateRoots:                   sr.Slice(),
			HistoricalRoots:              b.historicalRoots.Slice(),
			Eth1Data:                     b.eth1Data,
			Eth1DataVotes:                b.eth1DataVotes,
			Eth1DepositIndex:             b.eth1DepositIndex,
			Validators:                   b.validatorsVal(),
			Balances:                     b.balancesVal(),
			RandaoMixes:                  rm.Slice(),
			Slashings:                    b.slashings,
			PreviousEpochParticipation:   b.previousEpochParticipation,
			CurrentEpochParticipation:    b.currentEpochParticipation,
			JustificationBits:            b.justificationBits,
			PreviousJustifiedCheckpoint:  b.previousJustifiedCheckpoint,
			CurrentJustifiedCheckpoint:   b.currentJustifiedCheckpoint,
			FinalizedCheckpoint:          b.finalizedCheckpoint,
			InactivityScores:             b.inactivityScoresVal(),
			CurrentSyncCommittee:         b.currentSyncCommittee,
			NextSyncCommittee:            b.nextSyncCommittee,
			LatestExecutionPayloadHeader: b.latestExecutionPayloadHeaderCapella,
			NextWithdrawalIndex:          b.nextWithdrawalIndex,
			NextWithdrawalValidatorIndex: b.nextWithdrawalValidatorIndex,
			HistoricalSummaries:          b.historicalSummaries,
		}
	default:
		return nil
	}
}

// ToProto the beacon state into a protobuf for usage.
func (b *BeaconState) ToProto() interface{} {
	if b == nil {
		return nil
	}

	b.lock.RLock()
	defer b.lock.RUnlock()

	gvrCopy := b.genesisValidatorsRoot

	br := customtypes.BlockRoots(b.blockRootsVal())
	brSlice := br.Slice()
	brCopy := make([][]byte, len(br))
	for i, v := range brSlice {
		brCopy[i] = make([]byte, len(v))
		copy(brCopy[i], v)
	}

	sr := customtypes.StateRoots(b.stateRootsVal())
	srSlice := sr.Slice()
	srCopy := make([][]byte, len(sr))
	for i, v := range srSlice {
		srCopy[i] = make([]byte, len(v))
		copy(srCopy[i], v)
	}

	rm := customtypes.RandaoMixes(b.randaoMixesVal())
	rmSlice := rm.Slice()
	rmCopy := make([][]byte, len(rm))
	for i, v := range rmSlice {
		rmCopy[i] = make([]byte, len(v))
		copy(rmCopy[i], v)
	}

	balances := b.balancesVal()
	balancesCopy := make([]uint64, len(balances))
	copy(balancesCopy, balances)

	var inactivityScoresCopy []uint64
	if b.version > version.Phase0 {
		inactivityScores := b.inactivityScoresVal()
		inactivityScoresCopy = make([]uint64, len(inactivityScores))
		copy(inactivityScoresCopy, inactivityScores)
	}

	vals := b.validatorsVal()
	valsCopy := make([]*ethpb.Validator, len(vals))
	for i := 0; i < len(vals); i++ {
		val := vals[i]
		if val == nil {
			continue
		}
		valsCopy[i] = ethpb.CopyValidator(val)
	}

	switch b.version {
	case version.Phase0:
		return &ethpb.BeaconState{
			GenesisTime:                 b.genesisTime,
			GenesisValidatorsRoot:       gvrCopy[:],
			Slot:                        b.slot,
			Fork:                        b.forkVal(),
			LatestBlockHeader:           b.latestBlockHeaderVal(),
			BlockRoots:                  brCopy,
			StateRoots:                  srCopy,
			HistoricalRoots:             b.historicalRoots.Slice(),
			Eth1Data:                    b.eth1DataVal(),
			Eth1DataVotes:               b.eth1DataVotesVal(),
			Eth1DepositIndex:            b.eth1DepositIndex,
			Validators:                  valsCopy,
			Balances:                    balancesCopy,
			RandaoMixes:                 rmCopy,
			Slashings:                   b.slashingsVal(),
			PreviousEpochAttestations:   b.previousEpochAttestationsVal(),
			CurrentEpochAttestations:    b.currentEpochAttestationsVal(),
			JustificationBits:           b.justificationBitsVal(),
			PreviousJustifiedCheckpoint: b.previousJustifiedCheckpointVal(),
			CurrentJustifiedCheckpoint:  b.currentJustifiedCheckpointVal(),
			FinalizedCheckpoint:         b.finalizedCheckpointVal(),
		}
	case version.Altair:
		return &ethpb.BeaconStateAltair{
			GenesisTime:                 b.genesisTime,
			GenesisValidatorsRoot:       gvrCopy[:],
			Slot:                        b.slot,
			Fork:                        b.forkVal(),
			LatestBlockHeader:           b.latestBlockHeaderVal(),
			BlockRoots:                  brCopy,
			StateRoots:                  srCopy,
			HistoricalRoots:             b.historicalRoots.Slice(),
			Eth1Data:                    b.eth1DataVal(),
			Eth1DataVotes:               b.eth1DataVotesVal(),
			Eth1DepositIndex:            b.eth1DepositIndex,
			Validators:                  valsCopy,
			Balances:                    balancesCopy,
			RandaoMixes:                 rmCopy,
			Slashings:                   b.slashingsVal(),
			PreviousEpochParticipation:  b.previousEpochParticipationVal(),
			CurrentEpochParticipation:   b.currentEpochParticipationVal(),
			JustificationBits:           b.justificationBitsVal(),
			PreviousJustifiedCheckpoint: b.previousJustifiedCheckpointVal(),
			CurrentJustifiedCheckpoint:  b.currentJustifiedCheckpointVal(),
			FinalizedCheckpoint:         b.finalizedCheckpointVal(),
			InactivityScores:            inactivityScoresCopy,
			CurrentSyncCommittee:        b.currentSyncCommitteeVal(),
			NextSyncCommittee:           b.nextSyncCommitteeVal(),
		}
	case version.Bellatrix:
		return &ethpb.BeaconStateBellatrix{
			GenesisTime:                  b.genesisTime,
			GenesisValidatorsRoot:        gvrCopy[:],
			Slot:                         b.slot,
			Fork:                         b.forkVal(),
			LatestBlockHeader:            b.latestBlockHeaderVal(),
			BlockRoots:                   brCopy,
			StateRoots:                   srCopy,
			HistoricalRoots:              b.historicalRoots.Slice(),
			Eth1Data:                     b.eth1DataVal(),
			Eth1DataVotes:                b.eth1DataVotesVal(),
			Eth1DepositIndex:             b.eth1DepositIndex,
			Validators:                   valsCopy,
			Balances:                     balancesCopy,
			RandaoMixes:                  rmCopy,
			Slashings:                    b.slashingsVal(),
			PreviousEpochParticipation:   b.previousEpochParticipationVal(),
			CurrentEpochParticipation:    b.currentEpochParticipationVal(),
			JustificationBits:            b.justificationBitsVal(),
			PreviousJustifiedCheckpoint:  b.previousJustifiedCheckpointVal(),
			CurrentJustifiedCheckpoint:   b.currentJustifiedCheckpointVal(),
			FinalizedCheckpoint:          b.finalizedCheckpointVal(),
			InactivityScores:             inactivityScoresCopy,
			CurrentSyncCommittee:         b.currentSyncCommitteeVal(),
			NextSyncCommittee:            b.nextSyncCommitteeVal(),
			LatestExecutionPayloadHeader: b.latestExecutionPayloadHeaderVal(),
		}
	case version.Capella:
		return &ethpb.BeaconStateCapella{
			GenesisTime:                  b.genesisTime,
			GenesisValidatorsRoot:        gvrCopy[:],
			Slot:                         b.slot,
			Fork:                         b.forkVal(),
			LatestBlockHeader:            b.latestBlockHeaderVal(),
			BlockRoots:                   brCopy,
			StateRoots:                   srCopy,
			HistoricalRoots:              b.historicalRoots.Slice(),
			Eth1Data:                     b.eth1DataVal(),
			Eth1DataVotes:                b.eth1DataVotesVal(),
			Eth1DepositIndex:             b.eth1DepositIndex,
			Validators:                   valsCopy,
			Balances:                     balancesCopy,
			RandaoMixes:                  rmCopy,
			Slashings:                    b.slashingsVal(),
			PreviousEpochParticipation:   b.previousEpochParticipationVal(),
			CurrentEpochParticipation:    b.currentEpochParticipationVal(),
			JustificationBits:            b.justificationBitsVal(),
			PreviousJustifiedCheckpoint:  b.previousJustifiedCheckpointVal(),
			CurrentJustifiedCheckpoint:   b.currentJustifiedCheckpointVal(),
			FinalizedCheckpoint:          b.finalizedCheckpointVal(),
			InactivityScores:             inactivityScoresCopy,
			CurrentSyncCommittee:         b.currentSyncCommitteeVal(),
			NextSyncCommittee:            b.nextSyncCommitteeVal(),
			LatestExecutionPayloadHeader: b.latestExecutionPayloadHeaderCapellaVal(),
			NextWithdrawalIndex:          b.nextWithdrawalIndex,
			NextWithdrawalValidatorIndex: b.nextWithdrawalValidatorIndex,
			HistoricalSummaries:          b.historicalSummariesVal(),
		}
	default:
		return nil
	}
}

// StateRoots kept track of in the beacon state.
func (b *BeaconState) StateRoots() [][]byte {
	b.lock.RLock()
	defer b.lock.RUnlock()

	roots := b.stateRootsVal()
	rootsCopy := make([][]byte, len(roots))
	for i, r := range roots {
		rootsCopy[i] = make([]byte, 32)
		copy(rootsCopy[i], r[:])
	}
	return rootsCopy
}

func (b *BeaconState) stateRootsVal() [][32]byte {
	if features.Get().EnableExperimentalState {
		if b.stateRootsMultiValue == nil {
			return nil
		}
		return b.stateRootsMultiValue.Value(b)
	}
	if b.stateRoots == nil {
		return nil
	}
	return b.stateRoots
}

// StateRootAtIndex retrieves a specific state root based on an
// input index value.
func (b *BeaconState) StateRootAtIndex(idx uint64) ([]byte, error) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	if features.Get().EnableExperimentalState {
		if b.stateRootsMultiValue == nil {
			return []byte{}, nil
		}
		r, err := b.stateRootsMultiValue.At(b, idx)
		if err != nil {
			return nil, err
		}
		return r[:], nil
	}

	if b.stateRoots == nil {
		return []byte{}, nil
	}
	if uint64(len(b.stateRoots)) <= idx {
		return []byte{}, fmt.Errorf("index %d out of bounds", idx)
	}
	return b.stateRoots[idx][:], nil
}

// ProtobufBeaconStatePhase0 transforms an input into beacon state in the form of protobuf.
// Error is returned if the input is not type protobuf beacon state.
func ProtobufBeaconStatePhase0(s interface{}) (*ethpb.BeaconState, error) {
	pbState, ok := s.(*ethpb.BeaconState)
	if !ok {
		return nil, errors.New("input is not type ethpb.BeaconState")
	}
	return pbState, nil
}

// ProtobufBeaconStateAltair transforms an input into beacon state Altair in the form of protobuf.
// Error is returned if the input is not type protobuf beacon state.
func ProtobufBeaconStateAltair(s interface{}) (*ethpb.BeaconStateAltair, error) {
	pbState, ok := s.(*ethpb.BeaconStateAltair)
	if !ok {
		return nil, errors.New("input is not type pb.BeaconStateAltair")
	}
	return pbState, nil
}

// ProtobufBeaconStateBellatrix transforms an input into beacon state Bellatrix in the form of protobuf.
// Error is returned if the input is not type protobuf beacon state.
func ProtobufBeaconStateBellatrix(s interface{}) (*ethpb.BeaconStateBellatrix, error) {
	pbState, ok := s.(*ethpb.BeaconStateBellatrix)
	if !ok {
		return nil, errors.New("input is not type pb.BeaconStateBellatrix")
	}
	return pbState, nil
}

// ProtobufBeaconStateCapella transforms an input into beacon state Capella in the form of protobuf.
// Error is returned if the input is not type protobuf beacon state.
func ProtobufBeaconStateCapella(s interface{}) (*ethpb.BeaconStateCapella, error) {
	pbState, ok := s.(*ethpb.BeaconStateCapella)
	if !ok {
		return nil, errors.New("input is not type pb.BeaconStateCapella")
	}
	return pbState, nil
}
