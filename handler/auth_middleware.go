package handler

import (
	"net/http"

	"github.com/PudgeKim/go-holdem/service"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authService *service.AuthService
}

func NewAuthMiddleware(authService *service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

func (a *AuthMiddleware) ValidateToken(c *gin.Context) {
	token, err := c.Cookie("access_token"); if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	userId, err := a.authService.ValidateToken(token); if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	c.Set("userId", userId)
	c.Next()
}