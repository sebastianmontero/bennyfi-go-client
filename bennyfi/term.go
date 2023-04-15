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

type DefaultValue struct {
	Key   string         `json:"first"`
	Value *dto.FlexValue `json:"second"`
}

func (m *DefaultValue) String() string {
	return fmt.Sprintf("Key: %v, Value: %v", m.Key, m.Value)
}

type DefaultValues []*DefaultValue

func (m DefaultValues) ToMap() map[string]interface{} {
	defaultValueMap := make(map[string]interface{})
	for _, defaultValueEntry := range m {
		defaultValueMap[defaultValueEntry.Key] = defaultValueEntry.Value.Impl
	}
	return defaultValueMap
}

func (m DefaultValues) FindPos(key string) int {
	for i, attr := range m {
		if attr.Key == key {
			return i
		}
	}
	return -1
}

func (m DefaultValues) Find(key string) *DefaultValue {
	pos := m.FindPos(key)
	if pos >= 0 {
		return m[pos]
	}
	return nil
}

type Term struct {
	TermID                   uint64                  `json:"term_id"`
	TermName                 string                  `json:"term_name"`
	Authorizer               eos.AccountName         `json:"authorizer"`
	RoundType                eos.Name                `json:"round_type"`
	RoundAccess              eos.Name                `json:"round_access"`
	NumParticipants          uint32                  `json:"num_participants"`
	EntryStake               eos.Asset               `json:"entry_stake"`
	StakingPeriod            *dto.Microseconds       `json:"staking_period"`
	EnrollmentTimeOut        *dto.Microseconds       `json:"enrollment_time_out"`
	BeneficiaryEntryFeePerc  uint32                  `json:"beneficiary_entry_fee_perc_x100000"`
	RoundManagerEntryFeePerc uint32                  `json:"round_manager_entry_fee_perc_x100000"`
	DistributionDefinitions  DistributionDefinitions `json:"distribution_definitions"`
	DefaultValues            DefaultValues           `json:"default_values"`
	CreatedDate              eos.TimePoint           `json:"created_date"`
	UpdatedDate              eos.TimePoint           `json:"updated_date"`
	*Deletable
}

type TermCustomJSON struct {
	DistributionDefinitions map[eos.Name]interface{} `json:"distribution_definitions"`
	DefaultValues           map[string]interface{}   `json:"default_values"`
	Term
}

func (m Term) ToCustomJSON() TermCustomJSON {
	return TermCustomJSON{
		DistributionDefinitions: m.DistributionDefinitions.ToMap(),
		DefaultValues:           m.DefaultValues.ToMap(),
		Term:                    m,
	}
}

// func (m *Term) GetEntryStake() eos.Asset {
// 	entryStake, err := eos.NewAssetFromString(m.EntryStake)
// 	if err != nil {
// 		panic(fmt.Sprintf("Unable to parse entry stake: %v to asset", m.EntryStake))
// 	}
// 	return entryStake
// }

func (m *Term) UpsertDistributionDef(name eos.Name, definition interface{}) {
	if m.DistributionDefinitions == nil {
		m.DistributionDefinitions = make(DistributionDefinitions, 0, 1)
	}
	m.DistributionDefinitions.Upsert(name, definition)
}

func (m *Term) RemoveDistributionDef(name eos.Name) {
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

func (m *Term) RequiresBeneficiary() bool {
	for _, distDef := range m.DistributionDefinitions {
		if distDef.Value.HasBeneficiaryReward() {
			return true
		}
	}
	return m.DistributionDefinitions.Has(DistributionProjectToken) || m.DistributionDefinitions.Has(DistributionProjectNFT)
}

type NewTermArgs struct {
	Authorizer               eos.AccountName         `json:"authorizer"`
	TermName                 string                  `json:"term_name"`
	RoundType                eos.Name                `json:"round_type"`
	RoundAccess              eos.Name                `json:"round_access"`
	NumParticipants          uint32                  `json:"num_participants"`
	EntryStake               eos.Asset               `json:"entry_stake"`
	StakingPeriodHrs         uint32                  `json:"staking_period_hrs"`
	EnrollmentTimeOutHrs     uint32                  `json:"enrollment_time_out_hrs"`
	BeneficiaryEntryFeePerc  uint32                  `json:"beneficiary_entry_fee_perc_x100000"`
	RoundManagerEntryFeePerc uint32                  `json:"round_manager_entry_fee_perc_x100000"`
	DistributionDefinitions  DistributionDefinitions `json:"distribution_definitions"`
	DefaultValues            DefaultValues           `json:"default_values"`
}

func TermToNewTermArgs(terms *Term) *NewTermArgs {
	return &NewTermArgs{
		TermName:                 terms.TermName,
		Authorizer:               terms.Authorizer,
		RoundType:                terms.RoundType,
		RoundAccess:              terms.RoundAccess,
		NumParticipants:          terms.NumParticipants,
		EntryStake:               terms.EntryStake,
		StakingPeriodHrs:         uint32(terms.StakingPeriod.Hrs()),
		EnrollmentTimeOutHrs:     uint32(terms.EnrollmentTimeOut.Hrs()),
		BeneficiaryEntryFeePerc:  terms.BeneficiaryEntryFeePerc,
		RoundManagerEntryFeePerc: terms.RoundManagerEntryFeePerc,
		DistributionDefinitions:  terms.DistributionDefinitions,
		DefaultValues:            terms.DefaultValues,
	}
}

type EraseTermArgs struct {
	TermID     uint64          `json:"term_id"`
	Authorizer eos.AccountName `json:"authorizer"`
	Erase      bool            `json:"erase"`
}

func (m *BennyfiContract) NewTerm(term *Term) (string, error) {
	return m.NewTermFromTermArgs(TermToNewTermArgs(term))
}

func (m *BennyfiContract) NewTermFromTermArgs(termArgs *NewTermArgs) (string, error) {
	return m.ExecAction(termArgs.Authorizer, "newterm", termArgs)
}

func (m *BennyfiContract) EraseTerm(termId uint64, authorizer eos.AccountName, erase bool) (string, error) {
	actionData := &EraseTermArgs{
		TermID:     termId,
		Authorizer: authorizer,
		Erase:      erase,
	}
	return m.ExecAction(authorizer, "eraseterm", actionData)
}

func (m *BennyfiContract) GetAllTerms() ([]Term, error) {
	var terms []Term
	req := eos.GetTableRowsRequest{
		Table: "terms",
	}
	err := m.GetAllTableRows(req, "term_id", &terms)
	if err != nil {
		return nil, fmt.Errorf("failed getting all terms: %v", err)
	}
	return terms, nil

}

func (m *BennyfiContract) GetAllTermsAsMap() ([]map[string]interface{}, error) {
	req := eos.GetTableRowsRequest{
		Table: "terms",
	}
	return m.GetAllTableRowsAsMap(req, "term_id")
}

func (m *BennyfiContract) GetTerms() ([]Term, error) {
	return m.GetTermsReq(nil)
}

func (m *BennyfiContract) GetTermsById(termId uint64) (*Term, error) {
	terms, err := m.GetTermsReq(&eos.GetTableRowsRequest{
		LowerBound: fmt.Sprintf("%v", termId),
		UpperBound: fmt.Sprintf("%v", termId),
	})
	if err != nil {
		return nil, err
	}
	if len(terms) > 0 {
		return &terms[0], nil
	}
	return nil, nil
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
