package bennyfi

import (
	"fmt"

	"github.com/sebastianmontero/bennyfi-go-client/util/utype"
	eos "github.com/sebastianmontero/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/err"
)

type IDistributionWinners interface {
	FindPos(account interface{}) int
	Upsert(winner interface{}) interface{}
	Get(pos int) interface{}
	AssignPrizes(dist *Distribution) error
	Len() int
}

type DistributionWinnersFT []*WinnerFT

func (m DistributionWinnersFT) FindPos(account interface{}) int {
	for i, winner := range m {
		if winner.IsWinner(account) {
			return i
		}
	}
	return -1
}

func (m DistributionWinnersFT) Get(pos int) interface{} {
	return m[pos]
}

func (m DistributionWinnersFT) Upsert(winner interface{}) interface{} {
	winnerFT := winner.(*WinnerFT)
	pos := m.FindPos(winnerFT.Participant)
	if pos >= 0 {
		m[pos] = winnerFT
	} else {
		m = append(m, winnerFT)
	}
	return m
}

func (m DistributionWinnersFT) AssignPrizes(dist *Distribution) error {
	prizes := dist.DistributionFT().WinnerPrizes
	if len(m) != len(prizes) {
		return fmt.Errorf("failed FT assigning prizes, the numnber of winners: %v is different from the number of prizes: %v", len(m), len(prizes))
	}
	for i, prize := range prizes {
		m[i].Prize = prize
	}
	return nil
}

func (m DistributionWinnersFT) Len() int {
	return len(m)
}

type DistributionWinnersNFT []*WinnerNFT

func (m DistributionWinnersNFT) FindPos(account interface{}) int {
	for i, winner := range m {
		if winner.IsWinner(account) {
			return i
		}
	}
	return -1
}

func (m DistributionWinnersNFT) Get(pos int) interface{} {
	return m[pos]
}

func (m DistributionWinnersNFT) Upsert(winner interface{}) interface{} {
	winnerNFT := winner.(*WinnerNFT)
	pos := m.FindPos(winnerNFT.Participant)
	if pos >= 0 {
		m[pos] = winnerNFT
	} else {
		m = append(m, winnerNFT)
	}
	return m
}

func (m DistributionWinnersNFT) AssignPrizes(dist *Distribution) error {
	prizes := dist.DistributionNFT().WinnerPrizes
	if len(m) != len(prizes) {
		return fmt.Errorf("failed NFT assigning prizes, the numnber of winners: %v is different from the number of prizes: %v", len(m), len(prizes))
	}
	for i, prize := range prizes {
		m[i].Prize = prize
	}
	return nil
}

func (m DistributionWinnersNFT) Len() int {
	return len(m)
}

var DistributionWinnersVariant = eos.NewVariantDefinition([]eos.VariantType{
	{Name: "DistributionWinnersFT", Type: DistributionWinnersFT{}},
	{Name: "DistributionWinnersNFT", Type: DistributionWinnersNFT{}},
})

func GetDistributionWinnersVariants() *eos.VariantDefinition {
	return DistributionWinnersVariant
}

type DistributionWinners struct {
	eos.BaseVariant
}

func NewDistributionWinners(value interface{}) *DistributionWinners {
	return &DistributionWinners{
		BaseVariant: eos.BaseVariant{
			TypeID: GetDistributionWinnersVariants().TypeID(utype.TypeName(value)),
			Impl:   value,
		}}
}

func NewDistributionWinnersFromWinner(value interface{}) *DistributionWinners {

	switch v := value.(type) {
	case *WinnerFT:
		return NewDistributionWinners(DistributionWinnersFT{v})
	case *WinnerNFT:
		return NewDistributionWinners(DistributionWinnersNFT{v})
	default:
		panic(fmt.Sprintf("failed creating a DistributionWinners object from type %T", value))
	}

}

func (m *DistributionWinners) FindPos(account interface{}) int {
	return m.Impl.(IDistributionWinners).FindPos(account)
}

func (m *DistributionWinners) Upsert(winner interface{}) {
	m.Impl = m.Impl.(IDistributionWinners).Upsert(winner)
}

func (m *DistributionWinners) Len() int {
	return m.Impl.(IDistributionWinners).Len()
}

func (m *DistributionWinners) Find(account interface{}) interface{} {
	pos := m.FindPos(account)
	if pos >= 0 {
		return m.Impl.(IDistributionWinners).Get(pos)
	}
	return nil
}

