package states

type StateType = int

const (
	Empty StateType = iota
	ReadyToSync
	Sync
	Complete
)
