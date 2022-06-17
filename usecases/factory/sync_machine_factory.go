package factory

import (
	"context"
	"google-photos-sync/infrastructure/state"
	"google-photos-sync/infrastructure/web/auth"
	"google-photos-sync/usecases/machine"
)

type SyncMachineFactory struct {
	stateMachines map[string]state.Machine
	stateStorage  machine.StateStorage
}

func MakeSyncMachineFactory(storage machine.StateStorage) *SyncMachineFactory {
	return &SyncMachineFactory{
		stateMachines: make(map[string]state.Machine),
		stateStorage:  storage,
	}
}

func (f *SyncMachineFactory) Create(ctx context.Context) (state.Machine, error) {
	userId := auth.GetUserId(ctx)
	if machine, ok := f.stateMachines[userId]; ok {
		return machine, nil
	}

	syncMachine := machine.NewSyncMachine(userId, f.stateStorage)
	err := syncMachine.Load(ctx)

	return syncMachine, err
}
