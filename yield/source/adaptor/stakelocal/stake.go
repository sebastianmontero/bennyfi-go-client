package stakelocal

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/sebastianmontero/bennyfi-go-client/common/types"
	"github.com/sebastianmontero/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/dto"
)

type Stake struct {
	RoundID          uint64                 `json:"pool_id"`
	TotalStake       eos.Asset              `json:"total_stake"`
	TotalReturn      eos.Asset              `json:"total_return"`
	LastCycleReturn  eos.Asset              `json:"last_cycle_return"`
	State            eos.Name               `json:"state"`
	Cycle            uint32                 `json:"cycle"`
	StakingPeriod    *dto.Microseconds      `json:"staking_period"`
	StakedTime       eos.TimePoint          `json:"staked_time"`
	CycleStakedTime  eos.TimePoint          `json:"cycle_staked_time"`
	StakeEndTime     eos.TimePoint          `json:"stake_end_time"`
	AdditionalFields types.AdditionalFields `json:"additional_fields"`
}

func (m *Stake) String() string {
	result, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Sprintf("Failed marshalling round: %v", err))
	}
	return string(result)
}

func (m *StakeLocalContract) CheckStakeParameters(authorizer, tokenContract eos.AccountName, minStakeAmount eos.Asset, maxStakeAmount eos.Asset, stakingPeriodHrs uint32, yieldSourceName eos.Name) (string, error) {
	actionData := struct {
		TokenContract    eos.AccountName
		MinStakeAmount   eos.Asset
		MaxStakeAmount   eos.Asset
		StakingPeriodHrs uint32
		YieldSourceName  eos.Name
	}{tokenContract, minStakeAmount, maxStakeAmount, stakingPeriodHrs, yieldSourceName}
	return m.ExecAction(authorizer, "chckstkparam", actionData)
}

func (m *Stake) Clone() *Stake {
	return &Stake{
		RoundID:          m.RoundID,
		TotalStake:       m.TotalStake,
		TotalReturn:      m.TotalReturn,
		LastCycleReturn:  m.LastCycleReturn,
		State:            m.State,
		Cycle:            m.Cycle,
		StakingPeriod:    m.StakingPeriod,
		StakedTime:       m.StakedTime,
		StakeEndTime:     m.StakeEndTime,
		AdditionalFields: m.AdditionalFields,
	}
}

func (m *StakeLocalContract) GetStake(roundID uint64) (*Stake, error) {
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

func (m *StakeLocalContract) GetLastStake() (*Stake, error) {
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

func (m *StakeLocalContract) GetStakesReq(req *eos.GetTableRowsRequest) ([]Stake, error) {

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

func (m *StakeLocalContract) GetAllStakesByState(state eos.Name) ([]Stake, error) {
	allStakes := []Stake{}
	startPoolId := uint64(0)
	for {
		stakes, err := m.GetStakesByStateAndId(state, startPoolId)
		if err != nil {
			return nil, err
		}
		if len(stakes) == 0 {
			break
		}
		allStakes = append(allStakes, stakes...)
		startPoolId = stakes[len(stakes)-1].RoundID + 1
	}
	return allStakes, nil
}

func (m *StakeLocalContract) GetStakesByStateAndId(state eos.Name, startPoolId uint64) ([]Stake, error) {
	request := &eos.GetTableRowsRequest{}
	err := m.FilterStakesByStateAndId(request, state, startPoolId)
	if err != nil {
		return nil, err
	}
	return m.GetStakesReq(request)
}

func (m *StakeLocalContract) FilterStakesByStateAndId(req *eos.GetTableRowsRequest, state eos.Name, startPoolId uint64) error {

	req.Index = "2"
	req.KeyType = "i128"
	req.Reverse = false
	stateAndRndLB, err := m.EOS.GetComposedIndexValue(state, startPoolId)
	if err != nil {
		return fmt.Errorf("failed to generate lower bound composed index, err: %v", err)
	}
	stateAndRndUB, err := m.EOS.GetComposedIndexValue(state, uint64(18446744073709551615))
	if err != nil {
		return fmt.Errorf("failed to generate upper bound composed index, err: %v", err)
	}
	// fmt.Println("LB: ", stateAndRndLB, "UB: ", stateAndRndUB)
	req.LowerBound = stateAndRndLB
	req.UpperBound = stateAndRndUB
	return err
}

func (m *StakeLocalContract) GetStakesByYieldSourceAndId(yieldSource eos.Name, startPoolId uint64) ([]Stake, error) {
	request := &eos.GetTableRowsRequest{}
	err := m.FilterStakesByYieldSourceAndId(request, yieldSource, startPoolId)
	if err != nil {
		return nil, err
	}
	return m.GetStakesReq(request)
}

func (m *StakeLocalContract) FilterStakesByYieldSourceAndId(req *eos.GetTableRowsRequest, yieldSource eos.Name, startPoolId uint64) error {

	req.Index = "2"
	req.KeyType = "i128"
	req.Reverse = false
	stateAndRndLB, err := m.EOS.GetComposedIndexValue(yieldSource, startPoolId)
	if err != nil {
		return fmt.Errorf("failed to generate lower bound composed index, err: %v", err)
	}
	stateAndRndUB, err := m.EOS.GetComposedIndexValue(yieldSource, uint64(18446744073709551615))
	if err != nil {
		return fmt.Errorf("failed to generate upper bound composed index, err: %v", err)
	}
	// fmt.Println("LB: ", stateAndRndLB, "UB: ", stateAndRndUB)
	req.LowerBound = stateAndRndLB
	req.UpperBound = stateAndRndUB
	return err
}
