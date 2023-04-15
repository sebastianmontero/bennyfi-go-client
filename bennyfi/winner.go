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

	eos "github.com/sebastianmontero/eos-go"
)

type BaseWinner struct {
	Participant   eos.AccountName `json:"participant"`
	EntryPosition uint64          `json:"entry_position"`
}

func (m *BaseWinner) IsWinner(account interface{}) bool {
	return fmt.Sprintf("%v", m.Participant) == fmt.Sprintf("%v", account)
}

func NewBaseWinner(participant eos.AccountName, entryPosition uint64) *BaseWinner {
	return &BaseWinner{
		Participant:   participant,
		EntryPosition: entryPosition,
	}
}

type WinnerFT struct {
	*BaseWinner
	Prize eos.Asset `json:"prize"`
}

func NewWinnerFT(participant eos.AccountName, prize eos.Asset, entryPosition uint64) *WinnerFT {
	return &WinnerFT{
		BaseWinner: NewBaseWinner(participant, entryPosition),
		Prize:      prize,
	}
}

type WinnerNFT struct {
	*BaseWinner
	Prize uint16 `json:"prize"`
}

func NewWinnerNFT(participant eos.AccountName, prize uint16, entryPosition uint64) *WinnerNFT {
	return &WinnerNFT{
		BaseWinner: NewBaseWinner(participant, entryPosition),
		Prize:      prize,
	}
}
