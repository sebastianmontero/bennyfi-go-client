package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/sebastianmontero/bennyfi-go-client/yield/source/adaptor/eosrexremote"
	"github.com/sebastianmontero/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/test"
	"github.com/sebastianmontero/eos-go-toolbox/util"
	"gotest.tools/assert"
)

type TestUtil struct {
	t                    *testing.T
	eosRexRemoteClient   *eosrexremote.EosRexRemoteContract
	rexContract          eos.AccountName
	eosRexRemoteContract eos.AccountName
	tokenContract        eos.AccountName
	ibcRemoteContract    eos.AccountName
	eosTestUtil          *test.TestUtil
}

func NewTestUtil(t *testing.T, eosRexRemoteClient *eosrexremote.EosRexRemoteContract, rexContract, tokenContract, ibcRemoteContract eos.AccountName) *TestUtil {
	return &TestUtil{
		t:                    t,
		eosRexRemoteClient:   eosRexRemoteClient,
		rexContract:          rexContract,
		eosRexRemoteContract: eos.AN(eosRexRemoteClient.ContractName),
		tokenContract:        tokenContract,
		ibcRemoteContract:    ibcRemoteContract,
		eosTestUtil:          test.NewTestUtil(t, eosRexRemoteClient.EOS),
	}
}

func (m *TestUtil) AssertStakeByID(expected *eosrexremote.Stake, sellDelay bool) *eosrexremote.Stake {
	actual, err := m.eosRexRemoteClient.GetStake(expected.RoundID)
	assert.NilError(m.t, err)
	m.AssertStake(actual, expected, sellDelay)
	return actual
}

func (m *TestUtil) AssertStake(actual, expected *eosrexremote.Stake, sellDelay bool) {
	// time.Sleep(time.Millisecond * 500)
	assert.Check(m.t, actual != nil)
	assert.Equal(m.t, actual.RoundID, expected.RoundID)
	assert.Equal(m.t, actual.TotalStake, expected.TotalStake)
	assert.Equal(m.t, actual.RexBalance, expected.RexBalance)
	assert.Equal(m.t, actual.TotalReturn, expected.TotalReturn)
	assert.Equal(m.t, actual.RexState, expected.RexState)
	assert.DeepEqual(m.t, actual.StakingPeriod, expected.StakingPeriod)
	stakingPeriod := actual.StakingPeriod.AsTimeDuration()
	shift := stakingPeriod + time.Minute
	dateLimit := time.Now().Add(shift * -1)
	stakedTime := actual.StakedTime
	stakeEndTime := actual.StakeEndTime
	assert.Assert(m.t, stakedTime.Time().Add(-1*time.Second).Before(time.Now()))
	assert.Assert(m.t, stakedTime.Time().After(dateLimit), "Expected staked time to be set")
	if actual.RexState == eosrexremote.RexStateInSavings {
		assert.Assert(m.t, util.IsNullTimePoint(actual.MovedFromSavingsTime), "Expected moved from savings time not to be set")
		assert.Assert(m.t, stakeEndTime.Time().After(time.Now()))
		assert.Assert(m.t, stakeEndTime.Time().Equal(stakedTime.Time().Add(stakingPeriod)), "Expected stake end time to be set and be staking period apart from staked time")
	} else {
		movedFromSavingsTime := actual.MovedFromSavingsTime
		rexLockPeriod := time.Duration(int64(eosrexremote.RexLockPeriodDays) * int64(time.Hour) * 24)
		expectedMoveFromSavingsTime := stakeEndTime.Time().Add(rexLockPeriod * -1)
		assert.Equal(m.t, movedFromSavingsTime.Time(), expectedMoveFromSavingsTime)
		assert.Equal(m.t, movedFromSavingsTime.Time().Add(rexLockPeriod), stakeEndTime.Time())
		assert.Assert(m.t, time.Since(movedFromSavingsTime.Time().Add(-1*time.Second)) >= 0)
		if actual.RexState == eosrexremote.RexStateSold || actual.RexState == eosrexremote.RexStateWithdrawn {
			assert.Assert(m.t, time.Since(stakeEndTime.Time().Add(-1*time.Second)) >= 0)
		} else {
			if !sellDelay {
				assert.Assert(m.t, stakeEndTime.Time().After(time.Now()))
			}
		}
	}
}

func (m *TestUtil) AssertREXNotifications(account eos.AccountName, notifications [][]string, verifyExactQuantity bool) {
	m.eosTestUtil.BaseAssertNotifications(account, "notifytrex", notifications, verifyExactQuantity)
}

func (m *TestUtil) AssertInvestedInRex(stakeAmount, rexAmount eos.Asset) {
	m.eosTestUtil.AssertAction(m.rexContract, "deposit", map[string]interface{}{
		"owner":  m.eosRexRemoteContract.String(),
		"amount": stakeAmount.String(),
	}, 0)

	m.eosTestUtil.AssertAction(m.rexContract, "buyrex", map[string]interface{}{
		"from":   m.eosRexRemoteContract.String(),
		"amount": stakeAmount.String(),
	}, 0)

	m.eosTestUtil.AssertAction(m.rexContract, "mvtosavings", map[string]interface{}{
		"owner": m.eosRexRemoteContract.String(),
		"rex":   rexAmount.String(),
	}, 0)
}

func (m *TestUtil) AssertRexMovedFromSavings(rexBalance eos.Asset) {
	m.eosTestUtil.AssertAction(m.rexContract, "mvfrsavings", map[string]interface{}{
		"owner": m.eosRexRemoteContract.String(),
		"rex":   rexBalance.String(),
	}, 0)
}

func (m *TestUtil) AssertRexSold(rexBalance eos.Asset) {
	m.eosTestUtil.AssertAction(m.rexContract, "sellrex", map[string]interface{}{
		"from": m.eosRexRemoteContract.String(),
		"rex":  rexBalance.String(),
	}, 0)
}

func (m *TestUtil) AssertRexWithdrawn(amount eos.Asset) {
	m.eosTestUtil.AssertAction(m.rexContract, "withdraw", map[string]interface{}{
		"owner":  m.eosRexRemoteContract.String(),
		"amount": amount.String(),
	}, 0)
}

func (m *TestUtil) AssertYieldReturnTransfer(roundId uint64, quantity eos.Asset) map[string]interface{} {
	actionData := map[string]interface{}{
		"from":     m.eosRexRemoteContract.String(),
		"to":       m.ibcRemoteContract.String(),
		"quantity": quantity.String(),
		"memo":     fmt.Sprintf("pool id: %v", roundId),
	}
	return m.eosTestUtil.AssertAction(m.tokenContract, "transfer", actionData, 0)
}
