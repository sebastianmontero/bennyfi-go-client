package bennyfi

import (
	eos "github.com/sebastianmontero/eos-go"
)

const (
	SettingEntryFeeAccount               = "ENTRY_FEE_ACCOUNT"
	SettingEntryFeeBurnAccount           = "ENTRY_FEE_BURN_ACCOUNT"
	SettingEntryFeePercentageOfYield     = "ENTRY_FEE_PERCENTAGE_OF_YIELD"
	SettingEntryTokenTelosYieldDaily     = "ENTRY_TOKEN_TELOS_YIELD_DAILY"
	SettingEntryTokenValueTLOS           = "ENTRY_TOKEN_VALUE_TLOS"
	SettingEntryTokenValueBENY           = "ENTRY_TOKEN_VALUE_BENY"
	SettingEntryFeeSelffundedPeruserBeny = "ENTRY_FEE_SELFFUNDED_PERUSER_BENY"
	SettingEntryFeeRefundOnCancelPm      = "ENTRY_FEE_REFUND_ON_CANCEL_PM"
	SettingEntryFeeBurnYes               = "ENTRY_FEE_BURN_YES"
	SettingBenyToken                     = "BENY_TOKEN"
	SettingRoundManagerStakeAmount       = "POOL_MANAGER_STAKE_AMOUNT"
	SettingRoundManagerStakeRate         = "POOL_MANAGER_STAKE_RATE"
	SettingBeneficiaryStakeAmount        = "BENEFICIARY_STAKE_AMOUNT"
	SettingCreationFee                   = "CREATION_FEE"
	SettingPoolManagerUnstakingPeriodHrs = "POOL_MANAGER_UNSTAKING_PERIOD_HRS"
	SettingBeneficiaryUnstakingPeriodHrs = "BENEFICIARY_UNSTAKING_PERIOD_HRS"
	SettingIsPaused                      = "IS_PAUSED"
)

type EntryFeeSettings struct {
	PercOfYield       uint32
	DailyYield        uint32
	ValueTLOS         eos.Asset
	ValueBENY         eos.Asset
	SelfFundedPerUser eos.Asset
	BENYToken         eos.Asset
}

func (m *EntryFeeSettings) HourlyYield() uint32 {
	return m.DailyYield / 24
}

func GetEntryFeeSettings(contract *BennyfiContract) (*EntryFeeSettings, error) {
	percOfYield, err := contract.SettingAsUint32(SettingEntryFeePercentageOfYield)
	if err != nil {
		return nil, err
	}

	dailyYield, err := contract.SettingAsUint32(SettingEntryTokenTelosYieldDaily)
	if err != nil {
		return nil, err
	}

	valueTLOS, err := contract.SettingAsAsset(SettingEntryTokenValueTLOS)
	if err != nil {
		return nil, err
	}

	valueBENY, err := contract.SettingAsAsset(SettingEntryTokenValueBENY)
	if err != nil {
		return nil, err
	}
	selfFundedPerUser, err := contract.SettingAsAsset(SettingEntryFeeSelffundedPeruserBeny)
	if err != nil {
		return nil, err
	}
	benyToken, err := contract.SettingAsAsset(SettingBenyToken)
	if err != nil {
		return nil, err
	}
	return &EntryFeeSettings{
		PercOfYield:       percOfYield,
		DailyYield:        dailyYield,
		ValueTLOS:         valueTLOS,
		ValueBENY:         valueBENY,
		SelfFundedPerUser: selfFundedPerUser,
		BENYToken:         benyToken,
	}, nil
}

func (m *BennyfiContract) ShouldBurnFees() (bool, error) {
	shouldBurn, err := m.SettingAsUint32(SettingEntryFeeBurnYes)
	if err != nil {
		return false, err
	}
	return shouldBurn > 0, nil
}

func (m *BennyfiContract) GetActiveFeeAccount() (eos.AccountName, error) {
	shouldBurn, err := m.ShouldBurnFees()
	if err != nil {
		return "", err
	}
	feeAccountSettingName := SettingEntryFeeAccount
	if shouldBurn {
		feeAccountSettingName = SettingEntryFeeBurnAccount
	}
	feeAccount, err := m.SettingAsName(feeAccountSettingName)
	if err != nil {
		return "", err
	}
	return eos.AccountName(feeAccount), nil
}

func (m *BennyfiContract) GetAuthUnstakeWaitingPeriodHrs(authLevel uint64) (uint32, error) {
	settingKey := SettingPoolManagerUnstakingPeriodHrs
	if authLevel == Beneficiary {
		settingKey = SettingBeneficiaryUnstakingPeriodHrs
	}
	return m.SettingAsUint32(settingKey)
}
