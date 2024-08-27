// The MIT License (MIT)

// Copyright (c) 2020, Digital Scarcity

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
package rex

import (
	"fmt"

	"github.com/sebastianmontero/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/contract"
	"github.com/sebastianmontero/eos-go-toolbox/service"
)

var (
	RexSymbol = eos.Symbol{
		Precision: 4,
		Symbol:    "REX",
	}
	RexFundSymbol = eos.Symbol{
		Precision: 4,
		Symbol:    "BTLOS",
	}
	VersionEOS   = eos.Name("eos")
	VersionTELOS = eos.Name("telos")
)

type Config struct {
	TokenContract     eos.AccountName `json:"token_contract"`
	LendableIncrement uint64          `json:"lendable_increment"`
	Version           eos.Name        `json:"version"`
}

type Balance struct {
	Owner           eos.AccountName `json:"owner"`
	FundInBalance   eos.Asset       `json:"fund_in_balance"`
	RexBought       eos.Asset       `json:"rex_bought"`
	RexInSavings    eos.Asset       `json:"rex_in_savings"`
	RexLiquid       eos.Asset       `json:"rex_liquid"`
	RexInSellOrders eos.Asset       `json:"rex_in_sell_orders"`
	FundOutBalance  eos.Asset       `json:"fund_out_balance"`
}

type RexPool struct {
	TotalLent     eos.Asset `json:"total_lent"`
	TotalUnlent   eos.Asset `json:"total_unlent"`
	TotalLendable eos.Asset `json:"total_lendable"`
	TotalRex      eos.Asset `json:"total_rex"`
}

func (m *RexPool) ToInitialPool() *InitialPool {
	return &InitialPool{
		TotalLendable: m.TotalLendable,
		TotalRex:      m.TotalRex,
		DepositCount:  0,
	}
}

type InitialPool struct {
	TotalLendable eos.Asset `json:"total_lendable"`
	TotalRex      eos.Asset `json:"total_rex"`
	DepositCount  uint64    `json:"deposit_count"`
}

type PairTimePointSecInt64 struct {
	First  eos.TimePointSec `json:"first"`
	Second int64            `json:"second"`
}

type RexBalance struct {
	Version       uint8                    `json:"version"`
	Owner         eos.AccountName          `json:"owner"`
	VoteStake     eos.Asset                `json:"vote_stake"`
	RexBalance    eos.Asset                `json:"rex_balance"`
	MaturedRex    int64                    `json:"matured_rex"`
	RexMaturities []*PairTimePointSecInt64 `json:"rex_maturities"`
}

type RexOrder struct {
	Version      uint8           `json:"version"`
	Owner        eos.AccountName `json:"owner"`
	RexRequested eos.Asset       `json:"rex_requested"`
	Proceeds     eos.Asset       `json:"proceeds"`
	StakeChange  eos.Asset       `json:"stake_change"`
	OrderTime    eos.TimePoint   `json:"order_time"`
	IsOpen       uint8           `json:"is_open"`
}

type RexFund struct {
	Version uint8           `json:"version"`
	Owner   eos.AccountName `json:"owner"`
	Balance eos.Asset       `json:"balance"`
}

type SetInitialPoolArgs struct {
	TotalLendable eos.Asset `json:"total_lendable"`
	TotalRex      eos.Asset `json:"total_rex"`
}

type InitRexArgs struct {
	TotalLendable     eos.Asset       `json:"total_lendable"`
	TotalRex          eos.Asset       `json:"total_rex"`
	LendableIncrement uint64          `json:"lendable_increment"`
	TokenContract     eos.AccountName `json:"token_contract"`
	Version           eos.Name        `json:"version"`
}

type SetLentArgs struct {
	TotalLent   eos.Asset `json:"total_lent"`
	TotalUnlent eos.Asset `json:"total_unlent"`
}

type InitConfArgs struct {
	LendableIncrement uint64          `json:"lendable_increment"`
	TokenContract     eos.AccountName `json:"token_contract"`
	Version           eos.Name        `json:"version"`
}

type RexContract struct {
	*contract.Contract
}

func NewRexContract(eos *service.EOS, contractName string) *RexContract {
	return &RexContract{
		&contract.Contract{
			EOS:          eos,
			ContractName: contractName,
		},
	}
}

func (m *RexContract) ExecAction(permissionLevel interface{}, action string, actionData interface{}) (string, error) {
	resp, err := m.Contract.ExecAction(permissionLevel, action, actionData)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Tx ID: %v", resp.TransactionID), nil
}

func (m *RexContract) Init(totalLendable, totalRex eos.Asset, lendableIncrement uint64, tokenContract eos.AccountName, version eos.Name) (string, error) {
	actionData := &InitRexArgs{
		TotalLendable:     totalLendable,
		TotalRex:          totalRex,
		LendableIncrement: lendableIncrement,
		TokenContract:     tokenContract,
		Version:           version,
	}

	return m.ExecAction(m.ContractName, "init", actionData)
}

