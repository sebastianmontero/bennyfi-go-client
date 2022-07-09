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

type Term struct {
	TermID                   uint64                  `json:"term_id"`
	TermName                 string                  `json:"term_name"`
	RoundManager             eos.AccountName         `json:"round_manager"`
	Beneficiary              eos.AccountName         `json:"beneficiary"`
	RoundType                eos.Name                `json:"round_type"`
	RoundAccess              eos.Name                `json:"round_access"`
	BeneficiaryEntryFeePerc  uint32                  `json:"beneficiary_entry_fee_perc_x100000"`
	RoundManagerEntryFeePerc uint32                  `json:"round_manager_entry_fee_perc_x100000"`
	DistributionDefinitions  DistributionDefinitions `json:"distribution_definitions"`
	CreatedDate              string                  `json:"created_date"`
	UpdatedDate              string                  `json:"updated_date"`
}

func (m *Term) UpsertDistributionDef(name string, definition interface{}) {
	if m.DistributionDefinitions == nil {
		m.DistributionDefinitions = make(DistributionDefinitions, 0, 1)
	}
	m.DistributionDefinitions.Upsert(name, definition)
}

func (m *Term) RemoveDistributionDef(name string) {
	m.DistributionDefinitions.Remove(name)
}

func (m *Term) ClearDistributionDefs() {
	m.DistributionDefinitions = nil
}

func (m *Term) GetInitializedWinners() Winners {
	winners := make(Winners, 0)
	for _, distDef := range m.DistributionDefinitions {
		var distWinners *DistributionWinners
		if IsFTDistribution(distDef.Key) {
			distWinners = NewDistributionWinners(DistributionWinnersFT{})
		} else {
			distWinners = NewDistributionWinners(DistributionWinnersNFT{})
		}
		winners = append(winners, &DistributionWinnersEntry{
			Key:   distDef.Key,
			Value: distWinners,
		})
	}
	return winners
}

type NewTermArgs struct {
	RoundManager             eos.AccountName         `json:"round_manager"`
	TermName                 string                  `json:"term_name"`
	RoundType                eos.Name                `json:"round_type"`
	RoundAccess              eos.Name                `json:"round_access"`
	Beneficiary              eos.AccountName         `json:"beneficiary"`
	BeneficiaryEntryFeePerc  uint32                  `json:"beneficiary_entry_fee_perc_x100000"`
	RoundManagerEntryFeePerc uint32                  `json:"round_manager_entry_fee_perc_x100000"`
	DistributionDefinitions  DistributionDefinitions `json:"distribution_definitions"`
}

func TermToNewTermArgs(terms *Term) *NewTermArgs {
	return &NewTermArgs{
		TermName:                 terms.TermName,
		RoundManager:             terms.RoundManager,
		Beneficiary:              terms.Beneficiary,
		RoundType:                terms.RoundType,
		RoundAccess:              terms.RoundAccess,
		BeneficiaryEntryFeePerc:  terms.BeneficiaryEntryFeePerc,
		RoundManagerEntryFeePerc: terms.RoundManagerEntryFeePerc,
		DistributionDefinitions:  terms.DistributionDefinitions,
	}
}

func (m *BennyfiContract) NewTerm(term *Term) (string, error) {
	return m.NewTermFromTermArgs(TermToNewTermArgs(term))
}

func (m *BennyfiContract) NewTermFromTermArgs(termArgs *NewTermArgs) (string, error) {
	actionData := make(map[string]interface{})
	actionData["round_manager"] = termArgs.RoundManager
	actionData["term_name"] = termArgs.TermName
	actionData["beneficiary"] = termArgs.Beneficiary
	actionData["round_type"] = termArgs.RoundType
	actionData["round_access"] = termArgs.RoundAccess
	actionData["beneficiary_entry_fee_perc_x100000"] = termArgs.BeneficiaryEntryFeePerc
	actionData["round_manager_entry_fee_perc_x100000"] = termArgs.RoundManagerEntryFeePerc
	actionData["distribution_definitions"] = termArgs.DistributionDefinitions

	return m.ExecAction(termArgs.RoundManager, "newterm", actionData)
}

func (m *BennyfiContract) GetTerms() ([]Term, error) {
	return m.GetTermsReq(nil)
}

func (m *BennyfiContract) GetTermsbyManager(termManager eos.AccountName) ([]Term, error) {
	request := &eos.GetTableRowsRequest{
		Index:      "2",
		KeyType:    "name",
		LowerBound: string(termManager),
	}
	return m.GetTermsReq(request)
}

func (m *BennyfiContract) GetLastTerm() (*Term, error) {
	terms, err := m.GetTermsReq(&eos.GetTableRowsRequest{
		Reverse: true,
		Limit:   1,
	})
	if err != nil {
		return nil, err
	}
	if len(terms) > 0 {
		return &terms[0], nil
	}
	return nil, nil
}

func (m *BennyfiContract) GetTermsReq(req *eos.GetTableRowsRequest) ([]Term, error) {

	var terms []Term
	if req == nil {
		req = &eos.GetTableRowsRequest{}
	}
	req.Table = "terms"
	err := m.GetTableRows(*req, &terms)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	return terms, nil
}
