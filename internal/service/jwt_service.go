package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtService struct {
	Key            []byte
	ExpirationTime time.Duration
}

func NewJwtService(key string, exp time.Duration) *JwtService {
	var jwtService JwtService = JwtService{
		Key:            []byte(key),
		ExpirationTime: exp,
	}
	return &jwtService
}

type Claims struct {
	UserId int64    `json:"user_id"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

func (j *JwtService) GenerateJwtToken(id int64, roles []string) (string, error) {
	var claims Claims = Claims{
		UserId: id,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ExpirationTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "GoBank",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.Key)
}
