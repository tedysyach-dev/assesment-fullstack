package main

import (
	"fmt"
	"log"
	"wms/core/config"
	"wms/core/lib/marketplace"
)

func main() {
	// Setup Viper
	v := config.NewViper()
	logger := config.NewLogger(v)
	validate := config.NewValidator()
	db := config.NewBun(v, logger)
	defer db.Close()
	app := config.NewFiber(v)
	store := marketplace.NewMemorySessionStore()
	sessionKey := "shope-123"
	marketplaceClient, err := marketplace.NewClient(v,
		marketplace.WithSessionStore(store, sessionKey),
		marketplace.WithShopID("shopee-123"),
		marketplace.WithLogger(logger),
	)

	if err != nil {
		log.Fatalf("Failed to integrate with marketplace: %v", err)
	}

	config.Bootstrap(&config.BootstrapConfig{
		DB:                db,
		App:               app,
		Log:               logger,
		Validate:          validate,
		Config:            v,
		MarketplaceClient: marketplaceClient,
	})

	logger.Info("Application started")
	webPort := v.GetInt("web.port")
	err = app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
