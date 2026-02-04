package yield_source_registry

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"gotest.tools/assert"
)

func TestABI(t *testing.T) {
	parsedABI, err := abi.JSON(strings.NewReader(IYieldSourceRegistryABI))
	assert.NilError(t, err)

	// Verify 'getYieldSource'
	method, ok := parsedABI.Methods["getYieldSource"]
	assert.Assert(t, ok, "Method 'getYieldSource' not found")
	assert.Equal(t, len(method.Inputs), 1)
	assert.Equal(t, len(method.Outputs), 1)
	// Output should be a tuple (YieldSource struct)
	assert.Equal(t, method.Outputs[0].Type.T, abi.TupleTy)

	// Verify 'isYieldSourceRegistered'
	method, ok = parsedABI.Methods["isYieldSourceRegistered"]
	assert.Assert(t, ok, "Method 'isYieldSourceRegistered' not found")
	assert.Equal(t, len(method.Inputs), 1)
	assert.Equal(t, len(method.Outputs), 1)
	assert.Equal(t, method.Outputs[0].Type.T, abi.BoolTy)
}
