package usecase

import (
	"GoBank/internal/domain"
	"GoBank/internal/repository"
	"GoBank/internal/service"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type UserService struct {
	repo       repository.UserRepository
	jwtService *service.JwtService
}

func NewUserService(repo repository.UserRepository) *UserService {
	var userService UserService = UserService{
		repo: repo,
	}
	return &userService
}

func Register(ctx context.Context, name, email, password string) {

}

func (j *UserService) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := j.GetUserDetails(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("user not found")
		}
		return "", err
	}
	token, err := j.jwtService.GenerateJwtToken(user.Id, user.Roles)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (j *UserService) GetUserDetails(ctx context.Context, email string) (domain.User, error) {
	user, err := j.repo.CheckUserInDataBase(ctx, email)
	if err != nil {
		return user, err
	}
	return user, nil
}
