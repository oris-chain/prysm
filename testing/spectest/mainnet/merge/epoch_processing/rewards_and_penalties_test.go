package epoch_processing

import (
	"testing"

	"github.com/prysmaticlabs/prysm/testing/spectest/shared/merge/epoch_processing"
)

func TestMainnet_Merge_EpochProcessing_RewardsAndPenalties(t *testing.T) {
	epoch_processing.RunRewardsAndPenaltiesTests(t, "mainnet")
}
