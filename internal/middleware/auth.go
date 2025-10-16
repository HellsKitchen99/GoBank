package middleware

import (
	"GoBank/internal/repository"
	"GoBank/internal/service"
	"context"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwtService *service.JwtService
	userRepo   repository.UserRepository
}

func NewAuthMiddleware(jwtService *service.JwtService, userRepo repository.UserRepository) *AuthMiddleware {
	var authMiddleware AuthMiddleware = AuthMiddleware{
		jwtService: jwtService,
		userRepo:   userRepo,
	}
	return &authMiddleware
}

func (a *AuthMiddleware) Filter() gin.HandlerFunc {
	return func(c *gin.Context) {
		fullToken := c.GetHeader("Authorization")
		if fullToken == "" {
			c.JSON(401, gin.H{
				"error": "unauthorized",
			})
			c.Abort()
			return
		}
		if !strings.HasPrefix(fullToken, "Bearer ") {
			c.JSON(401, gin.H{
				"error": "unauthorized",
			})
			c.Abort()
			return
		}
		token, _ := strings.CutPrefix(fullToken, "Bearer ")
		id, valid := a.jwtService.TokenValidation(token)
		if !valid {
			c.JSON(401, gin.H{
				"error": "unauthorized",
			})
			c.Abort()
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if have := a.userRepo.CheckUserInDataBaseById(ctx, id); !have {
			c.JSON(404, gin.H{
				"error": "user not registered",
			})
			c.Abort()
			return
		}
		c.Set("user_id", id)
		c.Next()
	}
}
