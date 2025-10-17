package usecase

import (
	"GoBank/internal/domain"
	"GoBank/internal/repository"
	"GoBank/internal/service"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo       repository.UserRepository
	jwtService *service.JwtService
}

func NewUserService(repo repository.UserRepository, jwtService *service.JwtService) *UserService {
	var userService UserService = UserService{
		repo:       repo,
		jwtService: jwtService,
	}
	return &userService
}

func (j *UserService) Register(ctx context.Context, name, email, password string) (string, error) {
	_, err := j.repo.CheckUserInDataBase(ctx, email)
	if err == nil {
		return "", ErrUserAlreadyExists
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return "", err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	roles := []string{string(domain.RoleUser)}
	id, err := j.repo.CreateUser(ctx, name, email, string(hashedPassword), roles)
	if err != nil {
		return "", err
	}
	token, err := j.jwtService.GenerateJwtToken(id, roles)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (j *UserService) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := j.GetUserDetails(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrUserNotFound
		}
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", ErrWrongPassword
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

func (j *UserService) GetInfo(from int64) (domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	user, err := j.repo.GetInfoFromDataBase(ctx, from)
	if err != nil {
		return user, nil
	}
	return user, nil
}

func (j *UserService) MakeDeposit(deposit domain.Deposit, from int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := j.repo.AddDepositToDataBase(ctx, from, deposit.Amount); err != nil {
		return err
	}
	return nil
}
