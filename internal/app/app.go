package app

import (
	"GoBank/internal/config"
	"GoBank/internal/repository"
	"GoBank/internal/server"
	"GoBank/internal/service"
	"GoBank/internal/usecase"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	cdb := cfg.DB
	cjwt := cfg.Jwt
	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", cdb.User, cdb.Password, cdb.Host, cdb.Port, cdb.DbName)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return err
	}
	defer pool.Close()

	jwtService := service.NewJwtService(cjwt.Key, cjwt.Duration)
	repo := repository.NewUserRepo(pool)
	userService := usecase.NewUserService(repo, jwtService)
	srv := server.NewServer(userService)
	srv.Run()
	return nil
}
