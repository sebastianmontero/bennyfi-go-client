package bennyfi_test

import (
	"testing"

	"github.com/eoscanada/eos-go"
	"github.com/sebastianmontero/bennyfi-go-client/bennyfi"
	"github.com/sebastianmontero/bennyfi-go-client/bennyfi/test"
)

func TestWinnersUpsert(t *testing.T) {
	symbol := eos.Symbol{Precision: 4, Symbol: "TLOS"}
	actual := bennyfi.Winners{}
	key1 := eos.Name("key1")
	winnerFT := bennyfi.NewWinnerFT(eos.AN("account1"), eos.Asset{Amount: 100000, Symbol: symbol}, 1)

	distWinnersFT := bennyfi.DistributionWinnersFT{winnerFT}
	expected := bennyfi.Winners{
		{
			Key:   key1,
			Value: bennyfi.NewDistributionWinners(distWinnersFT),
		},
	}
	actual.Upsert(key1, winnerFT)
	test.AssertWinners(t, actual, expected)

	winnerFT = bennyfi.NewWinnerFT(eos.AN("account2"), eos.Asset{Amount: 90000, Symbol: symbol}, 2)
	distWinnersFT = append(distWinnersFT, winnerFT)
	expected[0].Value.Impl = bennyfi.NewDistributionWinners(distWinnersFT)
	actual.Upsert(key1, winnerFT)
	test.AssertWinners(t, actual, expected)

	key2 := eos.Name("key2")
	winnerNFT := bennyfi.NewWinnerNFT(eos.AN("account1"), 2, 1)
	distWinnersNFT := bennyfi.DistributionWinnersNFT{winnerNFT}
	expected = append(expected, &bennyfi.DistributionWinnersEntry{
		Key:   key2,
		Value: bennyfi.NewDistributionWinners(distWinnersNFT),
	})

	actual.Upsert(key2, winnerNFT)
	test.AssertWinners(t, actual, expected)

}
