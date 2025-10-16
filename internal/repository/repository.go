package repository

import (
	"GoBank/internal/domain"
	"context"
	"time"
)

type UserRepository interface {
	CheckUserInDataBase(ctx context.Context, email string) (domain.User, error)
	CreateUser(ctx context.Context, name, email, password string, roles []string) (int64, error)
	CheckUserInDataBaseById(ctx context.Context, id int64) bool
}

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, from, to int64, amount float64, timeOfCreation time.Time, status string) error
	GetAmountOfSenderFromDataBase(ctx context.Context, from int64) (float64, error)
	CheckUserToSendInDataBase() error
	MinusMoney(amount float64) error
	AddMoney(amount float64) error
}
