package params_test

import (
	"path"
	"path/filepath"
	"testing"

	"github.com/bazelbuild/rules_go/go/tools/bazel"
	"github.com/prysmaticlabs/prysm/config/params"
	"github.com/prysmaticlabs/prysm/io/file"
	"github.com/prysmaticlabs/prysm/testing/assert"
	"github.com/prysmaticlabs/prysm/testing/require"
)

func testnetConfigFilePath(t *testing.T, network string) string {
	filepath, err := bazel.Runfile("external/eth2_networks")
	require.NoError(t, err)
	configFilePath := path.Join(filepath, "shared", network, "config.yaml")
	return configFilePath
}

func TestE2EConfigParity(t *testing.T) {
	params.SetupTestConfigCleanup(t)
	testDir := bazel.TestTmpDir()
	yamlDir := filepath.Join(testDir, "config.yaml")

	testCfg := params.E2EMainnetTestConfig()
	yamlObj := params.ConfigToYaml(testCfg)
	assert.NoError(t, file.WriteFile(yamlDir, yamlObj))

	cfg := params.UnmarshalChainConfigFile(yamlDir, params.MainnetConfig().Copy())

	// compareConfigs makes it easier to figure out exactly what changed
	compareConfigs(t, cfg, testCfg)
	// failsafe in case compareConfigs is not updated when new fields are added
	require.DeepEqual(t, cfg, testCfg)
}

