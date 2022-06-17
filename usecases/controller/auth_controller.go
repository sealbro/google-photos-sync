package controller

import (
	"github.com/labstack/echo/v4"
	"google-photos-sync/domain/model"
	"google-photos-sync/infrastructure/google"
	router "google-photos-sync/infrastructure/web"
	"google-photos-sync/usecases/interactor"
	"net/http"
)

type AuthController struct {
	photosInteractor *interactor.PhotosInteractor
}

func MakeAuthController(photosInteractor *interactor.PhotosInteractor) router.Controller {
	return &AuthController{photosInteractor: photosInteractor}
}

func (c *AuthController) RegisterRoutes(e *echo.Group) {
	group := e.Group("/v1/auth")
	group.GET("/callback", c.callbackHandler)
	group.GET("/account/from", c.getUrlFromHandler)
	group.GET("/account/to", c.getUrlToHandler)
}

// Callback
// @Tags auth
// @Accept json
// @Produce json
// @Param state query string true "account type [from / to]"
// @Param code query string true "google response code"
// @Description Google callback endpoint
// @Success 200
// @Router /v1/auth/callback [get]
func (c *AuthController) callbackHandler(e echo.Context) error {
	// TODO validate
	var request = &model.GoogleCallback{}
	e.Bind(request)
	ctx := e.Request().Context()

	return c.photosInteractor.SaveAccount(ctx, request)
}

// GetUrlFrom
// @Tags auth
// @Accept json
// @Produce json
// @Description Get url from account
// @Success 200 {string} string
// @Router /v1/auth/account/from [get]
func (c *AuthController) getUrlFromHandler(e echo.Context) error {
	url, err := google.GetAuthUrl(model.From)
	if err != nil {
		return err
	}

	e.String(http.StatusMovedPermanently, url)

	return nil
}

// GetUrlTo
// @Tags auth
// @Accept json
// @Produce json
// @Description Get url to account
// @Success 200 {string} string
// @Router /v1/auth/account/to [get]
func (c *AuthController) getUrlToHandler(e echo.Context) error {
	url, err := google.GetAuthUrl(model.To)
	if err != nil {
		return err
	}

	e.String(http.StatusMovedPermanently, url)

	return nil
}
