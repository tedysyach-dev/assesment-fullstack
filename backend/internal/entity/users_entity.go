package entity

import "github.com/uptrace/bun"

type Users struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID           string `bun:"id,pk"`
	Email        string `bun:"email,unique,notnull"`
	PasswordHash string `bun:"password_hash,notnull"`
}
