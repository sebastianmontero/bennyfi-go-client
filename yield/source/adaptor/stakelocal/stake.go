package stakelocal

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/sebastianmontero/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/dto"
)

type Stake struct {
	RoundID       uint64            `json:"pool_id"`
	TotalStake    eos.Asset         `json:"total_stake"`
	TotalReturn   eos.Asset         `json:"total_return"`
	State         eos.Name          `json:"state"`
	StakingPeriod *dto.Microseconds `json:"staking_period"`
	StakedTime    eos.TimePoint     `json:"staked_time"`
	StakeEndTime  eos.TimePoint     `json:"stake_end_time"`
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

func (m *StakeLocalContract) CheckStakeParameters(authorizer, tokenContract eos.AccountName, stakeAmount eos.Asset, stakingPeriodHrs uint32, yieldSourceName eos.Name) (string, error) {
	actionData := struct {
		TokenContract    eos.AccountName
		StakeAmount      eos.Asset
		StakingPeriodHrs uint32
		YieldSourceName  eos.Name
	}{tokenContract, stakeAmount, stakingPeriodHrs, yieldSourceName}
	return m.ExecAction(authorizer, "chckstkparam", actionData)
}

func (m *Stake) Clone() *Stake {
	return &Stake{
		RoundID:       m.RoundID,
		TotalStake:    m.TotalStake,
		TotalReturn:   m.TotalReturn,
		State:         m.State,
		StakingPeriod: m.StakingPeriod,
		StakedTime:    m.StakedTime,
		StakeEndTime:  m.StakeEndTime,
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
