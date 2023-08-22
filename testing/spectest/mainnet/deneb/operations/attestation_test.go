package operations

import (
	"testing"

	"github.com/prysmaticlabs/prysm/v4/testing/spectest/shared/deneb/operations"
)

func TestMainnet_Deneb_Operations_Attestation(t *testing.T) {
	operations.RunAttestationTest(t, "mainnet")
}
