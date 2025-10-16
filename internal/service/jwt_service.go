package service

import (
	"errors"
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

func (j *JwtService) TokenValidation(token string) (int64, bool) {
	var jwtKey []byte = j.Key
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(jwtKey), nil
	})
	if err != nil {
		return -1, false
	}

	if !parsedToken.Valid {
		return -1, false
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return -1, false
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return -1, false
	}
	if time.Now().Unix() > int64(exp) {
		return -1, false
	}
	signer, ok := claims["iss"].(string)
	if !ok {
		return -1, false
	}
	if signer != "GoBank" {
		return -1, false
	}

	userId, ok := claims["user_id"].(float64)
	if !ok {
		return -1, false
	}
	return int64(userId), true
}

func (j *JwtService) GetUserIdFromToken(token string) int64 {
	parsedToken, _ := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return []byte(j.Key), nil
	})
	claims, _ := parsedToken.Claims.(jwt.MapClaims)
	userId := claims["user_id"].(float64)
	return int64(userId)
}
