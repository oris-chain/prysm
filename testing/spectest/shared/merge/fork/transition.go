package fork

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/snappy"
	types "github.com/prysmaticlabs/eth2-types"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/transition"
	"github.com/prysmaticlabs/prysm/beacon-chain/state"
	stateAltair "github.com/prysmaticlabs/prysm/beacon-chain/state/v2"
	v3 "github.com/prysmaticlabs/prysm/beacon-chain/state/v3"
	"github.com/prysmaticlabs/prysm/config/params"
	ethpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1/wrapper"
	wrapperv1 "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1/wrapper"
	"github.com/prysmaticlabs/prysm/testing/require"
	"github.com/prysmaticlabs/prysm/testing/spectest/utils"
	"github.com/prysmaticlabs/prysm/testing/util"
)

type ForkConfig struct {
	PostFork    string `json:"post_fork"`
	ForkEpoch   int    `json:"fork_epoch"`
	ForkBlock   int    `json:"fork_block"`
	BlocksCount int    `json:"blocks_count"`
}

// RunForkTransitionTest is a helper function that runs Merge's transition core tests.
func RunForkTransitionTest(t *testing.T, config string) {
	require.NoError(t, utils.SetConfig(t, config))

	testFolders, testsFolderPath := utils.TestFolders(t, config, "merge", "transition/core/pyspec_tests")
	for _, folder := range testFolders {
		t.Run(folder.Name(), func(t *testing.T) {
			helpers.ClearCache()
			file, err := util.BazelFileBytes(testsFolderPath, folder.Name(), "meta.yaml")
			require.NoError(t, err)
			config := &ForkConfig{}
			require.NoError(t, utils.UnmarshalYaml(file, config), "Failed to Unmarshal")

			preforkBlocks := make([]*ethpb.SignedBeaconBlockAltair, 0)
			postforkBlocks := make([]*ethpb.SignedBeaconBlockMerge, 0)
			// Fork happens without any pre-fork blocks.
			if config.ForkBlock == 0 {
				for i := 0; i < config.BlocksCount; i++ {
					fileName := fmt.Sprint("blocks_", i, ".ssz_snappy")
					blockFile, err := util.BazelFileBytes(testsFolderPath, folder.Name(), fileName)
					require.NoError(t, err)
					blockSSZ, err := snappy.Decode(nil /* dst */, blockFile)
					require.NoError(t, err, "Failed to decompress")
					block := &ethpb.SignedBeaconBlockMerge{}
					require.NoError(t, block.UnmarshalSSZ(blockSSZ), "Failed to unmarshal")
					postforkBlocks = append(postforkBlocks, block)
				}
				// Fork happens with pre-fork blocks.
			} else {
				for i := 0; i <= config.ForkBlock; i++ {
					fileName := fmt.Sprint("blocks_", i, ".ssz_snappy")
					blockFile, err := util.BazelFileBytes(testsFolderPath, folder.Name(), fileName)
					require.NoError(t, err)
					blockSSZ, err := snappy.Decode(nil /* dst */, blockFile)
					require.NoError(t, err, "Failed to decompress")
					block := &ethpb.SignedBeaconBlockAltair{}
					require.NoError(t, block.UnmarshalSSZ(blockSSZ), "Failed to unmarshal")
					preforkBlocks = append(preforkBlocks, block)
				}
				for i := config.ForkBlock + 1; i < config.BlocksCount; i++ {
					fileName := fmt.Sprint("blocks_", i, ".ssz_snappy")
					blockFile, err := util.BazelFileBytes(testsFolderPath, folder.Name(), fileName)
					require.NoError(t, err)
					blockSSZ, err := snappy.Decode(nil /* dst */, blockFile)
					require.NoError(t, err, "Failed to decompress")
					block := &ethpb.SignedBeaconBlockMerge{}
					require.NoError(t, block.UnmarshalSSZ(blockSSZ), "Failed to unmarshal")
					postforkBlocks = append(postforkBlocks, block)
				}
			}

			preBeaconStateFile, err := util.BazelFileBytes(testsFolderPath, folder.Name(), "pre.ssz_snappy")
			require.NoError(t, err)
			preBeaconStateSSZ, err := snappy.Decode(nil /* dst */, preBeaconStateFile)
			require.NoError(t, err, "Failed to decompress")
			beaconStateBase := &ethpb.BeaconStateAltair{}
			require.NoError(t, beaconStateBase.UnmarshalSSZ(preBeaconStateSSZ), "Failed to unmarshal")
			beaconState, err := stateAltair.InitializeFromProto(beaconStateBase)
			require.NoError(t, err)

			bc := params.BeaconConfig()
			bc.AltairForkEpoch = types.Epoch(config.ForkEpoch)
			params.OverrideBeaconConfig(bc)

			ctx := context.Background()
			var ok bool
			for _, b := range preforkBlocks {
				wsb, err := wrapperv1.WrappedAltairSignedBeaconBlock(b)
				require.NoError(t, err)
				st, err := transition.ExecuteStateTransition(ctx, beaconState, wsb)
				require.NoError(t, err)
				beaconState, ok = st.(*stateAltair.BeaconState)
				require.Equal(t, true, ok)
			}
			postState := state.BeaconState(beaconState)
			for _, b := range postforkBlocks {
				wsb, err := wrapper.WrappedMergeSignedBeaconBlock(b)
				require.NoError(t, err)
				st, err := transition.ExecuteStateTransition(ctx, postState, wsb)
				require.NoError(t, err)
				postState, ok = st.(*v3.BeaconState)
				require.Equal(t, true, ok)
			}

			postBeaconStateFile, err := util.BazelFileBytes(testsFolderPath, folder.Name(), "post.ssz_snappy")
			require.NoError(t, err)
			postBeaconStateSSZ, err := snappy.Decode(nil /* dst */, postBeaconStateFile)
			require.NoError(t, err, "Failed to decompress")
			postBeaconState := &ethpb.BeaconStateMerge{}
			require.NoError(t, postBeaconState.UnmarshalSSZ(postBeaconStateSSZ), "Failed to unmarshal")

			pbState, err := v3.ProtobufBeaconState(postState.CloneInnerState())
			require.NoError(t, err)
			require.DeepSSZEqual(t, pbState, postBeaconState)
		})
	}
}
