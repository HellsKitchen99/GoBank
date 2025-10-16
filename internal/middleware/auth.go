package middleware

import (
	"GoBank/internal/service"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwtService *service.JwtService
}

func NewAuthMiddleware(jwtService *service.JwtService) *AuthMiddleware {
	var authMiddleware AuthMiddleware = AuthMiddleware{
		jwtService: jwtService,
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
		if valid := a.jwtService.TokenValidation(token); !valid {
			c.JSON(401, gin.H{
				"error": "unauthorized",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
