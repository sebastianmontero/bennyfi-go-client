package eosrexflow

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/sebastianmontero/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/dto"
)

var RexLockPeriodDays = 21
var RexLockPeriodHours = int64(RexLockPeriodDays) * 24
var RexLockPeriod = time.Duration(RexLockPeriodHours) * time.Hour
var RexLockPeriodMicroseconds = dto.NewMicroseconds(RexLockPeriodHours)

func CalculateMinFinalCycleTime(finalCycleBufferPeriodDays uint32) time.Duration {
	finalCycleBufferPeriod := time.Duration(finalCycleBufferPeriodDays*24) * time.Hour
	return finalCycleBufferPeriod + RexLockPeriod
}

type Stake struct {
	RoundID              uint64            `json:"pool_id"`
	InitialStake         eos.Asset         `json:"initial_stake"`
	CycleStake           eos.Asset         `json:"cycle_stake"`
	RexBalance           eos.Asset         `json:"rex_balance"`
	MinimumReturn        eos.Asset         `json:"minimum_return"`
	CycleReturn          eos.Asset         `json:"cycle_return"`
	TotalReturn          eos.Asset         `json:"total_return"`
	RexState             eos.Name          `json:"rex_state"`
	Cycle                uint32            `json:"cycle"`
	CycleStakingPeriod   *dto.Microseconds `json:"cycle_staking_period"`
	StakingPeriod        *dto.Microseconds `json:"staking_period"`
	CycleStakedTime      eos.TimePoint     `json:"cycle_staked_time"`
	StakedTime           eos.TimePoint     `json:"staked_time"`
	MovedFromSavingsTime eos.TimePoint     `json:"moved_from_savings_time"`
	MaturityTime         eos.TimePoint     `json:"maturity_time"`
	StakeEndTime         eos.TimePoint     `json:"stake_end_time"`
	UpdatedDate          eos.TimePoint     `json:"updated_date"`
	// NOT USED AT THE MOMENT
	// AdditionalFields types.AdditionalFields `json:"additional_fields"`
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

func (m *Stake) FinalStakeEndTime() eos.TimePoint {
	return eos.TimePoint(m.StakedTime.Time().Add(time.Hour * time.Duration(m.StakingPeriod.Hrs())).UnixMicro())
}

func (m *Stake) IsFinalCycle() bool {
	return !m.StakeEndTime.Time().Before(m.FinalStakeEndTime().Time())
}

func (m *Stake) GetNextStakeEndTime(finalCycleBufferPeriodDays uint32) eos.TimePoint {
	if m.RexState == RexStateInSavings {
		nextStakeEndTime := eos.TimePoint(m.StakedTime.Time().Add(time.Hour * time.Duration(m.StakingPeriod.Hrs())).UnixMicro())
		minTimeForNextPartialCycle := eos.TimePoint(m.CycleStakedTime.Time().Add((time.Hour * time.Duration(m.CycleStakingPeriod.Hrs())) + CalculateMinFinalCycleTime(finalCycleBufferPeriodDays)).UnixMicro())
		if minTimeForNextPartialCycle.Time().Before(nextStakeEndTime.Time()) {
			nextStakeEndTime = eos.TimePoint(m.CycleStakedTime.Time().Add(time.Hour * time.Duration(m.CycleStakingPeriod.Hrs())).UnixMicro())
		}
		return nextStakeEndTime
	}
	return eos.TimePoint(m.MovedFromSavingsTime.Time().Add(RexLockPeriod).UnixMicro())
}

func (m *Stake) SplitReturn() *ReturnSplit {
	withdrawAmount := m.CycleReturn
	if !m.IsFinalCycle() {
		withdrawAmount = m.CycleReturn.Sub(m.InitialStake)
		if withdrawAmount.Amount < m.MinimumReturn.Amount {
			withdrawAmount.Amount = 0
		}
	}
	return &ReturnSplit{
		WithdrawAmount: withdrawAmount,
		StakeAmount:    m.CycleReturn.Sub(withdrawAmount),
	}
}

type ReturnSplit struct {
	WithdrawAmount eos.Asset
	StakeAmount    eos.Asset
}

func (m *ReturnSplit) HasWithdrawAmount() bool {
	return m.WithdrawAmount.Amount > 0
}

func (m *ReturnSplit) HasStakeAmount() bool {
	return m.StakeAmount.Amount > 0
}

func (m *ReturnSplit) String() string {
	result, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Sprintf("Failed marshalling round: %v", err))
	}
	return string(result)
}

