package test

import (
	"math"
	"testing"

	"github.com/sebastianmontero/bennyfi-go-client/bennyfi"
	"github.com/sebastianmontero/bennyfi-go-client/util/utype"
	"github.com/sebastianmontero/eos-go"
	"gotest.tools/assert"
)

func AssertWinners(t *testing.T, actual, expected bennyfi.Winners) {
	assert.Equal(t, len(actual), len(expected), "Different number of winner entries, actual: %v, expected: %v", len(actual), len(expected))
	for _, ev := range expected {
		av := actual.Find(ev.Key)
		assert.Assert(t, av != nil, "Expected winners key: %v not found", ev.Key)
		AssertDistributionWinnersEntry(t, av, ev)
	}
}

func AssertDistributionWinnersEntry(t *testing.T, actual, expected *bennyfi.DistributionWinnersEntry) {
	assert.Equal(t, actual.Key, expected.Key, "Distribution entry keys don't match actual: %v, expected: %v", actual.Key, expected.Key)
	assert.DeepEqual(t, actual.Value.Impl, expected.Value.Impl)
	assert.Assert(t, utype.AreSameType(actual.Value.Impl, expected.Value.Impl), "DistributionWinners for key: %v are of different type, actual: %T, expected: %T", expected.Key, actual.Value.Impl, expected.Value.Impl)
	AssertDistributionWinners(t, actual.Value, expected.Value)
}

func AssertDistributionWinners(t *testing.T, actual, expected *bennyfi.DistributionWinners) {
	assert.Check(t, actual != nil)
	assert.Assert(t, utype.AreSameType(actual.Impl, expected.Impl), "DistributionWinners are of different type, actual: %T, expected: %T", actual.Impl, expected.Impl)
	assert.DeepEqual(t, actual.Impl, expected.Impl)
}

func AssertDistributionDefs(t *testing.T, actual, expected bennyfi.DistributionDefinitions) {
	assert.Equal(t, len(actual), len(expected), "Number of Distribution definitions do not match, actual: %v, expected: %v", len(actual), len(expected))
	for _, ee := range expected {
		ae := actual.Find(ee.Key)
		assert.Assert(t, ae != nil, "Expected distribution definition: %v, not found", ee.Key)
		AssertDistributionDef(t, ae.Value, ee.Value)
	}
}

func AssertDistributionDef(t *testing.T, actual, expected *bennyfi.DistributionDefinition) {
	assert.Check(t, actual != nil)
	assert.Assert(t, utype.AreSameType(actual.Impl, expected.Impl), "DistributionDefinition are of different type, actual: %T, expected: %T", actual.Impl, expected.Impl)
	switch v := expected.Impl.(type) {
	case *bennyfi.DistributionDefinitionFT:
		AssertDistributionDefFT(t, actual.Impl.(*bennyfi.DistributionDefinitionFT), v)
	case *bennyfi.DistributionDefinitionNFT:
		AssertDistributionDefNFT(t, actual.Impl.(*bennyfi.DistributionDefinitionNFT), v)
	default:
		assert.Assert(t, false, "Unknown distribution definition type: %T", expected)
	}
}

func AssertCategory(t *testing.T, actual, expected *bennyfi.Category) {
	assert.Check(t, actual != nil)
	assert.Equal(t, actual.Category, expected.Category)
	assert.Equal(t, actual.CategoryName, expected.CategoryName)
	assert.Equal(t, actual.CategoryDescription, expected.CategoryDescription)
}

func AssertYieldSource(t *testing.T, actual, expected *bennyfi.YieldSource) {
	assert.Check(t, actual != nil)
	assert.Equal(t, actual.YieldSource, expected.YieldSource)
	assert.Equal(t, actual.YieldSourceName, expected.YieldSourceName)
	assert.Equal(t, actual.YieldSourceDescription, expected.YieldSourceDescription)
	assert.Equal(t, actual.StakeSymbol, expected.StakeSymbol)
	assert.Equal(t, actual.AdaptorContract, expected.AdaptorContract)
	assert.Equal(t, actual.YieldSourceCID, expected.YieldSourceCID)
	assert.Equal(t, actual.EntryFeePercentageOfYieldx100000, expected.EntryFeePercentageOfYieldx100000)
	assert.Equal(t, actual.DailyYieldx100000, expected.DailyYieldx100000)
	assert.Equal(t, actual.TokenValue, expected.TokenValue)
	assert.Equal(t, actual.BenyValue, expected.BenyValue)
}

