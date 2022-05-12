package handlers

import (
	"context"
	"net/http"

	"github.com/PudgeKim/go-holdem/gamerooms"
	"github.com/gin-gonic/gin"
)

const (
	Leave = "leave"
	Add   = "add"
)

type GameRoomHandler struct {
	rooms       gamerooms.GameRooms
}

func NewGameRoomHandler(rooms gamerooms.GameRooms) *GameRoomHandler {
	return &GameRoomHandler{
		rooms:       rooms,

	}
}

type AddPlayerReq struct {
	RoomId       string `json:"room_id"`
	PlayerId     string `json:"player_id"`
	PlayerName   string `json:"player_name"`
	TotalBalance uint64 `json:"total_balance"`
	GameBalance  uint64 `json:"game_balance"`
}

func (g GameRoomHandler) AddPlayer(c *gin.Context) {
	var req AddPlayerReq

	if err := c.ShouldBindJSON(&req); err != nil {
		badRequestWithError(c, err)
		return
	}

	user, err := g.grpcHandler.GetUser(context.Background(), req.PlayerId)
	if err != nil {
		badRequestWithError(c, err)
		return
	}
	if user.Id != req.PlayerId {
		c.JSON(http.StatusBadRequest, gin.H{"error": grpc_error.InvalidUser})
		return
	}

	room, err := g.rooms.GetGameRoomAfterParse(req.RoomId)
	if err != nil {
		badRequestWithError(c, err)
		return
	}

	err = room.AddPlayer(req.PlayerName, req.TotalBalance, req.GameBalance)
	if err != nil {
		badRequestWithError(c, err)
		return
	}

	statusOkWithSuccess(c, nil, nil)
}

type LeavePlayerReq struct {
	RoomId     string `json:"room_id"`
	PlayerId   string `json:"player_id"`
	PlayerName string `json:"player_name"`
}

func (g GameRoomHandler) LeavePlayer(c *gin.Context) {
	var req LeavePlayerReq

	if err := c.ShouldBindJSON(&req); err != nil {
		badRequestWithError(c, err)
		return
	}

	user, err := g.grpcHandler.GetUser(context.Background(), req.PlayerId)
	if err != nil {
		badRequestWithError(c, err)
		return
	}
	if user.Id != req.PlayerId {
		c.JSON(http.StatusBadRequest, gin.H{"error": grpc_error.InvalidUser})
		return
	}

	room, err := g.rooms.GetGameRoomAfterParse(req.RoomId)
	if err != nil {
		badRequestWithError(c, err)
		return
	}

	err = room.LeavePlayer(req.PlayerName)
	if err != nil {
		badRequestWithError(c, err)
		return
	}

	statusOkWithSuccess(c, nil, nil)
}
