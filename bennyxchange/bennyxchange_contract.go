package bennyxchange

import (
	eos "github.com/sebastianmontero/eos-go"
)

var (
	ExchangeEntry             = eos.Name("xentry")
	ExchangeBeneficiaryReward = eos.Name("xbenefreward")
	ExchangePoolManagerFee    = eos.Name("xpoolmgrfee")
)
