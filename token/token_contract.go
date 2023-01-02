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
package token

import (
	"github.com/eoscanada/eos-go/token"
	"github.com/sebastianmontero/eos-go-toolbox/contract"
	"github.com/sebastianmontero/eos-go-toolbox/service"
	"github.com/sebastianmontero/eos-go-toolbox/util"
)

type TokenContract struct {
	*contract.TokenContract
}

func NewTokenContract(eos *service.EOS) *TokenContract {
	return &TokenContract{
		contract.NewTokenContract(eos),
	}
}

func (m *TokenContract) CreateTokenRaw(contract, authorizerName, issuerName, maxSupply interface{}, failIfExists bool) (*service.PushTransactionFullResp, error) {

	authorizer, err := util.ToAccountName(authorizerName)
	if err != nil {
		return nil, err
	}

	issuer, err := util.ToAccountName(issuerName)
	if err != nil {
		return nil, err
	}

	supply, err := util.ToAsset(maxSupply)
	if err != nil {
		return nil, err
	}

	data := token.Create{
		Issuer:        issuer,
		MaximumSupply: supply,
	}
	resp, err := m.ExecActionC(contract, authorizer, "create", data)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
