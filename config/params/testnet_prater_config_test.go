package params_test

import (
	"testing"

	"github.com/prysmaticlabs/prysm/config/params"
)

func TestPraterConfigMatchesUpstreamYaml(t *testing.T) {
	presetFPs := presetsFilePath(t, "mainnet")
	configFP := testnetConfigFilePath(t, "prater")
	cfg := params.UnmarshalChainConfigFile(configFP, nil)
	fields := fieldsFromYamls(t, append(presetFPs, configFP))
	assertYamlFieldsMatch(t, "prater", fields, cfg, params.PraterConfig())
}