func (m *DistributionWinners) FindFT(account interface{}) *WinnerFT {
	winner := m.Find(account)
	if winner != nil {
		return winner.(*WinnerFT)
	}
	return nil
}

func (m *DistributionWinners) FindNFT(account interface{}) *WinnerNFT {
	winner := m.Find(account)
	if winner != nil {
		return winner.(*WinnerNFT)
	}
	return nil

}

func (m *DistributionWinners) AssignPrizes(dist *Distribution) error {
	return m.Impl.(IDistributionWinners).AssignPrizes(dist)
}

func (m *DistributionWinners) DistributionWinnersNFT() DistributionWinnersNFT {
	switch v := m.Impl.(type) {
	case DistributionWinnersNFT:
		return v
	default:
		panic(&err.InvalidTypeError{
			Label:        fmt.Sprintf("received an unexpected type %T for value: %v of variant %T", v, v, m),
			ExpectedType: "DistributionWinnersNFT",
			Value:        m,
		})
	}
}

func (m *DistributionWinners) DistributionWinnersFT() DistributionWinnersFT {
	switch v := m.Impl.(type) {
	case DistributionWinnersFT:
		return v
	default:
		panic(&err.InvalidTypeError{
			Label:        fmt.Sprintf("received an unexpected type %T for value: %v of variant %T", v, v, m),
			ExpectedType: "DistributionWinnersFT",
			Value:        m,
		})
	}
}

// MarshalJSON translates to []byte
func (m *DistributionWinners) MarshalJSON() ([]byte, error) {
	return m.BaseVariant.MarshalJSON(DistributionWinnersVariant)
}

// UnmarshalJSON translates WinnerVariant
func (m *DistributionWinners) UnmarshalJSON(data []byte) error {
	return m.BaseVariant.UnmarshalJSON(data, DistributionWinnersVariant)
}

// UnmarshalBinary ...
func (m *DistributionWinners) UnmarshalBinary(decoder *eos.Decoder) error {
	return m.BaseVariant.UnmarshalBinaryVariant(decoder, DistributionWinnersVariant)
}

type DistributionWinnersEntry struct {
	Key   eos.Name             `json:"first"`
	Value *DistributionWinners `json:"second"`
}

type Winners []*DistributionWinnersEntry

func (m Winners) ToMap() map[eos.Name]interface{} {
	winnerMap := make(map[eos.Name]interface{})
	for _, winnerEntry := range m {
		winnerMap[winnerEntry.Key] = winnerEntry.Value.Impl
	}
	return winnerMap
}

func (m Winners) FindPos(key eos.Name) int {
	for i, def := range m {
		if def.Key == key {
			return i
		}
	}
	return -1
}

func (m Winners) Find(key eos.Name) *DistributionWinnersEntry {
	pos := m.FindPos(key)
	if pos >= 0 {
		return m[pos]
	}
	return nil
}

func (m Winners) HasWinners(key eos.Name) bool {
	return m.NumWinners(key) > 0
}

func (m Winners) NumWinners(key eos.Name) int {
	we := m.Find(key)
	if we != nil {
		return we.Value.Len()
	}
	return 0
}

func (m Winners) FindFT(key eos.Name) DistributionWinnersFT {
	v := m.Find(key)
	if v != nil {
		return v.Value.DistributionWinnersFT()
	}
	return nil
}

func (m Winners) FindNFT(key eos.Name) DistributionWinnersNFT {
	v := m.Find(key)
	if v != nil {
		return v.Value.DistributionWinnersNFT()
	}
	return nil
}

func (m Winners) FindWinnerFT(key eos.Name, account interface{}) *WinnerFT {
	v := m.Find(key)
	if v != nil {
		return v.Value.FindFT(account)
	}
	return nil
}

func (m Winners) FindWinnerNFT(key eos.Name, account interface{}) *WinnerNFT {
	v := m.Find(key)
	if v != nil {
		return v.Value.FindNFT(account)
	}
	return nil
}

func (p *Winners) Upsert(key eos.Name, winner interface{}) {
	m := *p
	pos := m.FindPos(key)
	if pos >= 0 {
		winners := m[pos].Value
		winners.Upsert(winner)
	} else {
		m = append(m, &DistributionWinnersEntry{
			Key:   key,
			Value: NewDistributionWinnersFromWinner(winner),
		})
	}
	*p = m
}

func (p *Winners) Remove(key eos.Name) *DistributionWinnersEntry {
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
