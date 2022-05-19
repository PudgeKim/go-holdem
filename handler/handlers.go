package handler

import "github.com/gin-gonic/gin"

type Handlers struct {
	gameHandler  *GameHandler
}

func NewHandlers(gameHandler *GameHandler) *Handlers {
	return &Handlers{
		gameHandler:  gameHandler,
	}
}

func (h *Handlers) Routes() *gin.Engine {
	router := gin.Default()

	router.GET("/game/joinroom/:roomid", h.gameHandler.JoinRoom)

	return router
}
