package model

type TokenClaims struct {
	UID  string `json:"userId"`
	ROLE string `json:"role"`
	Exp  int64  `json:"exp"`
}

type Auth struct {
	UID    string
	ROLE   string
	Claims *TokenClaims
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role" validate:"required,oneof=ADMIN STAFF PICKER PACKER"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
	Role        string `json:"role"`
}
