package blocks_test

import (
	"testing"

	"github.com/prysmaticlabs/prysm/beacon-chain/core/blocks"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	"github.com/prysmaticlabs/prysm/config/params"
	"github.com/prysmaticlabs/prysm/encoding/bytesutil"
	ethpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1/wrapper"
	"github.com/prysmaticlabs/prysm/testing/assert"
	"github.com/prysmaticlabs/prysm/testing/require"
	"github.com/prysmaticlabs/prysm/testing/util"
)

func TestVerifyBlockSignatureUsingCurrentFork(t *testing.T) {
	params.SetupTestConfigCleanup(t)
	bCfg := params.BeaconConfig()
	bCfg.AltairForkEpoch = 100
	bCfg.ForkVersionSchedule[bytesutil.ToBytes4(bCfg.AltairForkVersion)] = 100
	params.OverrideBeaconConfig(bCfg)
	bState, keys := util.DeterministicGenesisState(t, 100)
	altairBlk := util.NewBeaconBlockAltair()
	altairBlk.Block.ProposerIndex = 0
	altairBlk.Block.Slot = params.BeaconConfig().SlotsPerEpoch * 100
	fData := &ethpb.Fork{
		Epoch:           100,
		CurrentVersion:  params.BeaconConfig().AltairForkVersion,
		PreviousVersion: params.BeaconConfig().GenesisForkVersion,
	}
	domain, err := helpers.Domain(fData, 100, params.BeaconConfig().DomainBeaconProposer, bState.GenesisValidatorRoot())
	assert.NoError(t, err)
	rt, err := helpers.ComputeSigningRoot(altairBlk.Block, domain)
	assert.NoError(t, err)
	sig := keys[0].Sign(rt[:]).Marshal()
	altairBlk.Signature = sig
	wsb, err := wrapper.WrappedAltairSignedBeaconBlock(altairBlk)
	require.NoError(t, err)
	assert.NoError(t, blocks.VerifyBlockSignatureUsingCurrentFork(bState, wsb))
}
