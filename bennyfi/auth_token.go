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
	"encoding/json"
	"fmt"

	eos "github.com/sebastianmontero/eos-go"
)

var (
	TokenRoleFee    = "fee"
	TokenRoleStake  = "stake"
	TokenRoleReward = "reward"
	TokenRoleValues = map[string]bool{
		TokenRoleFee:    true,
		TokenRoleReward: true,
		TokenRoleStake:  true,
	}
)

type TokenLimits struct {
	MinValue eos.Asset `json:"min_value"`
	MaxValue eos.Asset `json:"max_value"`
}

type TokenRole struct {
	Key   eos.Name     `json:"first"`
	Value *TokenLimits `json:"second"`
}

func NewTokenRole(tokenRole string, minValue, maxValue eos.Int64, symbol eos.Symbol) *TokenRole {
	return &TokenRole{
		Key: eos.Name(tokenRole),
		Value: &TokenLimits{
			MinValue: eos.Asset{Amount: minValue, Symbol: symbol},
			MaxValue: eos.Asset{Amount: maxValue, Symbol: symbol},
		},
	}
}

type TokenRoles []*TokenRole

type AuthToken struct {
	Authorizer    eos.AccountName `json:"authorizer"`
	Symbol        string          `json:"symbol"`
	TokenContract eos.AccountName `json:"token_contract"`
	ArtifactCID   string          `json:"artifact_cid"`
	TokenRoles    TokenRoles      `json:"token_roles"`
	// NOT USED AT THE MOMENT
	// AdditionalFields types.AdditionalFields `json:"additional_fields"`
}

func (m *AuthToken) ToSetTokenArgs() *SetTokenArgs {
	symb, err := eos.StringToSymbol(m.Symbol)
	if err != nil {
		panic(fmt.Sprintf("failed parsing symbol: %v, error: %v", m.Symbol, err))
	}
	return &SetTokenArgs{
		Authorizer:    m.Authorizer,
		Symbol:        symb,
		TokenContract: m.TokenContract,
		ArtifactCID:   m.ArtifactCID,
		TokenRoles:    m.TokenRoles,
	}
}

type SetTokenArgs struct {
	Authorizer    eos.AccountName `json:"authorizer"`
	Symbol        eos.Symbol      `json:"symbol"`
	TokenContract eos.AccountName `json:"token_contract"`
	ArtifactCID   string          `json:"artifact_cid"`
	TokenRoles    TokenRoles      `json:"token_roles"`
}

func (m *AuthToken) String() string {
	j, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Sprintf("Failed marshalling auth token, error: %v", err))
	}
	return string(j)

}

type SetTokenRoleArgs struct {
	Authorizer    eos.AccountName `json:"authorizer"`
	Symbol        eos.Symbol      `json:"symbol"`
	TokenContract eos.AccountName `json:"token_contract"`
	TokenRole     eos.Name        `json:"token_role"`
	MinValue      eos.Asset       `json:"min_value"`
	MaxValue      eos.Asset       `json:"max_value"`
}

func (m *SetTokenRoleArgs) SetSymbol() {
	m.Symbol = m.MinValue.Symbol
}

func (m *SetTokenRoleArgs) GetAuthToken() *AuthToken {
	return &AuthToken{
		Authorizer:    m.Authorizer,
		Symbol:        m.MaxValue.Symbol.String(),
		TokenContract: m.TokenContract,
		TokenRoles: TokenRoles{
			m.GetTokenRole(),
		},
	}
}

func (m *SetTokenRoleArgs) GetTokenRole() *TokenRole {
	return NewTokenRole(string(m.TokenRole), m.MinValue.Amount, m.MaxValue.Amount, m.MinValue.Symbol)
}

type EraseTokenArgs struct {
	Authorizer eos.AccountName `json:"authorizer"`
	Symbol     eos.Symbol      `json:"symbol"`
}

type EraseTokenRoleArgs struct {
	Authorizer eos.AccountName `json:"authorizer"`
	Symbol     eos.Symbol      `json:"symbol"`
	TokenRole  eos.Name        `json:"token_role"`
}

func (m *BennyfiContract) SetToken(authToken *AuthToken) (string, error) {
	return m.ExecAction(authToken.Authorizer, "settoken", authToken.ToSetTokenArgs())
}

func (m *BennyfiContract) SetTokenRole(args *SetTokenRoleArgs) (string, error) {
	args.SetSymbol()
	return m.ExecAction(args.Authorizer, "settokenrole", args)
}

func (m *BennyfiContract) EraseToken(authorizer eos.AccountName, symbol eos.Symbol) (string, error) {
	actionData := &EraseTokenArgs{
		Authorizer: authorizer,
		Symbol:     symbol,
	}
	return m.ExecAction(authorizer, "erasetoken", actionData)
}

func (m *BennyfiContract) EraseTokenRole(authorizer eos.AccountName, symbol eos.Symbol, tokenRole eos.Name) (string, error) {
	actionData := &EraseTokenRoleArgs{
		Authorizer: authorizer,
		Symbol:     symbol,
		TokenRole:  tokenRole,
	}

	return m.ExecAction(authorizer, "erasetknrole", actionData)
}

func (m *BennyfiContract) GetTokens() ([]AuthToken, error) {
	return m.GetTokensReq(nil)
}

func (m *BennyfiContract) GetToken(symbol eos.Symbol) (*AuthToken, error) {

	request := &eos.GetTableRowsRequest{}
	err := m.FilterTokensBySymbol(request, symbol, true)
	if err != nil {
		return nil, err
	}
	authTokens, err := m.GetTokensReq(request)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	if len(authTokens) > 0 {
		return &authTokens[0], nil
	}
	return nil, nil
}

func (m *BennyfiContract) FilterTokensBySymbol(req *eos.GetTableRowsRequest, symbol eos.Symbol, exclusive bool) error {

	symbolCode, err := symbol.SymbolCode()
	if err != nil {
		return fmt.Errorf("converting symbol to uint64 %v", err)
	}
	req.LowerBound = symbolCode.String()
	if exclusive {
		req.UpperBound = req.LowerBound
		req.Limit = 1
	}
	return nil
}

func (m *BennyfiContract) GetTokensReq(req *eos.GetTableRowsRequest) ([]AuthToken, error) {

	var authTokens []AuthToken
	if req == nil {
		req = &eos.GetTableRowsRequest{}
	}
	req.Table = "authtokens"
	err := m.GetTableRows(*req, &authTokens)
	if err != nil {
		return []AuthToken{}, fmt.Errorf("get table rows %v", err)
	}
	return authTokens, nil
}
