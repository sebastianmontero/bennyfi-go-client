package test

import (
	"testing"

	"github.com/sebastianmontero/bennyfi-go-client/yield/source/adaptor/common"
	"github.com/sebastianmontero/eos-go-toolbox/service"
	"github.com/sebastianmontero/eos-go-toolbox/test"
	"gotest.tools/assert"
)

type TestUtil struct {
	*test.TestUtil
	stoppedContract common.StoppedInterface
}

func NewTestUtil(t *testing.T, stoppedContract common.StoppedInterface, eos *service.EOS) *TestUtil {
	return &TestUtil{
		TestUtil:        test.NewTestUtil(t, eos),
		stoppedContract: stoppedContract,
	}
}

func (m *TestUtil) AssertAllStoppedStakes(expected []common.StoppedStake) {
	stoppedStakes, err := m.stoppedContract.GetAllStoppedStakes()
	assert.NilError(m.T, err)
	m.AssertStoppedStakes(stoppedStakes, expected)
}

func (m *TestUtil) AssertStoppedStakes(actual []common.StoppedStake, expected []common.StoppedStake) {
	assert.Equal(m.T, len(actual), len(expected))
	for i, actualStake := range actual {
		expectedStake := expected[i]
		m.AssertStoppedStake(actualStake, expectedStake)
	}
}

func (m *TestUtil) AssertStoppedStake(actual, expected common.StoppedStake) {
	assert.Equal(m.T, actual.RoundId, expected.RoundId)
	assert.Equal(m.T, actual.IsStopped, expected.IsStopped)
}
