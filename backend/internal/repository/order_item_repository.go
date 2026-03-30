package repository

import (
	"context"
	"wms/core/base"
	"wms/internal/entity"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

type OrderItemRepository struct {
	base.Repository[entity.OrderItem]
	Log *logrus.Logger
}

func NewOrderItemRepository(db *bun.DB, log *logrus.Logger) *OrderItemRepository {
	return &OrderItemRepository{
		Repository: base.Repository[entity.OrderItem]{DB: db},
		Log:        log,
	}
}

func (r *OrderItemRepository) DeleteByOrderSNs(ctx context.Context, orderSNs []string) error {
	_, err := r.IDB(ctx).NewDelete().
		Model((*entity.OrderItem)(nil)).
		Where("order_sn IN (?)", bun.List(orderSNs)).
		Exec(ctx)
	return err
}
