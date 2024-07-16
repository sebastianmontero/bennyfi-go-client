package ibc

import "github.com/sebastianmontero/eos-go"

type StakeParams struct {
	PoolId           uint64    `json:"pool_id"`
	YieldSource      eos.Name  `json:"yield_source"`
	Quantity         eos.Asset `json:"quantity"`
	StakingPeriodHrs uint32    `json:"staking_period_hrs"`
}

func (m *StakeParams) ToEmitStakeParams(tokenContract eos.AccountName) *EmitStakeParams {
	return &EmitStakeParams{
		PoolId:           m.PoolId,
		YieldSource:      m.YieldSource,
		Quantity:         eos.ExtendedAsset{Asset: m.Quantity, Contract: tokenContract},
		StakingPeriodHrs: m.StakingPeriodHrs,
	}
}

type EmitXferParams struct {
	Owner       eos.AccountName   `json:"owner"`
	Quantity    eos.ExtendedAsset `json:"quantity"`
	Beneficiary eos.AccountName   `json:"beneficiary"`
}

type EmitStakeParams struct {
	PoolId           uint64            `json:"pool_id"`
	YieldSource      eos.Name          `json:"yield_source"`
	Quantity         eos.ExtendedAsset `json:"quantity"`
	StakingPeriodHrs uint32            `json:"staking_period_hrs"`
}

type EmitUnstakeParams struct {
	PoolId      uint64            `json:"pool_id"`
	YieldSource eos.Name          `json:"yield_source"`
	TotalReturn eos.ExtendedAsset `json:"total_return"`
}
