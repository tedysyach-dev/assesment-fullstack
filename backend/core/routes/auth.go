package routes

func (c *RouteConfig) AuthRoute() {
	auth := c.App.Group("/auth")
	auth.Post("/register", c.AuthController.Register)
	auth.Post("/login", c.AuthController.Login)

	// auth.Post("/register", c.UserController.Register)
}
