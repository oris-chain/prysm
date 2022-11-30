package capella

import (
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/core/time"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/state"
	state_native "github.com/prysmaticlabs/prysm/v3/beacon-chain/state/state-native"
	"github.com/prysmaticlabs/prysm/v3/config/params"
	"github.com/prysmaticlabs/prysm/v3/consensus-types/blocks"
	enginev1 "github.com/prysmaticlabs/prysm/v3/proto/engine/v1"
	ethpb "github.com/prysmaticlabs/prysm/v3/proto/prysm/v1alpha1"
)

// UpgradeToEip4844 updates inputs a generic state to return the version Eip4844 state.
func UpgradeToEip4844(state state.BeaconState) (state.BeaconState, error) {
	if err := state.SetFork(&ethpb.Fork{
		PreviousVersion: state.Fork().CurrentVersion,
		CurrentVersion:  params.BeaconConfig().EIP4844ForkVersion,
		Epoch:           time.CurrentEpoch(state),
	}); err != nil {
		return nil, err
	}

	header, err := state.LatestExecutionPayloadHeader()
	if err != nil {
		return nil, err
	}
	txRoot, err := header.TransactionsRoot()
	if err != nil {
		return nil, err
	}
	withdrawalRoot, err := header.WithdrawalsRoot()
	if err != nil {
		return nil, err
	}
	h, err := blocks.WrappedExecutionPayloadHeaderEIP4844(&enginev1.ExecutionPayloadHeader4844{
		ParentHash:       header.ParentHash(),
		FeeRecipient:     header.FeeRecipient(),
		StateRoot:        header.StateRoot(),
		ReceiptsRoot:     header.ReceiptsRoot(),
		LogsBloom:        header.LogsBloom(),
		PrevRandao:       header.PrevRandao(),
		BlockNumber:      header.BlockNumber(),
		GasLimit:         header.GasLimit(),
		GasUsed:          header.GasUsed(),
		Timestamp:        header.Timestamp(),
		ExtraData:        header.ExtraData(),
		BaseFeePerGas:    header.BaseFeePerGas(),
		BlockHash:        header.BlockHash(),
		TransactionsRoot: txRoot,
		WithdrawalsRoot:  withdrawalRoot,
		ExcessDataGas:    make([]byte, 32),
	})
	if err != nil {
		return nil, err
	}
	err = state.SetLatestExecutionPayloadHeader(h)
	if err != nil {
		return nil, err
	}

	return state, nil
}

// UpgradeToCapella updates a generic state to return the version Capella state.
func UpgradeToCapella(state state.BeaconState) (state.BeaconState, error) {
	epoch := time.CurrentEpoch(state)

	currentSyncCommittee, err := state.CurrentSyncCommittee()
	if err != nil {
		return nil, err
	}
	nextSyncCommittee, err := state.NextSyncCommittee()
	if err != nil {
		return nil, err
	}
	prevEpochParticipation, err := state.PreviousEpochParticipation()
	if err != nil {
		return nil, err
	}
	currentEpochParticipation, err := state.CurrentEpochParticipation()
	if err != nil {
		return nil, err
	}
	inactivityScores, err := state.InactivityScores()
	if err != nil {
		return nil, err
	}
	payloadHeader, err := state.LatestExecutionPayloadHeader()
	if err != nil {
		return nil, err
	}
	txRoot, err := payloadHeader.TransactionsRoot()
	if err != nil {
		return nil, err
	}

	s := &ethpb.BeaconStateCapella{
		GenesisTime:           state.GenesisTime(),
		GenesisValidatorsRoot: state.GenesisValidatorsRoot(),
		Slot:                  state.Slot(),
		Fork: &ethpb.Fork{
			PreviousVersion: state.Fork().CurrentVersion,
			CurrentVersion:  params.BeaconConfig().CapellaForkVersion,
			Epoch:           epoch,
		},
		LatestBlockHeader:           state.LatestBlockHeader(),
		BlockRoots:                  state.BlockRoots(),
		StateRoots:                  state.StateRoots(),
		HistoricalRoots:             state.HistoricalRoots(),
		Eth1Data:                    state.Eth1Data(),
		Eth1DataVotes:               state.Eth1DataVotes(),
		Eth1DepositIndex:            state.Eth1DepositIndex(),
		Validators:                  state.Validators(),
		Balances:                    state.Balances(),
		RandaoMixes:                 state.RandaoMixes(),
		Slashings:                   state.Slashings(),
		PreviousEpochParticipation:  prevEpochParticipation,
		CurrentEpochParticipation:   currentEpochParticipation,
		JustificationBits:           state.JustificationBits(),
		PreviousJustifiedCheckpoint: state.PreviousJustifiedCheckpoint(),
		CurrentJustifiedCheckpoint:  state.CurrentJustifiedCheckpoint(),
		FinalizedCheckpoint:         state.FinalizedCheckpoint(),
		InactivityScores:            inactivityScores,
		CurrentSyncCommittee:        currentSyncCommittee,
		NextSyncCommittee:           nextSyncCommittee,
		LatestExecutionPayloadHeader: &enginev1.ExecutionPayloadHeaderCapella{
			ParentHash:       payloadHeader.ParentHash(),
			FeeRecipient:     payloadHeader.FeeRecipient(),
			StateRoot:        payloadHeader.StateRoot(),
			ReceiptsRoot:     payloadHeader.ReceiptsRoot(),
			LogsBloom:        payloadHeader.LogsBloom(),
			PrevRandao:       payloadHeader.PrevRandao(),
			BlockNumber:      payloadHeader.BlockNumber(),
			GasLimit:         payloadHeader.GasLimit(),
			GasUsed:          payloadHeader.GasUsed(),
			Timestamp:        payloadHeader.Timestamp(),
			ExtraData:        payloadHeader.ExtraData(),
			BaseFeePerGas:    payloadHeader.BaseFeePerGas(),
			BlockHash:        payloadHeader.BlockHash(),
			TransactionsRoot: txRoot,
			WithdrawalsRoot:  make([]byte, 32),
		},
		NextWithdrawalIndex:          0,
		NextWithdrawalValidatorIndex: 0,
	}

	return state_native.InitializeFromProtoUnsafeCapella(s)
}
