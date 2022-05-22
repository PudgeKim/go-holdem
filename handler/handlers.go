package handler

import "github.com/gin-gonic/gin"

type Handlers struct {
	gameHandler  *GameHandler
	authHandler *AuthHandler
	authMiddleware *AuthMiddleware
}

func NewHandlers(gameHandler *GameHandler, authHandler *AuthHandler, authMiddleware *AuthMiddleware) *Handlers {
	return &Handlers{
		gameHandler:  gameHandler,
		authHandler: authHandler,
		authMiddleware: authMiddleware,
	}
}

func (h *Handlers) Routes() *gin.Engine {
	router := gin.Default()

	router.POST("/auth/signup", h.authHandler.SignUp)
	router.POST("/auth/signin", h.authHandler.SignIn)

	router.GET("/game/joinroom/:roomid", h.authMiddleware.ValidateToken, h.gameHandler.JoinRoom)
	router.POST("/game", h.authMiddleware.ValidateToken, h.gameHandler.CreateGameRoom)
	

	return router
}
