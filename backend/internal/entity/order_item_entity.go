package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type OrderItem struct {
	bun.BaseModel `bun:"table:order_items,alias:oi"`

	ID        string    `bun:"id,pk"`
	OrderSn   string    `bun:"order_sn,notnull"`
	SKU       string    `bun:"sku,notnull"`
	Quantity  int       `bun:"quantity,notnull"`
	Price     float64   `bun:"price,notnull,type:decimal(15,2)"`
	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp"`

	// Relations
	Order *Order `bun:"rel:belongs-to,join:order_sn=order_sn"`
}
