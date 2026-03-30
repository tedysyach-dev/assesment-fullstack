package routes

func (c *RouteConfig) WebhookRoute() {
	webhook := c.App.Group("/webhook")
	webhook.Post("/order-status", c.OrderController.WebhookOrderStatus)
	webhook.Post("/shipping-status", c.OrderController.WebhookShippingStatus)
}
