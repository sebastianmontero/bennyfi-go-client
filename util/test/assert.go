package test

import (
	"testing"

	"github.com/sebastianmontero/bennyfi-go-client/common/types"
	"gotest.tools/assert"
)

func AssertAdditionalFields(t *testing.T, actual, expected types.AdditionalFields) {

	assert.Equal(t, len(actual), len(expected), "Additional fields arrays length don't match, actual length: %v, expected length: %v, actual: %v, expected: %v", len(actual), len(expected), actual.String(), expected.String())
	for _, e := range expected {
		a := actual.Find(e.Key)
		assert.Assert(t, a != nil, "Expected additional field %v, not found", e)
		assert.Assert(t, a.Value.IsEqual(e.Value), "Additional field with key: %v do not match actual: %v, expected: %v", e.Key, a, e)
	}
}
