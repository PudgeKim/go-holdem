package handlers

import (
	"github.com/PudgeKim/go-holdem/gamerooms"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GameRoomsHandler struct {
	rooms gamerooms.GameRooms
}

func NewGameRoomsHandler(rooms gamerooms.GameRooms) *GameRoomsHandler {
	return &GameRoomsHandler{
		rooms: rooms,
	}
}

type AddRoomReq struct {
	RoomName string `json:"room_name"`
	HostName string `json:"host_name,omitempty"`
}

func (g GameRoomsHandler) AddRoom(c *gin.Context) {
	var req AddRoomReq

	if err := c.ShouldBindJSON(&req); err != nil {
		badRequestWithError(c, err)
		return
	}

	room, err := g.rooms.AddRoom(req.RoomName, req.HostName)
	if err != nil {
		badRequestWithError(c, err)
		return
	}

	// 베팅 처리 후 프론트에 다시 응답하기 위한 채널을 설정
	// 각 게임마다 독립적인 채널을 가지고 있게됨
	betResponseChan := SetBetResponseChan(room.Id)
	room.Game.SetBetResponseChan(betResponseChan)

	values := make(map[string]interface{})
	values["room_id"] = room.Id
	values["room_name"] = room.Name
	values["hostname"] = room.HostName
	values["room_limit"] = room.Limit

	statusOkWithSuccess(c, values, nil)
}

type RemoveRoomReq struct {
	RoomId string `json:"room_id"`
}

func (g GameRoomsHandler) RemoveRoom(c *gin.Context) {
	var req RemoveRoomReq

	if err := c.ShouldBindJSON(&req); err != nil {
		badRequestWithError(c, err)
		return
	}

	roomId, err := uuid.Parse(req.RoomId)
	if err != nil {
		badRequestWithError(c, err)
		return
	}

	g.rooms.RemoveRoom(roomId)

	statusOkWithSuccess(c, nil, nil)
}
