package repository

import (
	"GoBank/internal/domain"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	var userRepo UserRepo = UserRepo{
		db: db,
	}
	return &userRepo
}

func (r *UserRepo) CheckUserInDataBase(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	query := `SELECT id, name, email, password, roles, amount FROM users WHERE email=$1`
	row := r.db.QueryRow(ctx, query, email)
	if err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Roles); err != nil {
		return user, err
	}
	return user, nil
}

func (r *UserRepo) CheckUserInDataBaseById(ctx context.Context, id int64) bool {
	var user domain.User
	query := `SELECT id, name, email, password, roles, amount FROM users WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	if err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Roles, &user.Amount); err != nil {
		return false
	}
	return true
}

func (r *UserRepo) CreateUser(ctx context.Context, name, email, password string, roles []string) (int64, error) {
	_, err := r.CheckUserInDataBase(ctx, email)
	if err == nil {
		return -1, fmt.Errorf("user already exists")
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return -1, err
	}
	query := `INSERT INTO users (name, email, password, roles) VALUES ($1, $2, $3, $4) RETURNING id`

	tag := r.db.QueryRow(ctx, query, name, email, password, roles)
	var id int64
	if err := tag.Scan(&id); err != nil {
		return -1, err
	}
	return id, nil
}

// TRANSACTIONS
type TransactionsRepo struct {
	db *pgxpool.Pool
}

func NewTransactionsRepo(db *pgxpool.Pool) *TransactionsRepo {
	var transactionsRepo TransactionsRepo = TransactionsRepo{
		db: db,
	}
	return &transactionsRepo
}

func (t *TransactionsRepo) CreateTransaction(ctx context.Context, from, to int64, amount float64, timeOfCreation time.Time, status string) error {
	query := `INSERT INTO transactions (from_user, to_user, amount, created_at, status) VALUES ($1, $2, $3, $4, $5)`
	tag, err := t.db.Exec(ctx, query, from, to, amount, timeOfCreation, status)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("rows affected is not 1")
	}
	return nil
}

func (t *TransactionsRepo) GetAmountOfSenderFromDataBase(ctx context.Context, from int64) (float64, error) {
	query := `SELECT amount FROM users WHERE id = $1`
	var amount float64
	row := t.db.QueryRow(ctx, query, from)
	if err := row.Scan(&amount); err != nil {
		return -1, err
	}
	return amount, nil
}

func (t *TransactionsRepo) CheckUserToSendInDataBase(ctx context.Context, to int64) error {
	var user domain.User
	query := `SELECT id, name, email, password, roles, amount FROM users WHERE id = $1`
	row := t.db.QueryRow(ctx, query, to)
	if err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Roles, &user.Amount); err != nil {
		return err
	}
	return nil
}

func (t *TransactionsRepo) BeginTransaction(ctx context.Context) (pgx.Tx, error) {
	return t.db.Begin(ctx)
}

func (t *TransactionsRepo) MinusMoneyTx(ctx context.Context, tx pgx.Tx, id int64, amount float64) error {
	query := `UPDATE users SET amount = amount - $1 WHERE id = $2 AND amount >= $1`
	tag, err := tx.Exec(ctx, query, amount, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("not enough funds or user not found")
	}
	return nil
}

func (t *TransactionsRepo) AddMoneyTx(ctx context.Context, tx pgx.Tx, id int64, amount float64) error {
	query := `UPDATE users SET amount = amount + $1 WHERE id = $2`
	tag, err := tx.Exec(ctx, query, amount, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("rows affected is not 1")
	}
	return nil
}
