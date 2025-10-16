package service

import (
	"GoBank/internal/repository"
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtService struct {
	Key            []byte
	ExpirationTime time.Duration
	repo           *repository.UserRepo
}

func NewJwtService(key string, exp time.Duration, repo *repository.UserRepo) *JwtService {
	var jwtService JwtService = JwtService{
		Key:            []byte(key),
		ExpirationTime: exp,
		repo:           repo,
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

func (j *JwtService) TokenValidation(token string) bool {
	var jwtKey []byte = j.Key
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(jwtKey), nil
	})
	if err != nil {
		return false
	}

	if !parsedToken.Valid {
		return false
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return false
	}
	if time.Now().Unix() > int64(exp) {
		return false
	}
	signer, ok := claims["iss"].(string)
	if !ok {
		return false
	}
	if signer != "GoBank" {
		return false
	}

	userId, ok := claims["user_id"].(float64)
	if !ok {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if find := j.repo.CheckUserInDataBaseById(ctx, int64(userId)); !find {
		return false
	}
	return true
}
