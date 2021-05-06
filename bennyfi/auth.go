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

type Auth struct {
	Authorizer eos.AccountName `json:"authorizer"`
	Account    eos.AccountName `json:"account"`
	Level      uint64          `json:"auth_level"`
	Notes      string          `json:"notes"`
}

func (m *BennyfiContract) SetAuth(authorizer, account eos.AccountName, level uint64, notes string) (string, error) {
	actionData := make(map[string]interface{})
	actionData["authorizer"] = authorizer
	actionData["account"] = account
	actionData["auth_level"] = level
	actionData["notes"] = notes
	_, err := m.Contract.ExecAction(string(authorizer), "setauth", actionData)
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

func (m *BennyfiContract) GetAuths() ([]Auth, error) {
	return m.GetAuthsReq(nil)
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
