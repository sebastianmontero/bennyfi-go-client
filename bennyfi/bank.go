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
package bennyfi

import (
	"fmt"

	eos "github.com/eoscanada/eos-go"
)

type Balance struct {
	ID            uint64          `json:"id"`
	TokenHolder   eos.AccountName `json:"token_holder"`
	Symbol        string          `json:"symbol"`
	LiquidBalance string          `json:"liquid_balance"`
	StakedBalance string          `json:"staked_balance"`
	TokenContract eos.AccountName `json:"token_contract"`
}

func (m *Balance) GetLiquidBalance() eos.Asset {
	liquidBalance, err := eos.NewAssetFromString(m.LiquidBalance)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse liquid balance: %v to asset", m.LiquidBalance))
	}
	return liquidBalance
}

func (m *Balance) GetStakedBalance() eos.Asset {
	stakedBalance, err := eos.NewAssetFromString(m.StakedBalance)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse staked balance: %v to asset", m.StakedBalance))
	}
	return stakedBalance
}

func (m *Balance) HasLiquidBalance() bool {
	return m.GetLiquidBalance().Amount > 0
}

func (m *Balance) HasStakedBalance() bool {
	return m.GetStakedBalance().Amount > 0
}

func (m *Balance) HasBalance() bool {
	return m.HasLiquidBalance() || m.HasStakedBalance()
}

func (m *BennyfiContract) Withdraw(from eos.AccountName, quantity eos.Asset) (string, error) {
	actionData := make(map[string]interface{})
	actionData["from"] = eos.Name(from)
	actionData["quantity"] = quantity

	return m.ExecAction(from, "withdraw", actionData)
}

func (m *BennyfiContract) WithdrawTot(from eos.AccountName, symbol eos.Symbol) (string, error) {
	actionData := make(map[string]interface{})
	actionData["from"] = eos.Name(from)
	actionData["symbol"] = symbol.String()

	return m.ExecAction(from, "withdrawtot", actionData)
}

func (m *BennyfiContract) GetBalances() ([]Balance, error) {
	return m.GetBalancesReq(nil)
}

func (m *BennyfiContract) GetBalancesReq(req *eos.GetTableRowsRequest) ([]Balance, error) {

	var balances []Balance
	if req == nil {
		req = &eos.GetTableRowsRequest{}
	}
	req.Table = "balances"
	err := m.GetTableRows(*req, &balances)
	if err != nil {
		return []Balance{}, fmt.Errorf("get table rows %v", err)
	}
	return balances, nil
}

func (m *BennyfiContract) GetBalancesByAccount(account interface{}) ([]Balance, error) {
	request := &eos.GetTableRowsRequest{}
	m.FilterBalancesbyAccount(request, account)
	return m.GetBalancesReq(request)
}

func (m *BennyfiContract) FilterBalancesbyAccount(req *eos.GetTableRowsRequest, account interface{}) {

	req.Index = "3"
	req.KeyType = "name"
	req.LowerBound = fmt.Sprintf("%v", account)
	req.UpperBound = fmt.Sprintf("%v", account)
}

func (m *BennyfiContract) GetBalance(tokenHolder eos.AccountName, symbol string) (*Balance, error) {

	var balances []Balance
	request := eos.GetTableRowsRequest{
		Table:      "balances",
		Index:      "3",
		KeyType:    "name",
		LowerBound: string(tokenHolder),
		UpperBound: string(tokenHolder),
		//Limit:      1,
	}

	err := m.GetTableRows(request, &balances)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	for _, balance := range balances {
		if balance.Symbol == symbol {
			return &balance, nil
		}
	}
	return nil, nil

}
