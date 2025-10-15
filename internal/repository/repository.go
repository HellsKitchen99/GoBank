package repository

import (
	"GoBank/internal/domain"
	"context"
)

type UserRepository interface {
	CheckUserInDataBase(ctx context.Context, email string) (domain.User, error)
	CreateUser(ctx context.Context, email, password string) (domain.User, error)
}
