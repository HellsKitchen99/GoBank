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
		if errors.Is(err, pgx.ErrNoRows) {
			return user, fmt.Errorf("user not found")
		}
		return user, err
	}
	return user, nil
}
