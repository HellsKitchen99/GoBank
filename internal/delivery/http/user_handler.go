package http

import (
	"GoBank/internal/usecase"

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

func (h *UserHandler) Register(ctx *gin.Context) {
	//заглушка
}

func (h *UserHandler) Login(ctx *gin.Context) {
	//заглушка
}
