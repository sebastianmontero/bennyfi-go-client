package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/sebastianmontero/bennyfi-go-client/yield/source/adaptor/eosrexflow"
	"github.com/sebastianmontero/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/dto"
	"github.com/sebastianmontero/eos-go-toolbox/test"
	"github.com/sebastianmontero/eos-go-toolbox/util"
	"gotest.tools/assert"
)

type TestUtil struct {
	eosRexClient *eosrexflow.EosRexFlowContract
	rexContract  eos.AccountName
	*test.TestUtil
}

func NewTestUtil(t *testing.T, eosRexClient *eosrexflow.EosRexFlowContract, rexContract eos.AccountName) *TestUtil {
	return &TestUtil{
		eosRexClient: eosRexClient,
		rexContract:  rexContract,
		TestUtil:     test.NewTestUtil(t, eosRexClient.EOS),
	}
}

func (m *TestUtil) AssertStakeByID(expected *eosrexflow.Stake, sellDelay bool) *eosrexflow.Stake {
	actual, err := m.eosRexClient.GetStake(expected.RoundID)
	assert.NilError(m.T, err)
	m.AssertStake(actual, expected, sellDelay)
	return actual

}

func (m *TestUtil) AssertStake(actual, expected *eosrexflow.Stake, sellDelay bool) {
	// time.Sleep(time.Millisecond * 500)
	assert.Check(m.T, actual != nil)
	assert.Equal(m.T, actual.RoundID, expected.RoundID)
	assert.Equal(m.T, actual.InitialStake, expected.InitialStake)
	assert.Equal(m.T, actual.CycleStake, expected.CycleStake)
	assert.Equal(m.T, actual.RexBalance, expected.RexBalance)
	assert.Equal(m.T, actual.CycleReturn, expected.CycleReturn)
	assert.Equal(m.T, actual.TotalReturn, expected.TotalReturn)
	assert.Equal(m.T, actual.RexState, expected.RexState)
	assert.Equal(m.T, actual.Cycle, expected.Cycle)
	assert.DeepEqual(m.T, actual.CycleStakingPeriod, expected.CycleStakingPeriod)
	assert.DeepEqual(m.T, actual.StakingPeriod, expected.StakingPeriod)
	stakingPeriod := actual.StakingPeriod.AsTimeDuration()
	shift := stakingPeriod + time.Minute
	dateLimit := time.Now().Add(shift * -1)
	stakedTime := actual.StakedTime
	cycleStakedTime := actual.CycleStakedTime
	finalCycleBufferPeriodDays, err := m.eosRexClient.SettingAsUint32(eosrexflow.SettingFinalCycleBufferPeriodDays)
	assert.NilError(m.T, err)
	nextStakeEndTime := actual.GetNextStakeEndTime(finalCycleBufferPeriodDays)
	stakeEndTime := actual.StakeEndTime
	assert.Assert(m.T, cycleStakedTime.Time().Add(-1*time.Second).Before(time.Now()))
	assert.Assert(m.T, stakedTime.Time().Add(-1*time.Second).Before(time.Now()))
	assert.Assert(m.T, cycleStakedTime.Time().After(dateLimit), "Expected cycle staked time to be set")
	assert.Assert(m.T, stakedTime.Time().After(dateLimit), "Expected staked time to be set")
	assert.Assert(m.T, !cycleStakedTime.Time().Before(stakedTime.Time()), "Expected cycle staked time to be after or equal to staked time")
	assert.Assert(m.T, stakeEndTime.Time().Equal(nextStakeEndTime.Time()), "Expected stake end time to be set and be equal to next stake end time")
	if actual.RexState == eosrexflow.RexStateInSavings {
		assert.Assert(m.T, util.IsNullTimePoint(actual.MovedFromSavingsTime), "Expected moved from savings time not to be set")
		assert.Assert(m.T, util.IsNullTimePoint(actual.MaturityTime), "Expected maturity time not to be set")
		assert.Assert(m.T, stakeEndTime.Time().After(time.Now()))
	} else {
		movedFromSavingsTime := actual.MovedFromSavingsTime
		maturityTime := actual.MaturityTime
		rexLockPeriod := time.Duration(int64(eosrexflow.RexLockPeriodDays) * int64(time.Hour) * 24)
		expectedMoveFromSavingsTime := stakeEndTime.Time().Add(rexLockPeriod * -1)
		assert.Equal(m.T, movedFromSavingsTime.Time(), expectedMoveFromSavingsTime)
		assert.Assert(m.T, time.Since(movedFromSavingsTime.Time().Add(-1*time.Second)) >= 0)
		if actual.RexState == eosrexflow.RexStateInLockPeriod {
			stakeEndTimeUTC := stakeEndTime.Time().UTC()
			assert.Equal(m.T, maturityTime.Time().UTC(), time.Date(stakeEndTimeUTC.Year(), stakeEndTimeUTC.Month(), stakeEndTimeUTC.Day(), 0, 0, 0, 0, stakeEndTimeUTC.Location()), fmt.Sprintf("Expected maturity time to be start of day UTC of stake end time, maturity time: %v, stake end time: %v", maturityTime.Time().UTC(), stakeEndTime.Time().UTC()))
			assert.Assert(m.T, time.Now().Before(maturityTime.Time()))
		} else if actual.RexState == eosrexflow.RexStateProceedsCalculated {
			bufferPeriod, err := m.eosRexClient.SettingAsUint32(eosrexflow.SettingProceedsCalculationBufferPeriodMins)
			assert.NilError(m.T, err)
			assert.Assert(m.T, time.Since(maturityTime.Time().Add(-1*time.Second).Add(-1*time.Minute*time.Duration(bufferPeriod))) >= 0)
		} else if actual.RexState == eosrexflow.RexStateWithdrawn {
			assert.Assert(m.T, time.Since(stakeEndTime.Time().Add(-1*time.Second)) >= 0)
		} else {
			if !sellDelay {
				assert.Assert(m.T, stakeEndTime.Time().After(time.Now()))
			}
		}
	}
}

