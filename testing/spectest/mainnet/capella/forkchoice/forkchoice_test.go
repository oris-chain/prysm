package forkchoice

import (
	"testing"

	"github.com/prysmaticlabs/prysm/v4/config/features"
	"github.com/prysmaticlabs/prysm/v4/runtime/version"
	"github.com/prysmaticlabs/prysm/v4/testing/spectest/shared/common/forkchoice"
)

func TestMainnet_Capella_Forkchoice(t *testing.T) {
	resetCfg := features.InitWithReset(&features.Flags{
		// Experimental features are disabled by default for spec tests.
		EnableDefensivePull: false,
		DisablePullTips:     true,
	})
	defer resetCfg()
	forkchoice.Run(t, "mainnet", version.Capella)
}
