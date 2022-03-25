package bennyfi_test

import (
	"testing"

	"github.com/sebastianmontero/bennyfi-go-client/bennyfi"
	"github.com/sebastianmontero/bennyfi-go-client/bennyfi/test"
)

func TestDistributionAssert(t *testing.T) {
	expected := bennyfi.NewDistribution(&bennyfi.DistributionFT{
		BeneficiaryReward:    "190.00 TLOS",
		MinParticipantReward: "290.00 TLOS",
		RoundManagerFee:      "390.00 TLOS",
		WinnerPrizes: []string{
			"380.00 TLOS",
			"370.00 TLOS",
			"350.00 TLOS",
		},
	})
	actual := bennyfi.NewDistribution(&bennyfi.DistributionFT{
		BeneficiaryReward:    "190.00 TLOS",
		MinParticipantReward: "290.00 TLOS",
		RoundManagerFee:      "390.00 TLOS",
		WinnerPrizes: []string{
			"380.00 TLOS",
			"370.00 TLOS",
			"350.00 TLOS",
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
