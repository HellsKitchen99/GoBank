package server

import (
	delivery "GoBank/internal/delivery/http"
	"GoBank/internal/usecase"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(userService *usecase.UserService) *Server {
	router := gin.Default()

	userHandler := delivery.NewUserHandler(userService)

	userHandler.RegisterRoutes(router)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	var server Server = Server{
		httpServer: srv,
	}
	return &server
}

func (s *Server) Run() {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			panic(err)
		}
	}()
	fmt.Println("server started on http://localhost:8080")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown - %v\n", err)
	}
	fmt.Println("server exited cleanly")
}
