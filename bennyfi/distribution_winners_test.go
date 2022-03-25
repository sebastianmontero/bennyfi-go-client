package bennyfi_test

import (
	"testing"

	"github.com/eoscanada/eos-go"
	"github.com/sebastianmontero/bennyfi-go-client/bennyfi"
	"github.com/sebastianmontero/bennyfi-go-client/bennyfi/test"
)

func TestWinnersUpsert(t *testing.T) {
	actual := bennyfi.Winners{}
	key1 := "key1"
	winnerFT := bennyfi.NewWinnerFT(eos.AN("account1"), "10.0000 TLOS", 1)

	distWinnersFT := bennyfi.DistributionWinnersFT{winnerFT}
	expected := bennyfi.Winners{
		{
			Key:   key1,
			Value: bennyfi.NewDistributionWinners(distWinnersFT),
		},
	}
	actual.Upsert(key1, winnerFT)
	test.AssertWinners(t, actual, expected)

	winnerFT = bennyfi.NewWinnerFT(eos.AN("account2"), "9.0000 TLOS", 2)
	distWinnersFT = append(distWinnersFT, winnerFT)
	expected[0].Value.Impl = bennyfi.NewDistributionWinners(distWinnersFT)
	actual.Upsert(key1, winnerFT)
	test.AssertWinners(t, actual, expected)

	key2 := "key2"
	winnerNFT := bennyfi.NewWinnerNFT(eos.AN("account1"), 2, 1)
	distWinnersNFT := bennyfi.DistributionWinnersNFT{winnerNFT}
	expected = append(expected, &bennyfi.DistributionWinnersEntry{
		Key:   key2,
		Value: bennyfi.NewDistributionWinners(distWinnersNFT),
	})

	actual.Upsert(key2, winnerNFT)
	test.AssertWinners(t, actual, expected)

}
