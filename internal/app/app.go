package app

import (
	"GoBank/internal/config"
	"GoBank/internal/middleware"
	"GoBank/internal/repository"
	"GoBank/internal/server"
	"GoBank/internal/service"
	"GoBank/internal/usecase"
	"context"
	"fmt"
	"net/url"
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
	password := url.QueryEscape(cdb.Password)
	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", cdb.User, password, cdb.Host, cdb.Port, cdb.DbName)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return err
	}
	defer pool.Close()

	userRepo := repository.NewUserRepo(pool)
	jwtService := service.NewJwtService(cjwt.Key, cjwt.Duration)
	userService := usecase.NewUserService(userRepo, jwtService)
	authMiddleware := middleware.NewAuthMiddleware(jwtService, userRepo)
	transactionRepo := repository.NewTransactionsRepo(pool)
	transactionService := usecase.NewTransactionService(transactionRepo)
	srv := server.NewServer(userService, authMiddleware, transactionService)
	srv.Run()
	return nil
}
