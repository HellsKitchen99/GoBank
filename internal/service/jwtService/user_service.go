package jwtservice

import (
	"GoBank/internal/domain"
	"GoBank/internal/repository"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type UserService struct {
	repo repository.UserRepository
}

func NewJwtService(repo repository.UserRepository) *UserService {
	var userService UserService = UserService{
		repo: repo,
	}
	return &userService
}

func (j *UserService) Login(ctx context.Context, email string, password string) error {
	_, err := j.GetUserDetails(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("user not found")
		}
		return err
	}
	//uенерация токена
	return nil
}

func (j *UserService) GetUserDetails(ctx context.Context, email string) (domain.User, error) {
	user, err := j.repo.CheckUserInDataBase(ctx, email)
	if err != nil {
		return user, err
	}
	return user, nil
}
