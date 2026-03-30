package utils

import (
	"context"
	"fmt"
	"time"
	"wms/internal/model"

	"github.com/golang-jwt/jwt/v5"
)

type TokenUtil struct {
	SecretKey string
}

func NewTokenUtil(secretKey string) *TokenUtil {
	return &TokenUtil{
		SecretKey: secretKey,
	}
}

func (t *TokenUtil) CreateAccessToken(ctx context.Context, auth *model.Auth) (string, error) {
	exp := time.Now().Add(time.Minute * 30).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": auth.UID,
		"exp": exp, // fix: was "expire"
	})

	jwtToken, err := token.SignedString([]byte(t.SecretKey))
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

func (t *TokenUtil) ParseToken(jwtToken string) (*model.TokenClaims, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(t.SecretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("Unauthorize")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	uid, ok := claims["uid"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid uid claim")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid exp claim")
	}

	return &model.TokenClaims{
		UID: uid,
		Exp: int64(exp),
	}, nil
}
