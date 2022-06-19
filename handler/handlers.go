package handler

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

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

	router.Use(cors.New(cors.Config{
        AllowOrigins: []string{"http://127.0.0.1:5500"},
        AllowMethods: []string{"POST", "PUT", "PATCH", "DELETE"},
        AllowHeaders: []string{"Content-Type,access-control-allow-origin, access-control-allow-headers"},
		AllowCredentials: true,
    }))

	router.POST("/auth/signup", h.authHandler.SignUp)
	router.POST("/auth/signin", h.authHandler.SignIn)

	router.GET("/game/joinroom/:roomid", h.authMiddleware.ValidateToken, h.gameHandler.JoinRoom)
	router.POST("/game", h.authMiddleware.ValidateToken, h.gameHandler.CreateGameRoom)
	
	router.GET("/check", func(c *gin.Context) {
		cookie, err := c.Cookie("access_token")
		if err != nil {
			fmt.Println("err: ", err.Error())
		}
		fmt.Println("cookie: ", cookie)

	})
	
	return router
}
