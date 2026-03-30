package routes

import (
	"wms/core/middlewares"
	"wms/internal/controller"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

var (
	BuildTime     = "unknown"
	Version       = "development"
	Commit        = "unknown"
	CommitMessage = "unknown"
)

type RouteConfig struct {
	App             *fiber.App
	Config          *viper.Viper
	AuthMiddleware  *middlewares.AuthMiddleware
	OrderController *controller.OrderController
	AuthController  *controller.AuthController
}

func (c *RouteConfig) Setup() {
	// c.App.Use(c.LogMiddleware)

	c.AuthRoute()
	c.WebhookRoute()
	c.OrderRoute(c.AuthMiddleware)
}