func (m *EosRexFlowContract) CheckStakeParameters(authorizer, tokenContract eos.AccountName, minStakeAmount eos.Asset, maxStakeAmount eos.Asset, stakingPeriodHrs uint32) (string, error) {
	actionData := struct {
		TokenContract    eos.AccountName
		MinStakeAmount   eos.Asset
		MaxStakeAmount   eos.Asset
		StakingPeriodHrs uint32
	}{tokenContract, minStakeAmount, maxStakeAmount, stakingPeriodHrs}
	return m.ExecAction(authorizer, "chckstkparam", actionData)
}

func (m *EosRexFlowContract) MoveRoundFromSavings(roundId uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "mvfrmsvngsrn", roundId)
}

func (m *EosRexFlowContract) CalculateProceedsRoundRex(roundId uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "clcproceedrn", roundId)
}

func (m *EosRexFlowContract) WithdrawRoundRex(roundId uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "withdrwrexrn", roundId)
}

func (m *EosRexFlowContract) MoveFromSavings(callCounter uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "mvfrmsavings", callCounter)
}

func (m *EosRexFlowContract) UpdateRex(callCounter uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "updaterex", callCounter)
}

func (m *EosRexFlowContract) CalculateProceeds(callCounter uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "calcproceeds", callCounter)
}

func (m *EosRexFlowContract) SellRex(callCounter uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "sellrex", callCounter)
}

func (m *EosRexFlowContract) WithdrawRex(callCounter uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "withdrawrex", callCounter)
}

func (m *EosRexFlowContract) SetStake(roundId uint64, rexBalance eos.Asset, totalReturn eos.Asset, authorizer interface{}) (string, error) {
	actionData := struct {
		RoundId     uint64
		RexBalance  eos.Asset
		TotalReturn eos.Asset
	}{roundId, rexBalance, totalReturn}
	return m.ExecAction(m.GetValueOrContract(authorizer), "setstake", actionData)
}

func (m *EosRexFlowContract) SetCycleReturn(roundId uint64, cycleReturn eos.Asset, updateState bool, authorizer interface{}) (string, error) {
	actionData := struct {
		RoundId     uint64
		CycleReturn eos.Asset
		UpdateState bool
	}{roundId, cycleReturn, updateState}
	return m.ExecAction(m.GetValueOrContract(authorizer), "setcyclertrn", actionData)
}

func (m *EosRexFlowContract) TstLapseTime(roundId uint64, inTests bool) (string, error) {
	actionData := struct {
		RoundId     uint64
		InTests     bool
		CallCounter uint64
	}{roundId, inTests, m.NextCallCounter()}
	return m.ExecAction(eos.AN(m.ContractName), "tstlapsetime", actionData)
}

func (m *Stake) Clone() *Stake {
	return &Stake{
		RoundID:              m.RoundID,
		InitialStake:         m.InitialStake,
		CycleStake:           m.CycleStake,
		RexBalance:           m.RexBalance,
		MinimumReturn:        m.MinimumReturn,
		CycleReturn:          m.CycleReturn,
		TotalReturn:          m.TotalReturn,
		RexState:             m.RexState,
		Cycle:                m.Cycle,
		CycleStakingPeriod:   m.CycleStakingPeriod,
		StakingPeriod:        m.StakingPeriod,
		CycleStakedTime:      m.CycleStakedTime,
		StakedTime:           m.StakedTime,
		MovedFromSavingsTime: m.MovedFromSavingsTime,
		MaturityTime:         m.MaturityTime,
		StakeEndTime:         m.StakeEndTime,
		UpdatedDate:          m.UpdatedDate,
	}
}

func (m *EosRexFlowContract) GetStake(roundID uint64) (*Stake, error) {
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

func (m *EosRexFlowContract) GetLastStake() (*Stake, error) {
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

func (m *EosRexFlowContract) GetStakesReq(req *eos.GetTableRowsRequest) ([]Stake, error) {

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

func (m *EosRexFlowContract) GetStakesByRexStateAndId(state eos.Name) ([]Stake, error) {
	request := &eos.GetTableRowsRequest{}
	err := m.FilterStakesByRexStateAndId(request, state)
	if err != nil {
		return nil, err
	}
	return m.GetStakesReq(request)
}

func (m *EosRexFlowContract) FilterStakesByRexStateAndId(req *eos.GetTableRowsRequest, state eos.Name) error {

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
