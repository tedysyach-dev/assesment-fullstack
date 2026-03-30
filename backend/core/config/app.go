package config

import (
	"wms/core/lib/marketplace"
	"wms/core/middlewares"
	"wms/core/routes"
	"wms/core/utils"
	"wms/internal/controller"
	"wms/internal/repository"
	"wms/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
)

type BootstrapConfig struct {
	DB                *bun.DB
	App               *fiber.App
	Log               *logrus.Logger
	Validate          *validator.Validate
	Config            *viper.Viper
	MarketplaceClient *marketplace.Client
}

func Bootstrap(config *BootstrapConfig) {

	setupCORS(config)

	tokenUtils := utils.NewTokenUtil(config.Config.GetString("secret_key"))
	// Auth middleware configuration
	authConfig := middlewares.AuthMiddlewareConfig{
		TokenUtil: tokenUtils,
	}

	// Middleware
	authMiddleware := middlewares.NewAuthMiddleware(authConfig)

	//Repository
	orderRepository := repository.NewOrderRepository(config.DB, config.Log)
	orderItemRepository := repository.NewOrderItemRepository(config.DB, config.Log)
	userRepository := repository.NewUsersRepository(config.DB, config.Log)

	//Service
	orderService := service.NewOrderService(config.DB, config.Log, config.Validate, config.MarketplaceClient, orderRepository, orderItemRepository)
	authService := service.NewAuthService(config.DB, config.Log, config.Validate, userRepository, tokenUtils)

	//Controller
	orderController := controller.NewOrderController(orderService, config.Log)
	authController := controller.NewAuthController(authService, config.Log)

	// Setup routes
	routeConfig := routes.RouteConfig{
		App:             config.App,
		Config:          config.Config,
		AuthMiddleware:  authMiddleware,
		OrderController: orderController,
		AuthController:  authController,
	}

	routeConfig.Setup()
}

func setupCORS(configs *BootstrapConfig) {
	isDev := true

	if isDev {
		configs.App.Use(cors.New(cors.Config{
			AllowOrigins:     "*",
			AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
			AllowHeaders:     "Origin,Content-Type,Accept,Authorization,Platform", // Add Platform
			ExposeHeaders:    "Content-Length",
			AllowCredentials: false,
		}))
	} else {
		configs.App.Use(cors.New(cors.Config{
			AllowOrigins:     configs.Config.GetString("allowedWeb"),
			AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
			AllowHeaders:     "Origin,Content-Type,Accept,Authorization,Platform", // Add Platform
			ExposeHeaders:    "Content-Length",
			AllowCredentials: true,
		}))
	}
}
