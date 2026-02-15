package base

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeABIs(t *testing.T) {
	abi1 := `[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"}]`
	abi2 := `[{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"}]`

	merged, err := MergeABIs(abi1, abi2)
	assert.NoError(t, err)

	var resultArray []interface{}
	err = json.Unmarshal([]byte(merged), &resultArray)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(resultArray))
}

func TestMergeABIs_InvalidJSON(t *testing.T) {
	_, err := MergeABIs("invalid-json")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid ABI JSON")
}

func TestMergeABIs_Empty(t *testing.T) {
	merged, err := MergeABIs()
	assert.NoError(t, err)
	assert.Equal(t, "[]", merged)
}
