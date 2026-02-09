package common

type StoppedStake struct {
	RoundId   uint64
	IsStopped bool
}

type StoppedInterface interface {
	GetAllStoppedStakes() ([]StoppedStake, error)
}
