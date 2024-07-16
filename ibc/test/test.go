package test

import (
	"testing"

	"github.com/sebastianmontero/bennyfi-go-client/ibc"
	"github.com/sebastianmontero/eos-go-toolbox/test"
	"gotest.tools/assert"
)

type TestUtil struct {
	ibcClient *ibc.IBCContract
	*test.TestUtil
}

func NewTestUtil(t *testing.T, ibcClient *ibc.IBCContract) *TestUtil {
	return &TestUtil{
		ibcClient: ibcClient,
		TestUtil:  test.NewTestUtil(t, ibcClient.EOS),
	}
}

func (m *TestUtil) AssertHeavyProofExists(exists bool) {
	actual, err := m.ibcClient.GetHeavyProof()
	assert.NilError(m.T, err)
	assert.Equal(m.T, actual != nil, exists)
}