func (m *RexContract) SetInitialPool(totalLendable, totalRex eos.Asset) (string, error) {
	actionData := &SetInitialPoolArgs{
		TotalLendable: totalLendable,
		TotalRex:      totalRex,
	}
	return m.ExecAction(m.ContractName, "setinitpool", actionData)
}

func (m *RexContract) SetLent(totalLent, totalUnlent eos.Asset) (string, error) {
	actionData := &SetLentArgs{
		TotalLent:   totalLent,
		TotalUnlent: totalUnlent,
	}
	return m.ExecAction(m.ContractName, "setlent", actionData)
}

func (m *RexContract) InitConf(lendableIncrement uint64, tokenContract eos.AccountName, version eos.Name) (string, error) {
	actionData := &InitConfArgs{
		LendableIncrement: lendableIncrement,
		TokenContract:     tokenContract,
		Version:           version,
	}

	return m.ExecAction(m.ContractName, "initconf", actionData)
}

func (m *RexContract) SetIncrement(lendableIncrement uint64) (string, error) {
	return m.ExecAction(m.ContractName, "setincrement", lendableIncrement)
}

func (m *RexContract) SetVersion(version eos.Name) (string, error) {
	actionData := struct {
		Version eos.Name
	}{version}
	return m.ExecAction(m.ContractName, "setversion", actionData)
}

func (m *RexContract) SetFund(owner eos.AccountName, balance eos.Asset) (string, error) {
	actionData := struct {
		Owner   eos.AccountName
		Balance eos.Asset
	}{owner, balance}
	return m.ExecAction(m.ContractName, "setfund", actionData)
}

func (m *RexContract) EraseFund(owner eos.AccountName) (string, error) {
	actionData := struct {
		Owner eos.AccountName
	}{owner}
	return m.ExecAction(m.ContractName, "erasefund", actionData)
}

func (m *RexContract) SetRexBalance(owner eos.AccountName, rexBalance eos.Asset, maturedRex int64) (string, error) {
	actionData := struct {
		Owner      eos.AccountName
		RexBalance eos.Asset
		MaturedRex int64
	}{owner, rexBalance, maturedRex}
	return m.ExecAction(m.ContractName, "setrexbal", actionData)
}

func (m *RexContract) EraseRexBalance(owner eos.AccountName) (string, error) {
	actionData := struct {
		Owner eos.AccountName
	}{owner}
	return m.ExecAction(m.ContractName, "eraserexbal", actionData)
}

func (m *RexContract) SetOrder(owner eos.AccountName, rexRequested eos.Asset, proceeds eos.Asset, isOpen uint8) (string, error) {
	actionData := struct {
		Owner      eos.AccountName
		RexRequest eos.Asset
		Proceeds   eos.Asset
		IsOpen     uint8
	}{owner, rexRequested, proceeds, isOpen}
	return m.ExecAction(m.ContractName, "setorder", actionData)
}

func (m *RexContract) EraseOrder(owner eos.AccountName) (string, error) {
	actionData := struct {
		Owner eos.AccountName
	}{owner}
	return m.ExecAction(m.ContractName, "eraseorder", actionData)
}

func (m *RexContract) UpdateRex(owner eos.AccountName) (string, error) {
	actionData := struct {
		Owner eos.AccountName
	}{owner}
	return m.ExecAction(m.ContractName, "updaterex", actionData)
}

func (m *RexContract) LapseMaturities(owner eos.AccountName, numDays uint32, processMaturities bool, callCounter uint64) (string, error) {
	actionData := struct {
		Owner             eos.AccountName
		NumDays           uint32
		ProcessMaturities bool
		CallCnt           uint64
	}{owner, numDays, processMaturities, callCounter}
	return m.ExecAction(m.ContractName, "lapsematrts", actionData)
}

func (m *RexContract) ResetConf() (string, error) {
	return m.ExecAction(m.ContractName, "resetconf", nil)
}

func (m *RexContract) ResetBalance() (string, error) {
	return m.ExecAction(m.ContractName, "resetbal", nil)
}

func (m *RexContract) ResetPool() (string, error) {
	return m.ExecAction(m.ContractName, "resetpool", nil)
}

func (m *RexContract) ResetInitialPool() (string, error) {
	return m.ExecAction(m.ContractName, "rsetinitpool", nil)
}

func (m *RexContract) Deposit(owner eos.AccountName, amount eos.Asset) (string, error) {
	actionData := struct {
		Owner  eos.AccountName
		Amount eos.Asset
	}{owner, amount}
	return m.ExecAction(owner, "deposit", actionData)
}

func (m *RexContract) BuyRex(from eos.AccountName, amount eos.Asset) (string, error) {
	actionData := struct {
		From   eos.AccountName
		Amount eos.Asset
	}{from, amount}

	return m.ExecAction(from, "buyrex", actionData)
}

func (m *RexContract) MoveToSavings(owner eos.AccountName, rex eos.Asset) (string, error) {
	actionData := struct {
		Owner eos.AccountName
		Rex   eos.Asset
	}{owner, rex}

	return m.ExecAction(owner, "mvtosavings", actionData)
}

