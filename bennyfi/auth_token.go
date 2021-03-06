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

type TokenLimits struct {
	MinValue string `json:"min_value"`
	MaxValue string `json:"max_value"`
}

type TokenRole struct {
	Key   eos.Name     `json:"key"`
	Value *TokenLimits `json:"value"`
}

func NewTokenRole(tokenRole string, minValue, maxValue eos.Int64, symbol eos.Symbol) *TokenRole {
	return &TokenRole{
		Key: eos.Name(tokenRole),
		Value: &TokenLimits{
			MinValue: eos.Asset{Amount: minValue, Symbol: symbol}.String(),
			MaxValue: eos.Asset{Amount: maxValue, Symbol: symbol}.String(),
		},
	}
}

type TokenRoles []*TokenRole

// type TokenRolesDTO []map[]

type AuthToken struct {
	Authorizer    eos.AccountName `json:"authorizer"`
	Symbol        string          `json:"symbol"`
	TokenContract eos.AccountName `json:"token_contract"`
	TokenRoles    TokenRoles      `json:"token_roles"`
}

type SetTokenRoleArgs struct {
	Authorizer    eos.AccountName
	TokenContract eos.AccountName
	TokenRole     eos.Name
	MinValue      eos.Asset
	MaxValue      eos.Asset
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

func (m *BennyfiContract) SetToken(authToken *AuthToken) (string, error) {
	actionData := make(map[string]interface{})
	actionData["authorizer"] = authToken.Authorizer
	actionData["symbol"] = authToken.Symbol
	actionData["token_contract"] = authToken.TokenContract
	actionData["token_roles"] = authToken.TokenRoles

	return m.ExecAction(authToken.Authorizer, "settoken", actionData)
}

func (m *BennyfiContract) SetTokenRole(args *SetTokenRoleArgs) (string, error) {
	actionData := make(map[string]interface{})
	actionData["authorizer"] = args.Authorizer
	actionData["symbol"] = args.MinValue.Symbol.String()
	actionData["token_contract"] = args.TokenContract
	actionData["token_role"] = args.TokenRole
	actionData["min_value"] = args.MinValue.String()
	actionData["max_value"] = args.MaxValue.String()

	return m.ExecAction(args.Authorizer, "settokenrole", actionData)
}

func (m *BennyfiContract) EraseToken(authorizer eos.AccountName, symbol eos.Symbol) (string, error) {
	actionData := make(map[string]interface{})
	actionData["authorizer"] = authorizer
	actionData["symbol"] = symbol.String()

	return m.ExecAction(authorizer, "erasetoken", actionData)
}

func (m *BennyfiContract) EraseTokenRole(authorizer eos.AccountName, symbol eos.Symbol, tokenRole eos.Name) (string, error) {
	actionData := make(map[string]interface{})
	actionData["authorizer"] = authorizer
	actionData["symbol"] = symbol.String()
	actionData["token_role"] = tokenRole

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
