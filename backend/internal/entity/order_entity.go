package entity

import (
	"encoding/json"
	"time"

	"github.com/uptrace/bun"
)

type WmsStatus string

const (
	WmsStatusReadyToPick WmsStatus = "READY_TO_PICK"
	WmsStatusPicking     WmsStatus = "PICKING"
	WmsStatusPacked      WmsStatus = "PACKED"
	WmsStatusShiped      WmsStatus = "SHIPED"
)

type Order struct {
	bun.BaseModel `bun:"table:orders,alias:o"`

	ID                    string          `bun:"id,pk"`
	OrderSN               string          `bun:"order_sn,notnull,unique"`
	ShopID                string          `bun:"shop_id,notnull"`
	MarketplaceStatus     string          `bun:"marketplace_status,notnull"`
	ShippingStatus        string          `bun:"shipping_status,notnull"`
	WmsStatus             WmsStatus       `bun:"wms_status,notnull"`
	TrackingNumber        *string         `bun:"tracking_number"`
	TotalAmount           float64         `bun:"total_amount,notnull,type:decimal(15,2)"`
	RawMarketplacePayload json.RawMessage `bun:"raw_marketplace_payload,type:jsonb"`
	CreatedAt             time.Time       `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt             time.Time       `bun:"updated_at,notnull,default:current_timestamp"`

	// Relations
	OrderItems []*OrderItem `bun:"rel:has-many,join:id=order_sn"`
}
