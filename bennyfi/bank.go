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

	eos "github.com/sebastianmontero/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/util"
)

type Balance struct {
	ID            uint64          `json:"id"`
	TokenHolder   eos.AccountName `json:"token_holder"`
	Symbol        eos.Symbol      `json:"symbol"`
	LiquidBalance eos.Asset       `json:"liquid_balance"`
	StakedBalance eos.Asset       `json:"staked_balance"`
	TokenContract eos.AccountName `json:"token_contract"`
	// NOT USED AT THE MOMENT
	// AdditionalFields types.AdditionalFields `json:"additional_fields"`
}

type WithdrawArgs struct {
	From     eos.AccountName `json:"from"`
	Quantity eos.Asset       `json:"quantity"`
}

type WithdrawTotalArgs struct {
	From   eos.AccountName `json:"from"`
	Symbol eos.Symbol      `json:"symbol"`
}

func (m *Balance) GetTotalBalance() eos.Asset {
	return m.LiquidBalance.Add(m.StakedBalance)
}

func (m *Balance) HasLiquidBalance() bool {
	return m.LiquidBalance.Amount > 0
}

func (m *Balance) HasStakedBalance() bool {
	return m.StakedBalance.Amount > 0
}

func (m *Balance) HasBalance() bool {
	return m.HasLiquidBalance() || m.HasStakedBalance()
}

func (m *Balance) AddLiquidBalance(amount interface{}, negative bool) eos.Asset {
	amnt, err := util.ToAsset(amount)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse amount: %v to asset", amount))
	}
	if negative {
		m.LiquidBalance = m.LiquidBalance.Sub(amnt)
	} else {
		m.LiquidBalance = m.LiquidBalance.Add(amnt)
	}
	return m.LiquidBalance
}

func (m *Balance) AddStakedBalance(amount interface{}, negative bool) eos.Asset {
	amnt, err := util.ToAsset(amount)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse amount: %v to asset", amount))
	}

	if negative {
		m.StakedBalance = m.StakedBalance.Sub(amnt)
	} else {
		m.StakedBalance = m.StakedBalance.Add(amnt)
	}
	return m.StakedBalance
}

func (m *Balance) Stake(amount interface{}) {
	amnt, err := util.ToAsset(amount)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse amount: %v to asset", amount))
	}

	if m.LiquidBalance.Amount < amnt.Amount {
		panic(fmt.Sprintf("Stake amount: %v can not be greater than liquid balance: %v", amnt, m.LiquidBalance))
	}
	m.AddLiquidBalance(amount, true)
	m.AddStakedBalance(amount, false)

}

func (m *Balance) Unstake(amount interface{}) {
	amnt, err := util.ToAsset(amount)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse amount: %v to asset", amount))
	}

	if m.StakedBalance.Amount < amnt.Amount {
		panic(fmt.Sprintf("Unstake amount: %v can not be greater than staked balance: %v", amnt, m.StakedBalance))
	}
	m.AddLiquidBalance(amount, false)
	m.AddStakedBalance(amount, true)

}

func (m *BennyfiContract) Pay(authorizer eos.AccountName, from interface{}, to interface{}, amount eos.Asset, fromEscrow bool) (string, error) {
	f, err := util.ToAccountName(from)
	if err != nil {
		return "", fmt.Errorf("failed parsing pay from account: %v, error: %v", from, err)
	}

	t, err := util.ToAccountName(to)
	if err != nil {
		return "", fmt.Errorf("failed parsing pay to account: %v, error: %v", to, err)
	}
	actionData := struct {
		From       eos.AccountName
		To         eos.AccountName
		Amount     eos.Asset
		FromEscrow bool
	}{f, t, amount, fromEscrow}

	return m.ExecAction(authorizer, "pay", actionData)
}

func (m *BennyfiContract) Escrow(authorizer eos.AccountName, account interface{}, amount eos.Asset) (string, error) {
	a, err := util.ToAccountName(account)
	if err != nil {
		return "", fmt.Errorf("failed parsing escrow account: %v, error: %v", account, err)
	}
	actionData := struct {
		Account eos.AccountName
		Amount  eos.Asset
	}{a, amount}

	return m.ExecAction(authorizer, "escrow", actionData)
}

func (m *BennyfiContract) Unescrow(authorizer eos.AccountName, account interface{}, amount eos.Asset) (string, error) {
	a, err := util.ToAccountName(account)
	if err != nil {
		return "", fmt.Errorf("failed parsing unescrow account: %v, error: %v", account, err)
	}
	actionData := struct {
		Account eos.AccountName
		Amount  eos.Asset
	}{a, amount}

	return m.ExecAction(authorizer, "unescrow", actionData)
}

func (m *BennyfiContract) Withdraw(from eos.AccountName, quantity eos.Asset) (string, error) {
	actionData := &WithdrawArgs{
		From:     from,
		Quantity: quantity,
	}

	return m.ExecAction(from, "withdraw", actionData)
}

func (m *BennyfiContract) WithdrawTot(from eos.AccountName, symbol eos.Symbol) (string, error) {
	actionData := &WithdrawTotalArgs{
		From:   from,
		Symbol: symbol,
	}

	return m.ExecAction(from, "withdrawtot", actionData)
}

func (m *BennyfiContract) GetBalances() ([]Balance, error) {
	return m.GetBalancesReq(nil)
}

func (m *BennyfiContract) GetAllBalances() ([]Balance, error) {
	var balances []Balance
	req := eos.GetTableRowsRequest{
		Table: "balances",
	}
	err := m.GetAllTableRows(req, "id", &balances)
	if err != nil {
		return nil, fmt.Errorf("failed getting all balances: %v", err)
	}
	return balances, nil
}

func (m *BennyfiContract) GetAllBalancesAsMap() ([]map[string]interface{}, error) {
	req := eos.GetTableRowsRequest{
		Table: "balances",
	}
	return m.GetAllTableRowsAsMap(req, "id")
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

func (m *BennyfiContract) GetBalance(tokenHolder eos.AccountName, symbol interface{}) (*Balance, error) {

	symb, err := util.ToSymbol(symbol)
	if err != nil {
		return nil, err
	}

	var balances []Balance
	request := eos.GetTableRowsRequest{
		Table:      "balances",
		Index:      "3",
		KeyType:    "name",
		LowerBound: string(tokenHolder),
		UpperBound: string(tokenHolder),
		//Limit:      1,
	}

	err = m.GetTableRows(request, &balances)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	for _, balance := range balances {
		if balance.Symbol.Equal(symb) {
			return &balance, nil
		}
	}
	return nil, nil

}

func (m *BennyfiContract) GetBalanceOrDefault(tokenHolder eos.AccountName, symbol, tokenContract interface{}) (*Balance, error) {

	symb, err := util.ToSymbol(symbol)
	if err != nil {
		return nil, err
	}

	tkc, err := util.ToAccountName(tokenContract)
	if err != nil {
		return nil, err
	}

	balance, err := m.GetBalance(tokenHolder, symb.String())
	if err != nil {
		return nil, err
	}
	if balance == nil {
		balance = &Balance{
			TokenHolder:   tokenHolder,
			Symbol:        symb,
			LiquidBalance: eos.Asset{Amount: 0, Symbol: symb},
			StakedBalance: eos.Asset{Amount: 0, Symbol: symb},
			TokenContract: tkc,
		}
	}
	return balance, nil

}
