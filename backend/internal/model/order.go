package model

import (
	"time"
	"wms/internal/entity"
)

type OrderResponse struct {
	ID                string           `json:"id"`
	OrderSN           string           `json:"order_sn"`
	ShopID            string           `json:"shop_id"`
	MarketplaceStatus string           `json:"marketplace_status"`
	ShippingStatus    string           `json:"shipping_status"`
	WmsStatus         entity.WmsStatus `json:"wms_status"`
	TrackingNumber    *string          `json:"tracking_number,omitempty"`
	OrderItem         []OrderItemRes   `json:"items,omitempty"`
	TotalAmount       float64          `json:"total_amount"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
}

type OrderItemRes struct {
	ID        string    `json:"id"`
	OrderID   string    `json:"order_id"`
	SKU       string    `json:"sku"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

type ShipOrderResponse struct {
	OrderSn        string           `json:"order_sn"`
	WmsStatus      entity.WmsStatus `json:"wms_status"`
	ShippingStatus string           `json:"shipping_status"`
	TrackingNumber string           `json:"tracking_number"`
}

type ShipOrderRequest struct {
	ChannelId string `json:"channelId" validate:"required"`
}

type WebhookOrderStatusBody struct {
	Message string                 `json:"message"`
	Data    WebhookOrderStatusData `json:"data"`
}
type WebhookOrderStatusData struct {
	OrderSn string `json:"order_sn"`
	Status  string `json:"status"`
}

type WebhookShippingStatusBody struct {
	Message string                    `json:"message"`
	Data    WebhookShippingStatusData `json:"data"`
}

type WebhookShippingStatusData struct {
	OrderSn       string `json:"order_sn"`
	ShippingState string `json:"shipping_state"`
}
