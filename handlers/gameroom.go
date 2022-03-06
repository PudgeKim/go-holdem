package handlers

import (
	"github.com/PudgeKim/go-holdem/gamerooms"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type LeavePlayerReq struct {
	RoomId     string `json:"room_id"`
	PlayerName string `json:"player_name"`
	LeaveOrAdd string `json:"leave_or_add"` // "leave" or "add"
}

const (
	Leave = "leave"
	Add   = "add"
)

func LeavePlayer(c *gin.Context) {
	var req LeavePlayerReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomId, err := uuid.Parse(req.RoomId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	room, err := gamerooms.GetGameRoom(roomId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch req.LeaveOrAdd {
	case Leave:
		err = room.LeavePlayer(req.PlayerName)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	case Add:
		// TODO
		// grpc 통신으로 플레이어 정보 갖고 와야함
		err = room.AddPlayer(req.PlayerName)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(200, gin.H{})
}
