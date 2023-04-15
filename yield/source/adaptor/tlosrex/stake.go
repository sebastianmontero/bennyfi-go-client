package tlosrex

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/eoscanada/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/dto"
)

var RexLockPeriodDays = 5

type Stake struct {
	RoundID              uint64            `json:"round_id"`
	TotalStake           eos.Asset         `json:"total_stake"`
	RexBalance           eos.Asset         `json:"rex_balance"`
	TotalReturn          eos.Asset         `json:"total_return"`
	RexState             eos.Name          `json:"rex_state"`
	StakingPeriod        *dto.Microseconds `json:"staking_period"`
	StakedTime           eos.TimePoint     `json:"staked_time"`
	MovedFromSavingsTime eos.TimePoint     `json:"moved_from_savings_time"`
	StakeEndTime         eos.TimePoint     `json:"stake_end_time"`
	LastNotifiedTime     eos.TimePoint     `json:"last_notified_time"`
	UpdatedDate          eos.TimePoint     `json:"updated_date"`
}

func (m *Stake) String() string {
	result, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Sprintf("Failed marshalling round: %v", err))
	}
	return string(result)
}

func (m *Stake) CalculateRexLockPeriodTime() eos.TimePoint {

	return eos.TimePoint(m.StakedTime.Time().Add(time.Hour * time.Duration(m.StakingPeriod.Hrs()-int64(24*RexLockPeriodDays))).UnixMicro())
}

func (m *Stake) CalculateSellRexTime() eos.TimePoint {

	return eos.TimePoint(m.MovedFromSavingsTime.Time().Add(time.Hour * time.Duration(24*RexLockPeriodDays)).UnixMicro())
}

func (m *Stake) CalculateStakeEndTime() eos.TimePoint {
	return eos.TimePoint(m.StakedTime.Time().Add(time.Hour * time.Duration(m.StakingPeriod.Hrs())).UnixMicro())
}

func (m *TlosRexContract) CheckStakeParameters(authorizer, tokenContract eos.AccountName, stakeAmount eos.Asset, stakingPeriodHrs uint32) (string, error) {
	actionData := struct {
		TokenContract    eos.AccountName
		StakeAmount      eos.Asset
		StakingPeriodHrs uint32
	}{tokenContract, stakeAmount, stakingPeriodHrs}
	return m.ExecAction(authorizer, "chckstkparam", actionData)
}

func (m *TlosRexContract) MoveRoundFromSavings(roundId uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "mvfrmsvngsrn", roundId)
}

func (m *TlosRexContract) SellRoundRex(roundId uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "sellrexrn", roundId)
}

func (m *TlosRexContract) WithdrawRoundRex(roundId uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "withdrwrexrn", roundId)
}

func (m *TlosRexContract) MoveFromSavings(callCounter uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "mvfrmsavings", callCounter)
}

func (m *TlosRexContract) SellRex(callCounter uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "sellrex", callCounter)
}

func (m *TlosRexContract) WithdrawRex(callCounter uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "withdrawrex", callCounter)
}

func (m *TlosRexContract) TstLapseTime(roundId uint64) (string, error) {
	actionData := struct {
		RoundId     uint64
		CallCounter uint64
	}{roundId, m.NextCallCounter()}
	return m.ExecAction(eos.AN(m.ContractName), "tstlapsetime", actionData)
}

func (m *TlosRexContract) TstSetLastNotifiedTime(roundId uint64, lastNotifiedTime eos.TimePoint) (string, error) {
	actionData := struct {
		RoundId          uint64
		LastNotifiedTime eos.TimePoint
	}{roundId, lastNotifiedTime}
	return m.ExecAction(eos.AN(m.ContractName), "setlastnotif", actionData)
}

func (m *Stake) Clone() *Stake {
	return &Stake{
		RoundID:              m.RoundID,
		TotalStake:           m.TotalStake,
		RexBalance:           m.RexBalance,
		TotalReturn:          m.TotalReturn,
		RexState:             m.RexState,
		StakingPeriod:        m.StakingPeriod,
		StakedTime:           m.StakedTime,
		MovedFromSavingsTime: m.MovedFromSavingsTime,
		StakeEndTime:         m.StakeEndTime,
		LastNotifiedTime:     m.LastNotifiedTime,
		UpdatedDate:          m.UpdatedDate,
	}
}

func (m *TlosRexContract) GetStake(roundID uint64) (*Stake, error) {
	stakes, err := m.GetStakesReq(&eos.GetTableRowsRequest{
		LowerBound: strconv.FormatUint(roundID, 10),
		UpperBound: strconv.FormatUint(roundID, 10),
		Limit:      1,
	})
	if err != nil {
		return nil, err
	}
	if len(stakes) > 0 {
		return &stakes[0], nil
	}
	return nil, nil
}

func (m *TlosRexContract) GetLastStake() (*Stake, error) {
	stakes, err := m.GetStakesReq(&eos.GetTableRowsRequest{
		Reverse: true,
		Limit:   1,
	})
	if err != nil {
		return nil, err
	}
	if len(stakes) > 0 {
		return &stakes[0], nil
	}
	return nil, nil
}

func (m *TlosRexContract) GetStakesReq(req *eos.GetTableRowsRequest) ([]Stake, error) {

	var stakes []Stake
	if req == nil {
		req = &eos.GetTableRowsRequest{}
	}
	req.Table = "stakes"
	err := m.GetTableRows(*req, &stakes)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	return stakes, nil
}

func (m *TlosRexContract) GetStakesByRexStateAndId(state eos.Name) ([]Stake, error) {
	request := &eos.GetTableRowsRequest{}
	err := m.FilterStakesByRexStateAndId(request, state)
	if err != nil {
		return nil, err
	}
	return m.GetStakesReq(request)
}

func (m *TlosRexContract) FilterStakesByRexStateAndId(req *eos.GetTableRowsRequest, state eos.Name) error {

	req.Index = "2"
	req.KeyType = "i128"
	req.Reverse = true
	stateAndRndLB, err := m.EOS.GetComposedIndexValue(state, 0)
	if err != nil {
		return fmt.Errorf("failed to generate lower bound composed index, err: %v", err)
	}
	stateAndRndUB, err := m.EOS.GetComposedIndexValue(state, uint64(18446744073709551615))
	if err != nil {
		return fmt.Errorf("failed to generate upper bound composed index, err: %v", err)
	}
	fmt.Println("LB: ", stateAndRndLB, "UB: ", stateAndRndUB)
	req.LowerBound = stateAndRndLB
	req.UpperBound = stateAndRndUB
	return err
}
