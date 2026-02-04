package base

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"gotest.tools/assert"
)

func TestABI(t *testing.T) {
	parsedABI, err := abi.JSON(strings.NewReader(IYieldSourceAdaptorABI))
	assert.NilError(t, err)

	method, ok := parsedABI.Methods["stopStake"]
	assert.Assert(t, ok, "Method 'stopStake' not found")
	assert.Equal(t, len(method.Inputs), 1)
	assert.Equal(t, method.Inputs[0].Type.T, abi.UintTy)
	assert.Equal(t, method.Inputs[0].Type.Size, 64)

	method, ok = parsedABI.Methods["stake"]
	assert.Assert(t, ok, "Method 'stake' not found")
}
