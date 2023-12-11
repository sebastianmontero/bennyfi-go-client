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

type YieldSourceStruct = YieldSource

type YieldSource struct {
	YieldSource                      eos.Name        `json:"yield_source"`
	YieldSourceName                  string          `json:"yield_source_name"`
	YieldSourceDescription           string          `json:"yield_source_description"`
	StakeSymbol                      eos.Symbol      `json:"stake_symbol"`
	AdaptorContract                  eos.AccountName `json:"adaptor_contract"`
	YieldSourceCID                   string          `json:"yield_source_cid"`
	EntryFeePercentageOfYieldx100000 uint32          `json:"entry_fee_percentage_of_yield_x100000"`
	DailyYieldx100000                uint32          `json:"daily_yield_x100000"`
	TokenValue                       eos.Asset       `json:"token_value"`
	BenyValue                        eos.Asset       `json:"beny_value"`
	Authorizer                       eos.AccountName `json:"authorizer"`
	// NOT USED AT THE MOMENT
	// AdditionalFields types.AdditionalFields `json:"additional_fields"`
}

func (m *YieldSource) Update(args *UpdateYieldSourceArgs) {
	m.DailyYieldx100000 = args.DailyYieldx100000
	m.TokenValue = args.TokenValue
	m.BenyValue = args.BenyValue
}

func (m *YieldSource) String() string {
	result, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Sprintf("Failed marshalling round: %v", err))
	}
	return string(result)
}

type UpdateYieldSourceArgs struct {
	YieldSource       eos.Name        `json:"yield_source"`
	DailyYieldx100000 uint32          `json:"daily_yield_x100000"`
	TokenValue        eos.Asset       `json:"token_value"`
	BenyValue         eos.Asset       `json:"beny_value"`
	Authorizer        eos.AccountName `json:"authorizer"`
}

func (m *UpdateYieldSourceArgs) String() string {
	result, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Sprintf("Failed marshalling round: %v", err))
	}
	return string(result)
}

type EraseYieldSourceArgs struct {
	YieldSource eos.Name        `json:"yield_source"`
	Authorizer  eos.AccountName `json:"authorizer"`
}

func (m *EraseYieldSourceArgs) String() string {
	result, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Sprintf("Failed marshalling round: %v", err))
	}
	return string(result)
}

func (m *BennyfiContract) SetYieldSource(yieldSource *YieldSource) (string, error) {

	return m.ExecAction(yieldSource.Authorizer, "setyieldsrc", yieldSource)
}

func (m *BennyfiContract) UpdateYieldSource(args *UpdateYieldSourceArgs) (string, error) {

	return m.ExecAction(args.Authorizer, "updyieldsrc", args)
}

func (m *BennyfiContract) EraseYieldSource(yieldSource eos.Name, authorizer eos.AccountName) (string, error) {
	actionData := &EraseYieldSourceArgs{
		YieldSource: yieldSource,
		Authorizer:  authorizer,
	}
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

func (m *BennyfiContract) GetAllYieldSourcesAsMap() ([]map[string]interface{}, error) {
	req := eos.GetTableRowsRequest{
		Table: "yieldsources",
	}
	return m.GetAllTableRowsAsMap(req, "yield_source")
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
