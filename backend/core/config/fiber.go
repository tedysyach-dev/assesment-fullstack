package config

import (
	"log"
	"wms/core/errors"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func NewFiber(config *viper.Viper) *fiber.App {
	var app = fiber.New(fiber.Config{
		AppName:      config.GetString("app.name"),
		ErrorHandler: NewErrorHandler(config),
		Prefork:      config.GetBool("web.prefork"),
	})

	return app
}

func NewErrorHandler(config *viper.Viper) fiber.ErrorHandler {
	isDevelopment := config.GetString("app.env") == "development"

	return func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		message := "Internal server error"
		var details interface{}

		traceID := ctx.Locals("trace_id")
		if traceID == nil {
			traceID = ctx.Get("X-Request-ID")
		}

		if appErr, ok := err.(*errors.AppError); ok {
			code = appErr.Code
			message = appErr.Message
			details = appErr.Details

			if appErr.Internal != nil {
				log.Printf("[ERROR] [TraceID: %v] %s: %v", traceID, message, appErr.Internal)
			}
		} else if fiberErr, ok := err.(*fiber.Error); ok {
			code = fiberErr.Code
			message = fiberErr.Message
		} else {
			log.Printf("[ERROR] [TraceID: %v] Unexpected error: %v", traceID, err)

			if !isDevelopment {
				message = "An unexpected error occurred"
			} else {
				message = err.Error()
			}
		}

		response := errors.ErrorResponse{
			Status:  false,
			Message: message,
			Errors:  details,
		}

		if isDevelopment && traceID != nil {
			response.TraceID = traceID.(string)
		}

		return ctx.Status(code).JSON(response)
	}
}
