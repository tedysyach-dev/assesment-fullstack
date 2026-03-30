package controller

import (
	"wms/core/errors"
	"wms/core/utils"
	"wms/internal/model"
	"wms/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AuthController struct {
	Log     *logrus.Logger
	Service *service.AuthService
}

func NewAuthController(service *service.AuthService, logger *logrus.Logger) *AuthController {
	return &AuthController{
		Log:     logger,
		Service: service,
	}
}

func (c *AuthController) Register(ctx *fiber.Ctx) error {
	request := new(model.RegisterRequest)

	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return errors.NewBadRequestError("Invalid request body", err)
	}

	if err := c.Service.RegisterUser(ctx.Context(), request); err != nil {
		return err
	}

	return ctx.JSON(utils.WebResponse[*string]{Status: true, Message: "Register success", Data: nil})
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	request := new(model.LoginRequest)

	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return errors.NewBadRequestError("Invalid request body", err)
	}

	res, err := c.Service.Login(ctx.Context(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(utils.WebResponse[*model.LoginResponse]{Status: true, Message: "Login success", Data: res})
}
