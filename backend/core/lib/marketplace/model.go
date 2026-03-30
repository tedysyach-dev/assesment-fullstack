package marketplace

import (
	"encoding/json"
	"time"
)

type BaseResponse[T any] struct {
	Message string    `json:"message"`
	Data    *T        `json:"data"`
	Error   *[]string `json:"error"`
}

type Authorize struct {
	Code   string `json:"code"`
	ShopID string `json:"shop_id"`
	State  string `json:"state"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type tokenRequest struct {
	GrantType    string `json:"grant_type"`
	Code         string `json:"code,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type Session struct {
	AccessToken  string
	RefreshToken string
	ExpiredAt    time.Time
}

type Order struct {
	OrderSn        string    `json:"order_sn"`
	ShopID         string    `json:"shop_id"`
	Status         string    `json:"status"`
	ShippingStatus string    `json:"shipping_status"`
	TrackingNumber string    `json:"tracking_number"`
	Items          []Item    `json:"items"`
	TotalAmount    float64   `json:"total_amount"`
	CreatedAt      time.Time `json:"created_at"`
}

func (o *Order) ToJson() json.RawMessage {
	b, _ := json.Marshal(o)
	return b
}

type Item struct {
	Sku      string  `json:"sku"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type ShipOrderRequest struct {
	OrderSn   string
	ChannelId string
}

type ShipOrderResponse struct {
	OrderSn        string `json:"order_sn"`
	TrackingNo     string `json:"tracking_no"`
	ShippingStatus string `json:"shipping_status"`
}

type LogisticChannelResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}