func compareConfigs(t *testing.T, expected, actual *params.BeaconChainConfig) {
	require.DeepEqual(t, expected.GenesisEpoch, actual.GenesisEpoch)
	require.DeepEqual(t, expected.FarFutureEpoch, actual.FarFutureEpoch)
	require.DeepEqual(t, expected.FarFutureSlot, actual.FarFutureSlot)
	require.DeepEqual(t, expected.BaseRewardsPerEpoch, actual.BaseRewardsPerEpoch)
	require.DeepEqual(t, expected.DepositContractTreeDepth, actual.DepositContractTreeDepth)
	require.DeepEqual(t, expected.JustificationBitsLength, actual.JustificationBitsLength)
	require.DeepEqual(t, expected.PresetBase, actual.PresetBase)
	require.DeepEqual(t, expected.ConfigName, actual.ConfigName)
	require.DeepEqual(t, expected.TargetCommitteeSize, actual.TargetCommitteeSize)
	require.DeepEqual(t, expected.MaxValidatorsPerCommittee, actual.MaxValidatorsPerCommittee)
	require.DeepEqual(t, expected.MaxCommitteesPerSlot, actual.MaxCommitteesPerSlot)
	require.DeepEqual(t, expected.MinPerEpochChurnLimit, actual.MinPerEpochChurnLimit)
	require.DeepEqual(t, expected.ChurnLimitQuotient, actual.ChurnLimitQuotient)
	require.DeepEqual(t, expected.ShuffleRoundCount, actual.ShuffleRoundCount)
	require.DeepEqual(t, expected.MinGenesisActiveValidatorCount, actual.MinGenesisActiveValidatorCount)
	require.DeepEqual(t, expected.MinGenesisTime, actual.MinGenesisTime)
	require.DeepEqual(t, expected.TargetAggregatorsPerCommittee, actual.TargetAggregatorsPerCommittee)
	require.DeepEqual(t, expected.HysteresisQuotient, actual.HysteresisQuotient)
	require.DeepEqual(t, expected.HysteresisDownwardMultiplier, actual.HysteresisDownwardMultiplier)
	require.DeepEqual(t, expected.HysteresisUpwardMultiplier, actual.HysteresisUpwardMultiplier)
	require.DeepEqual(t, expected.MinDepositAmount, actual.MinDepositAmount)
	require.DeepEqual(t, expected.MaxEffectiveBalance, actual.MaxEffectiveBalance)
	require.DeepEqual(t, expected.EjectionBalance, actual.EjectionBalance)
	require.DeepEqual(t, expected.EffectiveBalanceIncrement, actual.EffectiveBalanceIncrement)
	require.DeepEqual(t, expected.BLSWithdrawalPrefixByte, actual.BLSWithdrawalPrefixByte)
	require.DeepEqual(t, expected.ZeroHash, actual.ZeroHash)
	require.DeepEqual(t, expected.GenesisDelay, actual.GenesisDelay)
	require.DeepEqual(t, expected.MinAttestationInclusionDelay, actual.MinAttestationInclusionDelay)
	require.DeepEqual(t, expected.SecondsPerSlot, actual.SecondsPerSlot)
	require.DeepEqual(t, expected.SlotsPerEpoch, actual.SlotsPerEpoch)
	require.DeepEqual(t, expected.SqrRootSlotsPerEpoch, actual.SqrRootSlotsPerEpoch)
	require.DeepEqual(t, expected.MinSeedLookahead, actual.MinSeedLookahead)
	require.DeepEqual(t, expected.MaxSeedLookahead, actual.MaxSeedLookahead)
	require.DeepEqual(t, expected.EpochsPerEth1VotingPeriod, actual.EpochsPerEth1VotingPeriod)
	require.DeepEqual(t, expected.SlotsPerHistoricalRoot, actual.SlotsPerHistoricalRoot)
	require.DeepEqual(t, expected.MinValidatorWithdrawabilityDelay, actual.MinValidatorWithdrawabilityDelay)
	require.DeepEqual(t, expected.ShardCommitteePeriod, actual.ShardCommitteePeriod)
	require.DeepEqual(t, expected.MinEpochsToInactivityPenalty, actual.MinEpochsToInactivityPenalty)
	require.DeepEqual(t, expected.Eth1FollowDistance, actual.Eth1FollowDistance)
	require.DeepEqual(t, expected.SafeSlotsToUpdateJustified, actual.SafeSlotsToUpdateJustified)
	require.DeepEqual(t, expected.SafeSlotsToImportOptimistically, actual.SafeSlotsToImportOptimistically)
	require.DeepEqual(t, expected.SecondsPerETH1Block, actual.SecondsPerETH1Block)
	require.DeepEqual(t, expected.ProposerScoreBoost, actual.ProposerScoreBoost)
	require.DeepEqual(t, expected.IntervalsPerSlot, actual.IntervalsPerSlot)
	require.DeepEqual(t, expected.DepositChainID, actual.DepositChainID)
	require.DeepEqual(t, expected.DepositNetworkID, actual.DepositNetworkID)
	require.DeepEqual(t, expected.DepositContractAddress, actual.DepositContractAddress)
	require.DeepEqual(t, expected.RandomSubnetsPerValidator, actual.RandomSubnetsPerValidator)
	require.DeepEqual(t, expected.EpochsPerRandomSubnetSubscription, actual.EpochsPerRandomSubnetSubscription)
	require.DeepEqual(t, expected.EpochsPerHistoricalVector, actual.EpochsPerHistoricalVector)
	require.DeepEqual(t, expected.EpochsPerSlashingsVector, actual.EpochsPerSlashingsVector)
	require.DeepEqual(t, expected.HistoricalRootsLimit, actual.HistoricalRootsLimit)
	require.DeepEqual(t, expected.ValidatorRegistryLimit, actual.ValidatorRegistryLimit)
	require.DeepEqual(t, expected.BaseRewardFactor, actual.BaseRewardFactor)
	require.DeepEqual(t, expected.WhistleBlowerRewardQuotient, actual.WhistleBlowerRewardQuotient)
	require.DeepEqual(t, expected.ProposerRewardQuotient, actual.ProposerRewardQuotient)
	require.DeepEqual(t, expected.InactivityPenaltyQuotient, actual.InactivityPenaltyQuotient)
	require.DeepEqual(t, expected.MinSlashingPenaltyQuotient, actual.MinSlashingPenaltyQuotient)
	require.DeepEqual(t, expected.ProportionalSlashingMultiplier, actual.ProportionalSlashingMultiplier)
	require.DeepEqual(t, expected.MaxProposerSlashings, actual.MaxProposerSlashings)
	require.DeepEqual(t, expected.MaxAttesterSlashings, actual.MaxAttesterSlashings)
	require.DeepEqual(t, expected.MaxAttestations, actual.MaxAttestations)
	require.DeepEqual(t, expected.MaxDeposits, actual.MaxDeposits)
	require.DeepEqual(t, expected.MaxVoluntaryExits, actual.MaxVoluntaryExits)
	require.DeepEqual(t, expected.DomainBeaconProposer, actual.DomainBeaconProposer)
	require.DeepEqual(t, expected.DomainRandao, actual.DomainRandao)
	require.DeepEqual(t, expected.DomainBeaconAttester, actual.DomainBeaconAttester)
	require.DeepEqual(t, expected.DomainDeposit, actual.DomainDeposit)
	require.DeepEqual(t, expected.DomainVoluntaryExit, actual.DomainVoluntaryExit)
	require.DeepEqual(t, expected.DomainSelectionProof, actual.DomainSelectionProof)
	require.DeepEqual(t, expected.DomainAggregateAndProof, actual.DomainAggregateAndProof)
	require.DeepEqual(t, expected.DomainSyncCommittee, actual.DomainSyncCommittee)
	require.DeepEqual(t, expected.DomainSyncCommitteeSelectionProof, actual.DomainSyncCommitteeSelectionProof)
	require.DeepEqual(t, expected.DomainContributionAndProof, actual.DomainContributionAndProof)
	require.DeepEqual(t, expected.GweiPerEth, actual.GweiPerEth)
	require.DeepEqual(t, expected.BLSSecretKeyLength, actual.BLSSecretKeyLength)
	require.DeepEqual(t, expected.BLSPubkeyLength, actual.BLSPubkeyLength)
	require.DeepEqual(t, expected.DefaultBufferSize, actual.DefaultBufferSize)
	require.DeepEqual(t, expected.ValidatorPrivkeyFileName, actual.ValidatorPrivkeyFileName)
	require.DeepEqual(t, expected.WithdrawalPrivkeyFileName, actual.WithdrawalPrivkeyFileName)
	require.DeepEqual(t, expected.RPCSyncCheck, actual.RPCSyncCheck)
	require.DeepEqual(t, expected.EmptySignature, actual.EmptySignature)
	require.DeepEqual(t, expected.DefaultPageSize, actual.DefaultPageSize)
	require.DeepEqual(t, expected.MaxPeersToSync, actual.MaxPeersToSync)
	require.DeepEqual(t, expected.SlotsPerArchivedPoint, actual.SlotsPerArchivedPoint)
	require.DeepEqual(t, expected.GenesisCountdownInterval, actual.GenesisCountdownInterval)
	require.DeepEqual(t, expected.BeaconStateFieldCount, actual.BeaconStateFieldCount)
	require.DeepEqual(t, expected.BeaconStateAltairFieldCount, actual.BeaconStateAltairFieldCount)
	require.DeepEqual(t, expected.BeaconStateBellatrixFieldCount, actual.BeaconStateBellatrixFieldCount)
	require.DeepEqual(t, expected.WeakSubjectivityPeriod, actual.WeakSubjectivityPeriod)
	require.DeepEqual(t, expected.PruneSlasherStoragePeriod, actual.PruneSlasherStoragePeriod)
	require.DeepEqual(t, expected.SlashingProtectionPruningEpochs, actual.SlashingProtectionPruningEpochs)
	require.DeepEqual(t, expected.GenesisForkVersion, actual.GenesisForkVersion)
	require.DeepEqual(t, expected.AltairForkVersion, actual.AltairForkVersion)
	require.DeepEqual(t, expected.AltairForkEpoch, actual.AltairForkEpoch)
	require.DeepEqual(t, expected.BellatrixForkVersion, actual.BellatrixForkVersion)
	require.DeepEqual(t, expected.BellatrixForkEpoch, actual.BellatrixForkEpoch)
	require.DeepEqual(t, expected.ShardingForkVersion, actual.ShardingForkVersion)
	require.DeepEqual(t, expected.ShardingForkEpoch, actual.ShardingForkEpoch)
	require.DeepEqual(t, expected.ForkVersionSchedule, actual.ForkVersionSchedule)
	require.DeepEqual(t, expected.SafetyDecay, actual.SafetyDecay)
	require.DeepEqual(t, expected.TimelySourceFlagIndex, actual.TimelySourceFlagIndex)
	require.DeepEqual(t, expected.TimelyTargetFlagIndex, actual.TimelyTargetFlagIndex)
	require.DeepEqual(t, expected.TimelyHeadFlagIndex, actual.TimelyHeadFlagIndex)
	require.DeepEqual(t, expected.TimelySourceWeight, actual.TimelySourceWeight)
	require.DeepEqual(t, expected.TimelyTargetWeight, actual.TimelyTargetWeight)
	require.DeepEqual(t, expected.TimelyHeadWeight, actual.TimelyHeadWeight)
	require.DeepEqual(t, expected.SyncRewardWeight, actual.SyncRewardWeight)
	require.DeepEqual(t, expected.WeightDenominator, actual.WeightDenominator)
	require.DeepEqual(t, expected.ProposerWeight, actual.ProposerWeight)
	require.DeepEqual(t, expected.TargetAggregatorsPerSyncSubcommittee, actual.TargetAggregatorsPerSyncSubcommittee)
	require.DeepEqual(t, expected.SyncCommitteeSubnetCount, actual.SyncCommitteeSubnetCount)
	require.DeepEqual(t, expected.SyncCommitteeSize, actual.SyncCommitteeSize)
	require.DeepEqual(t, expected.InactivityScoreBias, actual.InactivityScoreBias)
	require.DeepEqual(t, expected.InactivityScoreRecoveryRate, actual.InactivityScoreRecoveryRate)
	require.DeepEqual(t, expected.EpochsPerSyncCommitteePeriod, actual.EpochsPerSyncCommitteePeriod)
	require.DeepEqual(t, expected.InactivityPenaltyQuotientAltair, actual.InactivityPenaltyQuotientAltair)
	require.DeepEqual(t, expected.MinSlashingPenaltyQuotientAltair, actual.MinSlashingPenaltyQuotientAltair)
	require.DeepEqual(t, expected.ProportionalSlashingMultiplierAltair, actual.ProportionalSlashingMultiplierAltair)
	require.DeepEqual(t, expected.MinSlashingPenaltyQuotientBellatrix, actual.MinSlashingPenaltyQuotientBellatrix)
	require.DeepEqual(t, expected.ProportionalSlashingMultiplierBellatrix, actual.ProportionalSlashingMultiplierBellatrix)
	require.DeepEqual(t, expected.InactivityPenaltyQuotientBellatrix, actual.InactivityPenaltyQuotientBellatrix)
	require.DeepEqual(t, expected.MinSyncCommitteeParticipants, actual.MinSyncCommitteeParticipants)
	require.DeepEqual(t, expected.TerminalBlockHash, actual.TerminalBlockHash)
	require.DeepEqual(t, expected.TerminalBlockHashActivationEpoch, actual.TerminalBlockHashActivationEpoch)
	require.DeepEqual(t, expected.TerminalTotalDifficulty, actual.TerminalTotalDifficulty)
	require.DeepEqual(t, expected.DefaultFeeRecipient, actual.DefaultFeeRecipient)
}
