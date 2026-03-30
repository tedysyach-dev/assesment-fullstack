package repository

import (
	"wms/core/base"
	"wms/internal/entity"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

type OrderRepository struct {
	base.Repository[entity.Order]
	Log *logrus.Logger
}

func NewOrderRepository(db *bun.DB, log *logrus.Logger) *OrderRepository {
	return &OrderRepository{
		Repository: base.Repository[entity.Order]{DB: db},
		Log:        log,
	}
}
