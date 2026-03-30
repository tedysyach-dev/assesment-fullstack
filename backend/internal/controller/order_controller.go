package controller

import (
	"wms/core/errors"
	"wms/core/utils"
	"wms/internal/model"
	"wms/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type OrderController struct {
	Log     *logrus.Logger
	Service *service.OrderService
}

func NewOrderController(service *service.OrderService, logger *logrus.Logger) *OrderController {
	return &OrderController{
		Log:     logger,
		Service: service,
	}
}

func (c *OrderController) WebhookOrderStatus(ctx *fiber.Ctx) error {
	request := new(model.WebhookOrderStatusBody)

	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return errors.NewBadRequestError("Invalid request body", err)
	}
	err := c.Service.WebhookOrderStatus(ctx.Context(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(utils.WebResponse[*string]{Status: true, Message: "Success", Data: nil})
}

func (c *OrderController) WebhookShippingStatus(ctx *fiber.Ctx) error {
	request := new(model.WebhookShippingStatusBody)

	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return errors.NewBadRequestError("Invalid request body", err)
	}
	err := c.Service.WebhookShippingStatus(ctx.Context(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(utils.WebResponse[*string]{Status: true, Message: "Success", Data: nil})
}

func (c *OrderController) OrderList(ctx *fiber.Ctx) error {

	err := c.Service.SyncOrders(ctx.Context())
	query := ctx.Query("wmsStatus")

	if err != nil {
		c.Log.WithError(err).Warn("[order] failed to get order list")
	}

	res := c.Service.OrderList(ctx.Context(), query)

	return ctx.JSON(utils.WebResponse[*[]model.OrderResponse]{Status: true, Message: "Success", Data: &res})
}

func (c *OrderController) OrderDetail(ctx *fiber.Ctx) error {
	orderSn := ctx.Params("order_sn")
	res, err := c.Service.OrderDetail(ctx.Context(), orderSn)
	if err != nil {
		return err
	}

	return ctx.JSON(utils.WebResponse[*model.OrderResponse]{Status: true, Message: "Success", Data: res})
}

func (c *OrderController) PickOrder(ctx *fiber.Ctx) error {
	orderSn := ctx.Params("order_sn")
	err := c.Service.PickOrder(ctx.Context(), orderSn)
	if err != nil {
		return err
	}

	return ctx.JSON(utils.WebResponse[*string]{Status: true, Message: "Success", Data: nil})
}

func (c *OrderController) PackOrder(ctx *fiber.Ctx) error {
	orderSn := ctx.Params("order_sn")
	err := c.Service.PackOrder(ctx.Context(), orderSn)
	if err != nil {
		return err
	}

	return ctx.JSON(utils.WebResponse[*string]{Status: true, Message: "Success", Data: nil})
}

func (c *OrderController) ShipOrder(ctx *fiber.Ctx) error {
	request := new(model.ShipOrderRequest)
	orderSn := ctx.Params("order_sn")

	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return errors.NewBadRequestError("Invalid request body", err)
	}

	res, err := c.Service.ShipOrder(ctx.Context(), orderSn, request)
	if err != nil {
		return err
	}

	return ctx.JSON(utils.WebResponse[*model.ShipOrderResponse]{Status: true, Message: "Success", Data: res})
}
