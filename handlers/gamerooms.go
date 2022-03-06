package handlers

import (
	"github.com/PudgeKim/go-holdem/gamerooms"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type AddRoomReq struct {
	RoomName string `json:"room_name"`
	HostName string `json:"host_name,omitempty"`
}

func AddRoom(c *gin.Context) {
	var req AddRoomReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := gamerooms.AddRoom(req.RoomName, req.HostName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

type RemoveRoomReq struct {
	RoomId string `json:"room_id"`
}

func RemoveRoom(c *gin.Context) {
	var req RemoveRoomReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomId, err := uuid.Parse(req.RoomId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	gamerooms.RemoveRoom(roomId)
	c.JSON(http.StatusOK, gin.H{})
}
