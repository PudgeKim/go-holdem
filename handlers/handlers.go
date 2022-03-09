package handlers

import "github.com/gin-gonic/gin"

type Handlers struct {
	roomsHandler GameRoomsHandler
	roomHandler  GameRoomHandler
	gameHandler  GameHandler
}

func NewHandlers(gameroomsHandler GameRoomsHandler, gameroomHandler GameRoomHandler, gameHandler GameHandler) *Handlers {
	return &Handlers{
		roomsHandler: gameroomsHandler,
		roomHandler:  gameroomHandler,
		gameHandler:  gameHandler,
	}
}

func (h *Handlers) Routes() *gin.Engine {
	router := gin.Default()

	router.POST("/gamerooms/add", h.roomsHandler.AddRoom)
	router.POST("/gamerooms/remove", h.roomsHandler.RemoveRoom)

	router.POST("/gameroom/add", h.roomHandler.AddPlayer)
	router.POST("/gameroom/leave", h.roomHandler.LeavePlayer)

	return router
}
