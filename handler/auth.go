package handler

import (
	"net/http"

	"github.com/PudgeKim/go-holdem/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type SignUpReq struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,len=8"`
	Nickname string `json:"nickname" binding:"required,len=2"`
}

func (a *AuthHandler) SignUp(c *gin.Context) {
	var signUpReq SignUpReq

	if err := c.ShouldBindJSON(&signUpReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return 
	}

	if err := a.authService.SignUp(c, signUpReq.Email, signUpReq.Password, signUpReq.Nickname); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return 
	}

	c.JSON(http.StatusCreated, gin.H{
		"nickname": signUpReq.Nickname,
	})
}

type SignInReq struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,len=8"`
}

func (a *AuthHandler) SignIn(c *gin.Context) {
	var signInReq SignInReq
	
	if err := c.ShouldBindJSON(&signInReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return 
	}

	token, err := a.authService.SignIn(c, signInReq.Email, signInReq.Password); if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return 
	}

	c.SetCookie("access_token", token, 3600, "/", "localhost", false, false)
	c.Status(http.StatusOK)
}