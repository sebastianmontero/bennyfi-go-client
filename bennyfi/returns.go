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

	"github.com/sebastianmontero/bennyfi-go-client/util/utype"
	eos "github.com/sebastianmontero/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/err"
	"github.com/sebastianmontero/eos-go-toolbox/util"
)

type IReturn interface {
	HasReturns() bool
}

type ReturnsFT struct {
	Prize              eos.Asset `json:"prize"`
	MinimumPayout      eos.Asset `json:"minimum_payout"`
	AmountPaidOut      eos.Asset `json:"amount_paid_out"`
	EarlyExitReturnFee eos.Asset `json:"early_exit_return_fee"`
}

func (m *ReturnsFT) HasReturns() bool {
	return m.GetTotalReturn().Amount > 0
}

func (m *ReturnsFT) GetTotalReturn() eos.Asset {
	return m.Prize.Add(m.MinimumPayout)
}

func (m *ReturnsFT) PaidTotalAmount() {
	m.AmountPaidOut = m.GetTotalReturn()
}

func (m *ReturnsFT) PaidAmount(amount interface{}) {
	amnt, err := util.ToAsset(amount)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse paid amount: %v error: %v", amount, err))
	}
	paid := m.AmountPaidOut.Add(amnt)
	if paid.Amount > m.GetTotalReturn().Amount {
		panic(fmt.Sprintf("Total Paid amount: %v is greater than round manager fee: %v, current payment: %v", paid, m.GetTotalReturn(), amount))
	}
	m.AmountPaidOut = paid
}

type ReturnsNFT struct {
	Prize         uint16 `json:"prize"`
	MinimumPayout uint16 `json:"minimum_payout"`
	AmountPaidOut uint16 `json:"amount_paid_out"`
}

func (m *ReturnsNFT) HasReturns() bool {
	return m.GetTotalReturn() > 0
}

func (m *ReturnsNFT) GetTotalReturn() uint16 {
	return m.Prize + m.MinimumPayout
}

func (m *ReturnsNFT) PaidTotalAmount() {
	m.AmountPaidOut = m.GetTotalReturn()
}

func (m *ReturnsNFT) PaidAmount(amount interface{}) {
	amnt := amount.(uint16)
	paid := m.AmountPaidOut + amnt
	if paid > m.GetTotalReturn() {
		panic(fmt.Sprintf("Total Paid amount: %v is greater than round manager fee: %v, current payment: %v", paid, m.GetTotalReturn(), amount))
	}
	m.AmountPaidOut = paid
}

var ReturnsVariant = eos.NewVariantDefinition([]eos.VariantType{
	{Name: "ReturnsFT", Type: &ReturnsFT{}},
	{Name: "ReturnsNFT", Type: &ReturnsNFT{}},
})

func GetReturnsVariants() *eos.VariantDefinition {
	return ReturnsVariant
}

type Returns struct {
	eos.BaseVariant
}

func NewReturn(value interface{}) *Returns {
	return &Returns{
		BaseVariant: eos.BaseVariant{
			TypeID: GetReturnsVariants().TypeID(utype.TypeName(value)),
			Impl:   value,
		}}
}

func (m *Returns) HasReturns() bool {
	return m.Impl.(IReturn).HasReturns()
}

func (m *Returns) ReturnsNFT() *ReturnsNFT {
	switch v := m.Impl.(type) {
	case *ReturnsNFT:
		return v
	default:
		panic(&err.InvalidTypeError{
			Label:        fmt.Sprintf("received an unexpected type %T for value: %v of variant %T", v, v, m),
			ExpectedType: "ReturnsNFT",
			Value:        m,
		})
	}
}

func (m *Returns) ReturnsFT() *ReturnsFT {
	switch v := m.Impl.(type) {
	case *ReturnsFT:
		return v
	default:
		panic(&err.InvalidTypeError{
			Label:        fmt.Sprintf("received1 an unexpected type %T for value: %v of variant %T", v, v, m),
			ExpectedType: "ReturnsFT",
			Value:        m,
		})
	}
}

// MarshalJSON translates to []byte
func (m *Returns) MarshalJSON() ([]byte, error) {
	return m.BaseVariant.MarshalJSON(ReturnsVariant)
}

// UnmarshalJSON translates WinnerVariant
func (m *Returns) UnmarshalJSON(data []byte) error {
	return m.BaseVariant.UnmarshalJSON(data, ReturnsVariant)
}

// UnmarshalBinary ...
func (m *Returns) UnmarshalBinary(decoder *eos.Decoder) error {
	return m.BaseVariant.UnmarshalBinaryVariant(decoder, ReturnsVariant)
}

type ReturnsEntry struct {
	Key   eos.Name `json:"first"`
	Value *Returns `json:"second"`
}

type ReturnEntries []*ReturnsEntry

func (m ReturnEntries) ToMap() map[eos.Name]interface{} {
	returnsMap := make(map[eos.Name]interface{})
	for _, returnsEntry := range m {
		returnsMap[returnsEntry.Key] = returnsEntry.Value.Impl
	}
	return returnsMap
}

func (m ReturnEntries) FindPos(key eos.Name) int {
	for i, def := range m {
		if def.Key == key {
			return i
		}
	}
	return -1
}

func (m ReturnEntries) Find(key eos.Name) *ReturnsEntry {
	pos := m.FindPos(key)
	if pos >= 0 {
		return m[pos]
	}
	return nil
}

func (m ReturnEntries) FindFT(key eos.Name) *ReturnsFT {
	v := m.Find(key)
	if v != nil {
		return v.Value.ReturnsFT()
	}
	return nil
}

func (m ReturnEntries) FindNFT(key eos.Name) *ReturnsNFT {
	v := m.Find(key)
	if v != nil {
		return v.Value.ReturnsNFT()
	}
	return nil
}

func (p *ReturnEntries) Upsert(key eos.Name, ret interface{}) {
	m := *p
	pos := m.FindPos(key)
	defEntry := &ReturnsEntry{
		Key:   key,
		Value: NewReturn(ret),
	}
	if pos >= 0 {
		m[pos] = defEntry
	} else {
		m = append(m, defEntry)
	}
	*p = m
}

func (p *ReturnEntries) Remove(key eos.Name) *ReturnsEntry {
	m := *p
	pos := m.FindPos(key)
	if pos >= 0 {
		def := m[pos]
		m[pos] = m[len(m)-1]
		m = m[:len(m)-1]
		*p = m
		return def
	}
	return nil
}