func AssertDistributionDefFT(t *testing.T, actual, expected *bennyfi.DistributionDefinitionFT) {
	assert.Check(t, actual != nil)
	assert.Equal(t, actual.AllParticipantsPerc, expected.AllParticipantsPerc)
	assert.Equal(t, actual.RoundManagerPerc, expected.RoundManagerPerc)
	assert.Equal(t, actual.BeneficiaryPerc, expected.BeneficiaryPerc)
	assert.DeepEqual(t, actual.WinnersPerc, expected.WinnersPerc)
	AssertBaseDistributionDef(t, actual.BaseDistributionDefinition, expected.BaseDistributionDefinition)
}

func AssertDistributionDefNFT(t *testing.T, actual, expected *bennyfi.DistributionDefinitionNFT) {
	assert.Check(t, actual != nil)
	assert.Equal(t, actual.EachParticipantReward, expected.EachParticipantReward)
	assert.Equal(t, actual.RoundManagerFee, expected.RoundManagerFee)
	assert.Equal(t, actual.BeneficiaryReward, expected.BeneficiaryReward)
	assert.DeepEqual(t, actual.WinnerPrizes, expected.WinnerPrizes)
	AssertNFTConfig(t, actual.NFTConfig, expected.NFTConfig)
	AssertBaseDistributionDef(t, actual.BaseDistributionDefinition, expected.BaseDistributionDefinition)
}

func AssertBaseDistributionDef(t *testing.T, actual, expected *bennyfi.BaseDistributionDefinition) {
	assert.Check(t, actual != nil)
	AssertVestingConfig(t, actual.VestingConfig, expected.VestingConfig)
}

func AssertVestingConfig(t *testing.T, actual, expected *bennyfi.VestingConfig) {
	assert.Check(t, actual != nil)
	AssertConfig(t, actual.Config, expected.Config)
}

func AssertNFTConfig(t *testing.T, actual, expected *bennyfi.NFTConfig) {
	assert.Check(t, actual != nil)
	AssertConfig(t, actual.Config, expected.Config)
}

func AssertConfig(t *testing.T, actual, expected bennyfi.Config) {
	assert.Equal(t, len(actual), len(expected), "Different number of config parameters, actual: %v, expected: %v", len(actual), len(expected))
	for _, ev := range expected {
		AssertConfigEntry(t, actual.FindEntry(ev.Key), ev)
	}
}

func AssertConfigEntry(t *testing.T, actual, expected *bennyfi.ConfigEntry) {
	assert.Check(t, actual != nil)
	assert.Equal(t, actual.Key, expected.Key)
	assert.DeepEqual(t, actual.Value, expected.Value)
}

func AssertDistributions(t *testing.T, actual, expected bennyfi.Distributions) {
	assert.Equal(t, len(actual), len(expected), "Number of distributions do not match, actual: %v, expected: %v", len(actual), len(expected))
	for _, ee := range expected {
		ae := actual.Find(ee.Key)
		assert.Assert(t, ae != nil, "Expected distribution: %v, not found", ee.Key)
		AssertDistribution(t, ae.Value, ee.Value)
	}
}

func AssertDistribution(t *testing.T, actual, expected *bennyfi.Distribution) {
	assert.Check(t, actual != nil)
	assert.DeepEqual(t, actual, expected)
}

func AssertReturns(t *testing.T, actual, expected bennyfi.ReturnEntries) {
	assert.Equal(t, len(actual), len(expected), "Number of returns do not match, actual: %v, expected: %v", len(actual), len(expected))
	for _, ee := range expected {
		ae := actual.Find(ee.Key)
		assert.Assert(t, ae != nil, "Expected return: %v, not found", ee.Key)
		AssertReturn(t, ae.Value, ee.Value)
	}
}

func AssertReturn(t *testing.T, actual, expected *bennyfi.Returns) {
	assert.Check(t, actual != nil)
	assert.DeepEqual(t, actual, expected)

}

func AssertRewards(t *testing.T, actual, expected bennyfi.Rewards) {
	assert.Equal(t, len(actual), len(expected), "Number of rewards do not match, actual: %v, expected: %v", len(actual), len(expected))
	for _, ee := range expected {
		ae := actual.Find(ee.Key)
		assert.Assert(t, ae != nil, "Expected reward: %v, not found", ee.Key)
		AssertReward(t, ae.Value, ee.Value)
	}
}

func AssertReward(t *testing.T, actual, expected *bennyfi.Reward) {
	assert.Check(t, actual != nil)
	assert.DeepEqual(t, actual, expected)
}

func AssertTime(t *testing.T, actual, expected eos.TimePoint) {
	assert.Assert(t, math.Abs(float64(actual)-float64(expected)) < float64(1000))
}
