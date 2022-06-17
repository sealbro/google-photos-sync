package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"google-photos-sync/infrastructure/web/auth"
	"google-photos-sync/pkg/logger"
)

type WebServer interface {
	ListenAndServe() error
}

type Controller interface {
	RegisterRoutes(e *echo.Group)
}

type Route struct {
	Method  string
	Path    string
	Handler echo.HandlerFunc
}

type EchoWebServer struct {
	echo      *echo.Echo
	port      string
	apiPrefix string
}

func MakeEchoWebServer(controllers []Controller) WebServer {
	server := &EchoWebServer{
		echo:      echo.New(),
		port:      "8080",
		apiPrefix: "/api",
	}

	server.echo.GET("/swagger/*", echoSwagger.WrapHandler)

	// Auth Middleware
	server.echo.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// TODO from cookies or ENV
			userId := "user"
			ctx := auth.WrapWithUserId(c.Request().Context(), userId)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	})
	server.echo.Use(middleware.LoggerWithConfig(logger.EchoLoggerConfig))
	server.echo.Use(middleware.Recover())

	for _, c := range controllers {
		c.RegisterRoutes(server.echo.Group(server.apiPrefix))
	}

	return server
}

func (server *EchoWebServer) ListenAndServe() error {
	return server.echo.Start(":" + server.port)
}
