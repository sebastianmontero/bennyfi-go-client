package common

type BasicStake struct {
	RoundId   uint64
	IsStopped bool
}

type StoppedInterface interface {
	GetAllStoppedStakes() ([]BasicStake, error)
}
