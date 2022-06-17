package state

import (
	"context"
	"google-photos-sync/usecases/machine/states"
)

type State interface {
	Action(ctx context.Context) error
	GetType() int
}

type Machine interface {
	SetState(ctx context.Context, stateType states.StateType) error
	Action(ctx context.Context) error
}
