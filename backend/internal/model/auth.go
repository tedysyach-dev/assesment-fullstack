package model

type TokenClaims struct {
	UID string `json:"userId"`
	Exp int64  `json:"exp"`
}

type Auth struct {
	UID    string
	Claims *TokenClaims
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
}
