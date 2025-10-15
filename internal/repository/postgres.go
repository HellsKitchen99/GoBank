package repository

import (
	"GoBank/internal/domain"
	"context"
	"errors"
	"fmt"

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
	query := `SELECT id, name, email, password, roles FROM users WHERE email=$1`
	row := r.db.QueryRow(ctx, query, email)
	if err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Roles); err != nil {
		return user, err
	}
	return user, nil
}

func (r *UserRepo) CreateUser(ctx context.Context, name, email, password string, roles []string) error {
	_, err := r.CheckUserInDataBase(ctx, email)
	if err == nil {
		return fmt.Errorf("user already exists")
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	query := `INSERT INTO users (name, email, password, roles) VALUES ($1, $2, $3, $4)`
	tag, err := r.db.Exec(ctx, query, name, email, password, roles)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("bad insert")
	}
	return nil
}
