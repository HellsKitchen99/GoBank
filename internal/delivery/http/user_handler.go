package http

import (
	"GoBank/internal/domain"
	"GoBank/internal/usecase"
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *usecase.UserService
}

func NewUserHandler(service *usecase.UserService) *UserHandler {
	var userHandler UserHandler = UserHandler{
		service: service,
	}
	return &userHandler
}

func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)
}

func (h *UserHandler) Register(c *gin.Context) {
	var userFromFront domain.UserFromFront

	if err := c.ShouldBindJSON(&userFromFront); err != nil {
		c.JSON(400, gin.H{
			"error": "bad json",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := h.service.Register(ctx, userFromFront.Name, userFromFront.Email, userFromFront.Password)

	if errors.Is(err, usecase.ErrUserAlreadyExists) {
		c.JSON(409, gin.H{
			"error": err.Error(),
		})
		return
	} else if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(201, gin.H{
		"answer": "user has been created succesfully",
	})
}

func (h *UserHandler) Login(ctx *gin.Context) {
	//заглушка
}
