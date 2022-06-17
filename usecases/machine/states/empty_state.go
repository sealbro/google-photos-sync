package states

import (
	"context"
)

type EmptyState struct {
	Type   int
	UserId string
}

func NewEmptyState(userId string) EmptyState {
	return EmptyState{
		Type:   Empty,
		UserId: userId,
	}
}

func (state *EmptyState) Action(ctx context.Context) error {
	return nil
}

func (state *EmptyState) GetType() int {
	return state.Type
}
