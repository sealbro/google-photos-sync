package states

import (
	"context"
)

type CompleteState struct {
	EmptyState
}

func NewCompleteState(state EmptyState) *CompleteState {
	state.Type = Complete

	return &CompleteState{
		EmptyState: state,
	}
}

func (state *CompleteState) Action(ctx context.Context) error {

	return nil
}
