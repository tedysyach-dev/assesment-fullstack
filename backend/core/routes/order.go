package routes

import "wms/core/middlewares"

func (c *RouteConfig) OrderRoute(auth *middlewares.AuthMiddleware) {
	order := c.App.Group("/order")

	// semua route wajib login
	order.Use(auth.Authenticate())

	order.Post("/:order_sn/pick",
		auth.Authorize("PICKER"),
		c.OrderController.PickOrder,
	)

	order.Post("/:order_sn/pack",
		auth.Authorize("PACKER"),
		c.OrderController.PackOrder,
	)

	order.Post("/:order_sn/ship",
		auth.Authorize("ADMIN"),
		c.OrderController.ShipOrder,
	)

	// semua role boleh akses
	order.Get("/:order_sn", c.OrderController.OrderDetail)
	order.Get("/", c.OrderController.OrderList)
}
