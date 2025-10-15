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

func (r *UserRepo) CreateUser(ctx context.Context, name, email, password string, roles []string) (int64, error) {
	_, err := r.CheckUserInDataBase(ctx, email)
	if err == nil {
		return 0, fmt.Errorf("user already exists")
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return 0, err
	}
	query := `INSERT INTO users (name, email, password, roles) VALUES ($1, $2, $3, $4) RETURNING id`
	tag := r.db.QueryRow(ctx, query, name, email, password, roles)
	var id int64
	if err := tag.Scan(&id); err != nil {
		return id, err
	}
	return id, nil
}
