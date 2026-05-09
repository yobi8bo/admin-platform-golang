package jwt

import (
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID  uint   `json:"userId"`
	TokenID string `json:"tokenId"`
	jwtlib.RegisteredClaims
}

func Sign(secret string, userID uint, tokenID string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:  userID,
		TokenID: tokenID,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwtlib.NewNumericDate(now),
		},
	}
	return jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims).SignedString([]byte(secret))
}

func Parse(secret, token string) (*Claims, error) {
	claims := &Claims{}
	parsed, err := jwtlib.ParseWithClaims(token, claims, func(token *jwtlib.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !parsed.Valid {
		return nil, jwtlib.ErrTokenInvalidClaims
	}
	return claims, nil
}
