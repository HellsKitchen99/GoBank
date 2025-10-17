package http

import (
	"GoBank/internal/domain"
	"GoBank/internal/middleware"
	"GoBank/internal/usecase"
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service            *usecase.UserService
	authMiddleware     *middleware.AuthMiddleware
	transactionService *usecase.TransactionService
}

func NewUserHandler(service *usecase.UserService, authMiddleware *middleware.AuthMiddleware, transactionService *usecase.TransactionService) *UserHandler {
	var userHandler UserHandler = UserHandler{
		service:            service,
		authMiddleware:     authMiddleware,
		transactionService: transactionService,
	}
	return &userHandler
}

func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)

	userGroup := r.Group("/user")
	userGroup.Use(h.authMiddleware.Filter())
	{
		userGroup.POST("/transaction", h.Transaction)
	}

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
	token, err := h.service.Register(ctx, userFromFront.Name, userFromFront.Email, userFromFront.Password)

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
		"answer": token,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var userFromFront domain.UserFromFront
	if err := c.ShouldBindJSON(&userFromFront); err != nil {
		c.JSON(400, gin.H{
			"error": "bad json",
		})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	token, err := h.service.Login(ctx, userFromFront.Email, userFromFront.Password)
	if errors.Is(err, usecase.ErrUserNotFound) {
		c.JSON(404, gin.H{
			"error": "user not registered",
		})
		return
	} else if errors.Is(err, usecase.ErrWrongPassword) {
		c.JSON(400, gin.H{
			"error": "wrong password",
		})
		return
	}
	c.JSON(200, gin.H{
		"answer": token,
	})
}

func (h *UserHandler) Transaction(c *gin.Context) {
	var transaction domain.TransactionFromFront
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(400, gin.H{
			"error": "bad json",
		})
		return
	}
	from := c.GetInt64("user_id")
	err := h.transactionService.ValidateTransaction(transaction, from)
	if err != nil {
		if errors.Is(err, usecase.ErrAmountToTransfer) {
			c.JSON(400, gin.H{
				"error": "not enough balance to make transfer",
			})
			return
		} else if errors.Is(err, usecase.ErrSameUser) {
			c.JSON(400, gin.H{
				"error": "trying to transfer the same user",
			})
			return
		}
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}
	if err := h.transactionService.CreateTransaction(transaction, from); err != nil {
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}
	if err := h.transactionService.MakeTransaction(from, transaction.To, transaction.Amount); err != nil {
		c.JSON(500, gin.H{
			"error": "transaction error",
		})
		return
	}
	c.JSON(200, gin.H{
		"answer": "transaction success",
	})
}
