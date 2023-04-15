package bennyfi_test

import (
	"testing"

	"github.com/sebastianmontero/bennyfi-go-client/bennyfi"
	"github.com/sebastianmontero/bennyfi-go-client/bennyfi/test"
	"github.com/sebastianmontero/eos-go"
)

func TestDistributionAssert(t *testing.T) {
	symbol := eos.Symbol{Precision: 2, Symbol: "TLOS"}
	expected := bennyfi.NewDistribution(&bennyfi.DistributionFT{
		BeneficiaryReward:    eos.Asset{Amount: 19000, Symbol: symbol},
		MinParticipantReward: eos.Asset{Amount: 29000, Symbol: symbol},
		RoundManagerFee:      eos.Asset{Amount: 39000, Symbol: symbol},
		WinnerPrizes: []eos.Asset{
			{Amount: 38000, Symbol: symbol},
			{Amount: 37000, Symbol: symbol},
			{Amount: 35000, Symbol: symbol},
		},
	})
	actual := bennyfi.NewDistribution(&bennyfi.DistributionFT{
		BeneficiaryReward:    eos.Asset{Amount: 19000, Symbol: symbol},
		MinParticipantReward: eos.Asset{Amount: 29000, Symbol: symbol},
		RoundManagerFee:      eos.Asset{Amount: 39000, Symbol: symbol},
		WinnerPrizes: []eos.Asset{
			{Amount: 38000, Symbol: symbol},
			{Amount: 37000, Symbol: symbol},
			{Amount: 35000, Symbol: symbol},
		},
	})

	test.AssertDistribution(t, actual, expected)

	expected = bennyfi.NewDistribution(&bennyfi.DistributionNFT{
		BeneficiaryReward:    1,
		MinParticipantReward: 2,
		RoundManagerFee:      3,
		WinnerPrizes: []uint16{
			4,
			5,
			6,
		},
	})
	actual = bennyfi.NewDistribution(&bennyfi.DistributionNFT{
		BeneficiaryReward:    1,
		MinParticipantReward: 2,
		RoundManagerFee:      3,
		WinnerPrizes: []uint16{
			4,
			5,
			6,
		},
	})

	test.AssertDistribution(t, actual, expected)

}
