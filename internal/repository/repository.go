package repository

import (
	"GoBank/internal/domain"
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	CheckUserInDataBase(ctx context.Context, email string) (domain.User, error)
	CreateUser(ctx context.Context, name, email, password string, roles []string) (int64, error)
	CheckUserInDataBaseById(ctx context.Context, id int64) bool
	GetInfoFromDataBase(ctx context.Context, id int64) (domain.User, error)
	AddDepositToDataBase(ctx context.Context, id int64, amount float64) error
}

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, from, to int64, amount float64, timeOfCreation time.Time, status string) error
	GetAmountOfSenderFromDataBase(ctx context.Context, from int64) (float64, error)
	CheckUserToSendInDataBase(ctx context.Context, to int64) error
	BeginTransaction(ctx context.Context) (pgx.Tx, error)
	MinusMoneyTx(ctx context.Context, tx pgx.Tx, id int64, amount float64) error
	AddMoneyTx(ctx context.Context, tx pgx.Tx, id int64, amount float64) error
}
