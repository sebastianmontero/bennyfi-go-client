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

var (
	Sudo         uint64 = 0
	Admin        uint64 = 10
	Enroller     uint64 = 20
	TermManager  uint64 = 30
	RoundManager uint64 = 40
	Player       uint64 = 50
	Beneficiary  uint64 = 60
	Disabled     uint64 = 70
)

type Auth struct {
	Authorizer   eos.AccountName `json:"authorizer"`
	Account      eos.AccountName `json:"account"`
	Level        uint64          `json:"auth_level"`
	DisplayName  string          `json:"display_name"`
	ArtifactCID  string          `json:"artifact_cid"`
	Notes        string          `json:"notes"`
	StakedAmount string          `json:"staked_amount"`
}

func (m *BennyfiContract) SetAuth(auth *Auth) (string, error) {
	_, err := m.Contract.ExecAction(auth.Authorizer, "setauth", auth)
	if err != nil {
		return "", err
	}
	return "", nil
}

func (m *BennyfiContract) SetAuthLevel(authorizer, account eos.AccountName, level uint64, notes string) (string, error) {
	actionData := make(map[string]interface{})
	actionData["authorizer"] = authorizer
	actionData["account"] = account
	actionData["auth_level"] = level
	actionData["notes"] = notes
	_, err := m.Contract.ExecAction(authorizer, "setauthlevel", actionData)
	if err != nil {
		return "", err
	}
	return "", nil
}

func (m *BennyfiContract) SetProfile(account eos.AccountName, displayName, artifactCID string) (string, error) {
	actionData := make(map[string]interface{})
	actionData["account"] = account
	actionData["display_name"] = displayName
	actionData["artifact_cid"] = artifactCID
	_, err := m.Contract.ExecAction(account, "setprofile", actionData)
	if err != nil {
		return "", err
	}
	return "", nil
}

func (m *BennyfiContract) EraseAuth(authorizer, account eos.AccountName) (string, error) {
	actionData := make(map[string]interface{})
	actionData["authorizer"] = authorizer
	actionData["account"] = account
	_, err := m.Contract.ExecAction(string(authorizer), "eraseauth", actionData)
	if err != nil {
		return "", err
	}
	return "", nil
}

func (m *BennyfiContract) StakeAuth(account eos.AccountName, level uint64) (string, error) {
	actionData := make(map[string]interface{})
	actionData["account"] = account
	actionData["auth_level"] = level
	_, err := m.Contract.ExecAction(string(account), "stakeauth", actionData)
	if err != nil {
		return "", err
	}
	return "", nil
}

func (m *BennyfiContract) UnstakeAuth(authorizer, account eos.AccountName) (string, error) {
	actionData := make(map[string]interface{})
	actionData["authorizer"] = authorizer
	actionData["account"] = account
	_, err := m.Contract.ExecAction(string(authorizer), "unstakeauth", actionData)
	if err != nil {
		return "", err
	}
	return "", nil
}

func (m *BennyfiContract) CanCreateToken(user eos.AccountName) (string, error) {
	actionData := make(map[string]interface{})
	actionData["user"] = user
	_, err := m.Contract.ExecAction(string(user), "cancreatetkn", actionData)
	if err != nil {
		return "", err
	}
	return "", nil
}

func (m *BennyfiContract) GetAuths() ([]Auth, error) {
	return m.GetAuthsReq(nil)
}

func (m *BennyfiContract) GetAuth(accountName interface{}) (*Auth, error) {
	request := &eos.GetTableRowsRequest{}
	m.FilterAuthsByAccount(request, accountName, true)
	auths, err := m.GetAuthsReq(request)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	if len(auths) > 0 {
		return &auths[0], nil
	}
	return nil, nil
}

func (m *BennyfiContract) FilterAuthsByAccount(req *eos.GetTableRowsRequest, account interface{}, exclusive bool) {
	req.LowerBound = fmt.Sprintf("%v", account)
	if exclusive {
		req.UpperBound = req.LowerBound
		req.Limit = 1
	}
}

func (m *BennyfiContract) GetAuthsReq(req *eos.GetTableRowsRequest) ([]Auth, error) {
	var auths []Auth
	if req == nil {
		req = &eos.GetTableRowsRequest{}
	}
	req.Table = "auths"
	err := m.GetTableRows(*req, &auths)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	return auths, nil
}
