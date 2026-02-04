package util

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret   []byte
	issuer   string
	expireIn time.Duration
}

func NewJWTManager(secret, issuer string, expireIn time.Duration) *JWTManager {
	return &JWTManager{
		secret:   []byte(secret),
		issuer:   issuer,
		expireIn: expireIn,
	}
}

type Claims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func (m *JWTManager) Generate(userID int64, email string) (token string, expiresAt time.Time, err error) {
	now := time.Now()
	exp := now.Add(m.expireIn)

	claims := &Claims{
		UserID: userID,
		Email:  email,
		// Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   email,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := t.SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, exp, nil
}

func (m *JWTManager) Verify(tokenString string) (*Claims, error) {
	tok, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := tok.Claims.(*Claims)
	if !ok || !tok.Valid {
		return nil, errors.New("invalid token")
	}
	if claims.Issuer != m.issuer {
		return nil, errors.New("invalid issuer")
	}
	return claims, nil
}
