package jwt

import (
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

// Claims 是系统 JWT 的业务载荷，TokenID 用于刷新令牌轮换和失效控制。
type Claims struct {
	UserID  uint   `json:"userId"`  // 登录用户 ID。
	TokenID string `json:"tokenId"` // 令牌唯一 ID，刷新令牌会用该值映射 Redis key。
	jwtlib.RegisteredClaims
}

// Sign 使用 HS256 签发 JWT，ttl 控制当前令牌的过期时间。
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

// Parse 校验 JWT 签名和有效期，并返回系统业务载荷。
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
