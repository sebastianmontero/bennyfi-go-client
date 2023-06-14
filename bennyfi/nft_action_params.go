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

type NFTActionParams struct {
	TermId          uint64          `json:"term_id"`
	TermName        string          `json:"term_name"`
	ProjectId       uint64          `json:"project_id"`
	RoundId         uint64          `json:"round_id"`
	RoundName       string          `json:"round_name"`
	Distribution    eos.Name        `json:"distribution"`
	NumParticipants uint32          `json:"num_participants"`
	Funder          eos.AccountName `json:"funder"`
	Recipient       eos.AccountName `json:"recipient"`
	Amount          uint16          `json:"amount"`
}

func (m *NFTActionParams) ToMap() map[string]interface{} {
	jsonStr, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Sprintf("failed transforming NFTActionParams to map, error marshalling: %v", err))
	}
	var paramsMap map[string]interface{}
	err = json.Unmarshal([]byte(jsonStr), &paramsMap)
	if err != nil {
		panic(fmt.Sprintf("failed transforming NFTActionParams to map, error unmarshalling: %v", err))
	}
	return paramsMap
}
