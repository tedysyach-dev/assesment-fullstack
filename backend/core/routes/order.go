package routes

import "github.com/gofiber/fiber/v2"

func (c *RouteConfig) OrderRoute(authMiddleware fiber.Handler) {
	order := c.App.Group("/order")
	order.Use(authMiddleware)
	order.Post("/:order_sn/pick", c.OrderController.PickOrder)
	order.Post("/:order_sn/pack", c.OrderController.PackOrder)
	order.Post("/:order_sn/ship", c.OrderController.ShipOrder)
	order.Get("/:order_sn", c.OrderController.OrderDetail)
	order.Get("/", c.OrderController.OrderList)
}