func (m *RexContract) MoveFromSavings(owner eos.AccountName, rex eos.Asset) (string, error) {
	actionData := struct {
		Owner eos.AccountName
		Rex   eos.Asset
	}{owner, rex}

	return m.ExecAction(owner, "mvfrsavings", actionData)
}

func (m *RexContract) SellRex(from eos.AccountName, rex eos.Asset) (string, error) {
	actionData := struct {
		From eos.AccountName
		Rex  eos.Asset
	}{from, rex}
	return m.ExecAction(from, "sellrex", actionData)
}

func (m *RexContract) Withdraw(owner eos.AccountName, amount eos.Asset) (string, error) {
	actionData := struct {
		Owner  eos.AccountName
		Amount eos.Asset
	}{owner, amount}

	return m.ExecAction(owner, "withdraw", actionData)
}

func (m *RexContract) GetConfig() (*Config, error) {
	var config []Config
	req := &eos.GetTableRowsRequest{
		Table: "config",
	}
	err := m.GetTableRows(*req, &config)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	if len(config) > 0 {
		return &config[0], nil
	}
	return nil, nil
}

func (m *RexContract) GetBalance(owner eos.AccountName) (*Balance, error) {
	entries, err := m.GetBalancesReq(&eos.GetTableRowsRequest{
		LowerBound: string(owner),
		UpperBound: string(owner),
		Limit:      1,
	})
	if err != nil {
		return nil, err
	}
	if len(entries) > 0 {
		return &entries[0], nil
	}
	return nil, nil
}

func (m *RexContract) GetBalancesReq(req *eos.GetTableRowsRequest) ([]Balance, error) {

	var balances []Balance
	if req == nil {
		req = &eos.GetTableRowsRequest{}
	}
	req.Table = "balance"
	err := m.GetTableRows(*req, &balances)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	return balances, nil
}

func (m *RexContract) GetPool() (*RexPool, error) {
	var pool []RexPool
	req := &eos.GetTableRowsRequest{
		Table: "rexpool",
	}
	err := m.GetTableRows(*req, &pool)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	if len(pool) > 0 {
		return &pool[0], nil
	}
	return nil, nil
}

func (m *RexContract) GetInitialPool() (*InitialPool, error) {
	var pool []InitialPool
	req := &eos.GetTableRowsRequest{
		Table: "initialpool",
	}
	err := m.GetTableRows(*req, &pool)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	if len(pool) > 0 {
		return &pool[0], nil
	}
	return nil, nil
}

func (m *RexContract) GetRexBalance(owner eos.AccountName) (*RexBalance, error) {
	entries, err := m.GetRexBalancesReq(&eos.GetTableRowsRequest{
		LowerBound: string(owner),
		UpperBound: string(owner),
		Limit:      1,
	})
	if err != nil {
		return nil, err
	}
	if len(entries) > 0 {
		return &entries[0], nil
	}
	return nil, nil
}

func (m *RexContract) GetRexBalancesReq(req *eos.GetTableRowsRequest) ([]RexBalance, error) {

	var balances []RexBalance
	if req == nil {
		req = &eos.GetTableRowsRequest{}
	}
	req.Table = "rexbal"
	err := m.GetTableRows(*req, &balances)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	return balances, nil
}

func (m *RexContract) GetOrder(owner eos.AccountName) (*RexOrder, error) {
	entries, err := m.GetOrdersReq(&eos.GetTableRowsRequest{
		LowerBound: string(owner),
		UpperBound: string(owner),
		Limit:      1,
	})
	if err != nil {
		return nil, err
	}
	if len(entries) > 0 {
		return &entries[0], nil
	}
	return nil, nil
}

func (m *RexContract) GetOrdersReq(req *eos.GetTableRowsRequest) ([]RexOrder, error) {

	var orders []RexOrder
	if req == nil {
		req = &eos.GetTableRowsRequest{}
	}
	req.Table = "rexqueue"
	err := m.GetTableRows(*req, &orders)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	return orders, nil
}

func (m *RexContract) GetFund(owner eos.AccountName) (*RexFund, error) {
	entries, err := m.GetFundsReq(&eos.GetTableRowsRequest{
		LowerBound: string(owner),
		UpperBound: string(owner),
		Limit:      1,
	})
	if err != nil {
		return nil, err
	}
	if len(entries) > 0 {
		return &entries[0], nil
	}
	return nil, nil
}

func (m *RexContract) GetFundsReq(req *eos.GetTableRowsRequest) ([]RexFund, error) {

	var funds []RexFund
	if req == nil {
		req = &eos.GetTableRowsRequest{}
	}
	req.Table = "rexfund"
	err := m.GetTableRows(*req, &funds)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	return funds, nil
}

func (m *RexContract) GetMinUnlent(totalLent, proceedsAmount eos.Asset) eos.Asset {
	return eos.Asset{Amount: totalLent.Amount/10 + proceedsAmount.Amount, Symbol: totalLent.Symbol}
}
