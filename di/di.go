package di

import (
	"context"
	"go.uber.org/dig"
	"google-photos-sync/infrastructure/db"
	router "google-photos-sync/infrastructure/web"
	"google-photos-sync/pkg/closer"
	"google-photos-sync/pkg/env"
	"google-photos-sync/pkg/graceful"
	"google-photos-sync/usecases/controller"
	"google-photos-sync/usecases/factory"
	"google-photos-sync/usecases/interactor"
	"google-photos-sync/usecases/machine"
	"google-photos-sync/usecases/repository"
	"log"
)

type ControllerGroup struct {
	dig.In

	Controllers []router.Controller `group:"controller"`
}

func Build() *dig.Container {
	container := dig.New()

	provideOrPanic(container, factory.MakeSyncMachineFactory)
	// API
	provideOrPanic(container, controller.MakeAuthController, dig.Group("controller"))
	provideOrPanic(container, controller.MakeStateController, dig.Group("controller"))
	provideOrPanic(container, func(group ControllerGroup) []router.Controller { return group.Controllers })
	provideOrPanic(container, router.MakeEchoWebServer)
	// DB
	provideOrPanic(container, func() *db.PostgresConfig {
		return &db.PostgresConfig{Schema: "photos", Connection: env.EnvOrDefault("PG_CONNECTION", "localhost")}
	})
	provideOrPanic(container, db.MakePostgresConnection)
	provideOrPanic(container, repository.MakeSqlCredentialsRepository)
	provideOrPanic(container, repository.MakeSqlStateRepository)
	provideOrPanic(container, machine.MakeSqlStorage)
	// UseCases
	provideOrPanic(container, interactor.MakePhotosInteractor)
	provideOrPanic(container, closer.MakeCloserCollection)
	provideOrPanic(container, MakeApplication)

	return container
}

func provideOrPanic(container *dig.Container, constructor interface{}, opts ...dig.ProvideOption) {
	err := container.Provide(constructor, opts...)
	if err != nil {
		log.Fatalf("container.Provide: %v", err)
	}
}

type App struct {
	*graceful.Graceful
}

func MakeApplication(collection *closer.CloserCollection, server router.WebServer) graceful.Application {

	var graceful = &graceful.Graceful{
		StartAction: func() error {
			return server.ListenAndServe()
		},
		DeferAction: func(ctx context.Context) error {
			return collection.Close(ctx)
		},
		ShutdownAction: func(ctx context.Context) error {
			return nil
		},
	}

	return &App{
		Graceful: graceful,
	}
}
