package test

import (
	"testing"
	"time"

	ltest "github.com/sebastianmontero/bennyfi-go-client/util/test"
	"github.com/sebastianmontero/bennyfi-go-client/yield/source/adaptor/stakelocal"
	"github.com/sebastianmontero/eos-go-toolbox/test"
	"github.com/sebastianmontero/eos-go-toolbox/util"
	"gotest.tools/assert"
)

type TestUtil struct {
	stakeLocalClient *stakelocal.StakeLocalContract
	*test.TestUtil
}

func NewTestUtil(t *testing.T, stakeLocalClient *stakelocal.StakeLocalContract) *TestUtil {
	return &TestUtil{
		stakeLocalClient: stakeLocalClient,
		TestUtil:         test.NewTestUtil(t, stakeLocalClient.EOS),
	}
}

func (m *TestUtil) AssertStakeByID(expected *stakelocal.Stake) *stakelocal.Stake {
	actual, err := m.stakeLocalClient.GetStake(expected.RoundID)
	assert.NilError(m.T, err)
	m.AssertStake(actual, expected)
	return actual

}

func (m *TestUtil) AssertStake(actual, expected *stakelocal.Stake) {
	assert.Check(m.T, actual != nil)
	assert.Equal(m.T, actual.RoundID, expected.RoundID)
	assert.Equal(m.T, actual.Cycle, expected.Cycle)
	assert.Equal(m.T, actual.State, expected.State)
	assert.Equal(m.T, actual.LastCycleReturn, expected.LastCycleReturn)
	assert.Equal(m.T, actual.State, expected.State)
	assert.Equal(m.T, actual.TotalReturn, expected.TotalReturn)
	assert.Equal(m.T, actual.TotalStake, expected.TotalStake)
	assert.DeepEqual(m.T, actual.StakingPeriod, expected.StakingPeriod)
	stakingPeriod := actual.StakingPeriod.AsTimeDuration()
	shift := stakingPeriod + time.Minute*2
	dateLimit := time.Now().Add(shift * -1)
	assert.Assert(m.T, actual.StakedTime.Time().After(dateLimit), "Expected staked time to be set staked time: %v, date limit: %v", actual.StakedTime, dateLimit)
	assert.Assert(m.T, actual.CycleStakedTime.Time().After(dateLimit), "Expected cycle staked time to be set cycle staked time: %v, date limit: %v", actual.CycleStakedTime, dateLimit)

	if actual.State == stakelocal.StakeStateUnstaked {
		assert.Assert(m.T, actual.StakeEndTime.Time().After(dateLimit), "Expected stake end time to be set stake end time: %v, date limit: %v", actual.StakeEndTime, dateLimit)
	} else {
		assert.Assert(m.T, util.IsNullTimePoint(actual.StakeEndTime), "Expected stake end time not to be set")
	}
	ltest.AssertAdditionalFields(m.T, actual.AdditionalFields, expected.AdditionalFields)
}

func (m *TestUtil) CheckStakesAre(stakes []stakelocal.Stake, ids []uint64) {
	assert.Equal(m.T, len(ids), len(stakes))
	for i, id := range ids {
		assert.Equal(m.T, id, stakes[i].RoundID)
	}
}
