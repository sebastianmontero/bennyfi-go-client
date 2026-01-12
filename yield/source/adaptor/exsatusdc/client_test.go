package exsatusdc

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"gotest.tools/assert"
)

func TestABI(t *testing.T) {
	parsedABI, err := abi.JSON(strings.NewReader(ExSatBankYieldSourceAdaptorABI))
	assert.NilError(t, err)

	// Verify 'stakes' method
	method, ok := parsedABI.Methods["stakes"]
	assert.Assert(t, ok, "Method 'stakes' not found in ABI")
	assert.Equal(t, len(method.Inputs), 1)
	assert.Equal(t, len(method.Inputs), 1)
	// Note: method.Inputs[0].Type.String() would be "uint64"

	assert.Equal(t, len(method.Outputs), 7)

	// Verify 'triggerUnstake' method
	method, ok = parsedABI.Methods["triggerUnstake"]
	assert.Assert(t, ok, "Method 'triggerUnstake' not found in ABI")
	assert.Equal(t, len(method.Inputs), 1)
	assert.Equal(t, len(method.Outputs), 0)
}
