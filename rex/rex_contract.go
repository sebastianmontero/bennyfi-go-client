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

	"github.com/eoscanada/eos-go"
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
)

type Config struct {
	TokenContract     eos.Name `json:"token_contract"`
	LendableIncrement uint64   `json:"lendable_increment"`
}

type Balance struct {
	Owner           eos.AccountName `json:"owner"`
	FundInBalance   string          `json:"fund_in_balance"`
	RexBought       string          `json:"rex_bought"`
	RexInSavings    string          `json:"rex_in_savings"`
	RexLiquid       string          `json:"rex_liquid"`
	RexInSellOrders string          `json:"rex_in_sell_orders"`
	FundOutBalance  string          `json:"fund_out_balance"`
}

func (m *Balance) GetFundOutBalance() eos.Asset {
	fundOutBalance, err := eos.NewAssetFromString(m.FundOutBalance)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse entry stake: %v to asset", m.FundOutBalance))
	}
	return fundOutBalance
}

type RexPool struct {
	TotalLendable string `json:"total_lendable"`
	TotalRex      string `json:"total_rex"`
}

func (m *RexPool) GetTotalLendable() eos.Asset {
	totalLendable, err := eos.NewAssetFromString(m.TotalLendable)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse entry stake: %v to asset", m.TotalLendable))
	}
	return totalLendable
}

func (m *RexPool) GetTotalRex() eos.Asset {
	totalRex, err := eos.NewAssetFromString(m.TotalRex)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse entry stake: %v to asset", m.TotalRex))
	}
	return totalRex
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

func (m *RexContract) Init(totalLendable, totalRex eos.Asset, lendableIncrement uint64, tokenContract eos.Name) (string, error) {
	actionData := make(map[string]interface{})
	actionData["total_lendable"] = totalLendable
	actionData["total_rex"] = totalRex
	actionData["lendable_increment"] = lendableIncrement
	actionData["token_contract"] = tokenContract

	return m.ExecAction(m.ContractName, "init", actionData)
}

func (m *RexContract) InitConf(lendableIncrement uint64, tokenContract eos.Name) (string, error) {
	actionData := make(map[string]interface{})
	actionData["lendable_increment"] = lendableIncrement
	actionData["token_contract"] = tokenContract

	return m.ExecAction(m.ContractName, "initconf", actionData)
}

func (m *RexContract) SetIncrement(lendableIncrement uint64) (string, error) {
	actionData := make(map[string]interface{})
	actionData["lendable_increment"] = lendableIncrement

	return m.ExecAction(m.ContractName, "setincrement", actionData)
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

func (m *RexContract) Deposit(owner eos.AccountName, amount eos.Asset) (string, error) {
	actionData := make(map[string]interface{})
	actionData["owner"] = owner
	actionData["amount"] = amount

	return m.ExecAction(owner, "deposit", actionData)
}

func (m *RexContract) BuyRex(from eos.AccountName, amount eos.Asset) (string, error) {
	actionData := make(map[string]interface{})
	actionData["from"] = from
	actionData["amount"] = amount

	return m.ExecAction(from, "buyrex", actionData)
}

func (m *RexContract) MoveToSavings(owner eos.AccountName, rex eos.Asset) (string, error) {
	actionData := make(map[string]interface{})
	actionData["owner"] = owner
	actionData["rex"] = rex

	return m.ExecAction(owner, "mvtosavings", actionData)
}

func (m *RexContract) MoveFromSavings(owner eos.AccountName, rex eos.Asset) (string, error) {
	actionData := make(map[string]interface{})
	actionData["owner"] = owner
	actionData["rex"] = rex

	return m.ExecAction(owner, "mvfrsavings", actionData)
}

func (m *RexContract) SellRex(from eos.AccountName, rex eos.Asset) (string, error) {
	actionData := make(map[string]interface{})
	actionData["from"] = from
	actionData["rex"] = rex

	return m.ExecAction(from, "sellrex", actionData)
}

func (m *RexContract) Withdraw(owner eos.AccountName, amount eos.Asset) (string, error) {
	actionData := make(map[string]interface{})
	actionData["owner"] = owner
	actionData["amount"] = amount

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

func (m *RexContract) GetBalance(owner eos.Name) (*Balance, error) {
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
