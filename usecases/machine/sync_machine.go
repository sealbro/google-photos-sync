package machine

import (
	"context"
	"fmt"
	"google-photos-sync/infrastructure/state"
	"google-photos-sync/infrastructure/web/auth"
	"google-photos-sync/usecases/machine/states"
)

type SyncMachine struct {
	current      state.State
	stateStorage StateStorage
}

func NewSyncMachine(userId string, storage StateStorage) *SyncMachine {
	emptyState := states.NewEmptyState(userId)
	return &SyncMachine{
		current:      &emptyState,
		stateStorage: storage,
	}
}

func (machine *SyncMachine) Load(ctx context.Context) error {
	lastState, err := machine.stateStorage.GetLastState(ctx)
	if err != nil {
		return err
	}

	machine.current = lastState

	return nil
}

func (machine *SyncMachine) Action(ctx context.Context) error {
	return machine.current.Action(ctx)
}

func (machine *SyncMachine) SetState(ctx context.Context, stateType states.StateType) error {
	switch machine.current.(type) {
	case *states.ReadyToSyncState:
		switch stateType {
		case states.Complete:
			return machine.setState(ctx, stateType)
		case states.Sync:
			return machine.setState(ctx, stateType)
		case states.ReadyToSync:
			return nil
		}
	case *states.SyncState:
		switch stateType {
		case states.Complete:
			return machine.setState(ctx, stateType)
		}
	case *states.CompleteState:
		switch stateType {
		case states.ReadyToSync:
			return machine.setState(ctx, stateType)
		}
	case *states.EmptyState:
		return machine.setState(ctx, stateType)
	}

	return fmt.Errorf("cann't change state from %v to %v", machine.current.GetType(), stateType)
}

func (machine *SyncMachine) setState(ctx context.Context, stateType states.StateType) error {
	if machine.current.GetType() == stateType {
		return nil
	}

	state := machine.stateStorage.Convert(auth.GetUserId(ctx), stateType)
	machine.current = state
	return machine.stateStorage.AddNewState(ctx, machine.current, state)
}
