package controller

import (
	"github.com/labstack/echo/v4"
	router "google-photos-sync/infrastructure/web"
	"google-photos-sync/usecases/interactor"
)

type StateController struct {
	photosInteractor *interactor.PhotosInteractor
}

func MakeStateController(photosInteractor *interactor.PhotosInteractor) router.Controller {
	return &StateController{
		photosInteractor: photosInteractor,
	}
}

func (c *StateController) RegisterRoutes(e *echo.Group) {
	group := e.Group("/v1/state")

	group.POST("/sync", c.syncHandler)
	group.POST("/complete", c.completeHandler)
	group.GET("/stats", c.statsHandler)

}

// Set Sync state
// @Tags state
// @Accept json
// @Produce json
// @Description Check and set Sync state
// @Success 200
// @Router /v1/state/sync [post]
func (c *StateController) syncHandler(e echo.Context) error {
	ctx := e.Request().Context()

	return c.photosInteractor.Sync(ctx)
}

// Set Complete state
// @Tags state
// @Accept json
// @Produce json
// @Description Check and set Complete state
// @Success 200
// @Router /v1/state/complete [post]
func (c *StateController) completeHandler(e echo.Context) error {
	ctx := e.Request().Context()

	return c.photosInteractor.Complete(ctx)
}

// Get current statistics
// @Tags state
// @Accept json
// @Produce json
// @Description Get current sync statistics
// @Success 200
// @Router /v1/state/stats [get]
func (c *StateController) statsHandler(e echo.Context) error {
	ctx := e.Request().Context()

	return c.photosInteractor.Statistic(ctx)
}
