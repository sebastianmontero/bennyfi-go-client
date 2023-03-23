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
	"github.com/sebastianmontero/eos-go-toolbox/dto"
)

type AdditionalParam struct {
	Key   string         `json:"first"`
	Value *dto.FlexValue `json:"second"`
}

type AdditionalParams []*AdditionalParam

func (m AdditionalParams) FindPos(key string) int {
	for i, attr := range m {
		if attr.Key == key {
			return i
		}
	}
	return -1
}

func (m AdditionalParams) Find(key string) *AdditionalParam {
	pos := m.FindPos(key)
	if pos >= 0 {
		return m[pos]
	}
	return nil
}

type YieldSource struct {
	YieldSource                      eos.Name        `json:"yield_source"`
	YieldSourceName                  string          `json:"yield_source_name"`
	YieldSourceDescription           string          `json:"yield_source_description"`
	StakeSymbol                      string          `json:"stake_symbol"`
	AdaptorContract                  eos.AccountName `json:"adaptor_contract"`
	YieldSourceCID                   string          `json:"yield_source_cid"`
	EntryFeePercentageOfYieldx100000 uint32          `json:"entry_fee_percentage_of_yield_x100000"`
	DailyYieldx100000                uint32          `json:"daily_yield_x100000"`
	TokenValue                       string          `json:"token_value"`
	BenyValue                        string          `json:"beny_value"`
	// NOT USED AT THE MOMENT
	// AdditionalParams                 AdditionalParams `json:"additional_params"`
}

type SetYieldSourceArgs struct {
	*YieldSource
	Authorizer eos.AccountName `json:"authorizer"`
}

func (m *BennyfiContract) SetYieldSource(yieldSourceArgs *SetYieldSourceArgs) (string, error) {
	return m.ExecAction(yieldSourceArgs.Authorizer, "setyieldsrc", yieldSourceArgs)
}

func (m *BennyfiContract) EraseYieldSource(yieldSource eos.Name, authorizer eos.AccountName) (string, error) {
	actionData := make(map[string]interface{})
	actionData["yield_source"] = yieldSource
	actionData["authorizer"] = authorizer
	return m.ExecAction(authorizer, "eraseyldsrc", actionData)
}

func (m *BennyfiContract) GetYieldSourcesReq(req *eos.GetTableRowsRequest) ([]*YieldSource, error) {

	var yieldSources []*YieldSource
	if req == nil {
		req = &eos.GetTableRowsRequest{}
	}
	req.Table = "yieldsources"
	err := m.GetTableRows(*req, &yieldSources)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	return yieldSources, nil
}

func (m *BennyfiContract) GetYieldSourceById(yieldSource eos.Name) (*YieldSource, error) {
	request := &eos.GetTableRowsRequest{}
	m.FilterYieldSourcesById(request, yieldSource)
	yieldSources, err := m.GetYieldSourcesReq(request)
	if err != nil {
		return nil, err
	}
	if len(yieldSources) > 0 {
		return yieldSources[0], nil
	}
	return nil, nil
}

func (m *BennyfiContract) FilterYieldSourcesById(req *eos.GetTableRowsRequest, yieldSource eos.Name) {
	req.LowerBound = string(yieldSource)
	req.UpperBound = string(yieldSource)
}
