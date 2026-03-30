package repository

import (
	"wms/core/base"
	"wms/internal/entity"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

type UsersRepository struct {
	base.Repository[entity.Users]
	Log *logrus.Logger
}

func NewUsersRepository(db *bun.DB, log *logrus.Logger) *UsersRepository {
	return &UsersRepository{
		Repository: base.Repository[entity.Users]{DB: db},
		Log:        log,
	}
}
