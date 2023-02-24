package rewards

import (
	"testing"

	"github.com/prysmaticlabs/prysm/v3/testing/spectest/shared/bellatrix/rewards"
)

func TestMainnet_Bellatrix_Rewards(t *testing.T) {
	rewards.RunPrecomputeRewardsAndPenaltiesTests(t, "mainnet")
}
