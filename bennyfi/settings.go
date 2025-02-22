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
	SettingMaxEntriesPerParticipant      = "MAX_ENTRIES_PER_PARTICIPANT"
	SettingMaxPoolStartPrdHrs            = "MAX_POOL_START_PRD_HRS"
)

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
