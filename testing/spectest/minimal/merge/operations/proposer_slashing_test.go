package operations

import (
	"testing"

	"github.com/prysmaticlabs/prysm/testing/spectest/shared/merge/operations"
)

func TestMinimal_Merge_Operations_ProposerSlashing(t *testing.T) {
	operations.RunProposerSlashingTest(t, "minimal")
}