func (m *TestUtil) AssertTREXNotifications(account eos.AccountName, notifications [][]string, verifyExactQuantity bool) {
	m.BaseAssertNotifications(account, "notifytrex", notifications, verifyExactQuantity)
}

func (m *TestUtil) AssertInvestedInRex(stakeAmount, rexAmount eos.Asset) {
	m.AssertAction(m.rexContract, "deposit", map[string]interface{}{
		"owner":  m.eosRexClient.ContractName,
		"amount": stakeAmount.String(),
	}, 0)

	m.AssertAction(m.rexContract, "buyrex", map[string]interface{}{
		"from":   m.eosRexClient.ContractName,
		"amount": stakeAmount.String(),
	}, 0)

	m.AssertNumActions(m.rexContract, "mvtosavings", 0)
}

func (m *TestUtil) AssertRexMovedFromSavings(rexBalance eos.Asset) {
	m.AssertAction(m.rexContract, "mvfrsavings", map[string]interface{}{
		"owner": m.eosRexClient.ContractName,
		"rex":   rexBalance.String(),
	}, 0)
}

func (m *TestUtil) AssertUpdateRexCalled() {
	m.AssertAction(m.rexContract, "updaterex", map[string]interface{}{
		"owner": m.eosRexClient.ContractName,
	}, 0)
}

func (m *TestUtil) AssertUpdateRexCalledNTimes(times int) {
	m.AssertNumActions(m.rexContract, "updaterex", times)
}

func (m *TestUtil) AssertRexSold(rexBalance eos.Asset) {
	m.AssertAction(m.rexContract, "sellrex", map[string]interface{}{
		"from": m.eosRexClient.ContractName,
		"rex":  rexBalance.String(),
	}, 0)
}

func (m *TestUtil) AssertRexSellCalledNTimes(times int) {
	m.AssertNumActions(m.rexContract, "sellrex", times)
}

func (m *TestUtil) AssertRexWithdrawn(amount eos.Asset) {
	m.AssertAction(m.rexContract, "withdraw", map[string]interface{}{
		"owner":  m.eosRexClient.ContractName,
		"amount": amount.String(),
	}, 0)
}

func (m *TestUtil) AssertRexWithdrawCalledNTimes(times int) {
	m.AssertNumActions(m.rexContract, "withdraw", times)
}

func (m *TestUtil) LapseLastNotifiedTime(lastNotifiedSetting string, shift time.Duration) {
	notificationPeriodMins, err := m.eosRexClient.SettingAsUint32(eosrexflow.SettingNotificationPeriodMins)
	assert.NilError(m.T, err)
	lastNotified := time.Now().Add(time.Duration(notificationPeriodMins)*time.Minute*-1 + shift)
	_, err = m.eosRexClient.SetSetting(eos.AccountName(m.eosRexClient.ContractName), lastNotifiedSetting, dto.FlexValueFromTime(lastNotified))
	assert.NilError(m.T, err)
}
