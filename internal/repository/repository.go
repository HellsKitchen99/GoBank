package repository

import (
	"GoBank/internal/domain"
	"context"
)

type UserRepository interface {
	CheckUserInDataBase(ctx context.Context, email string) (domain.User, error)
	CreateUser(ctx context.Context, name, email, password string, roles []string) (int64, error)
}
