package jwtservice

import (
	"GoBank/internal/domain"
	"GoBank/internal/repository"
	"context"
)

type JwtService struct {
	repo repository.UserRepository
}

func NewJwtService(repo repository.UserRepository) *JwtService {
	var jwtService JwtService = JwtService{
		repo: repo,
	}
	return &jwtService
}

func (j *JwtService) Login(ctx context.Context, email string, password string) {

}

func (j *JwtService) GetUserDetails(ctx context.Context, email string) (domain.User, error) {
	user, err := j.repo.CheckUserInDataBase(ctx, email)
	if err != nil {
		return user, err
	}
	return user, nil
}
