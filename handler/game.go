package handler

import (
	"fmt"
	"net/http"

	"github.com/PudgeKim/go-holdem/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type GameHandler struct {
	upgrader *websocket.Upgrader
	chatService service.ChatService
}

func NewGameHandler(upgrader *websocket.Upgrader, chatService service.ChatService) *GameHandler {
	return &GameHandler{
		upgrader: upgrader,
		chatService: chatService,
	}
}

type JoinRoomReq struct {
	RoomId string `uri:"roomid" binding:"required"`
}

func (g *GameHandler) JoinRoom(c *gin.Context) {
	var joinRoomReq JoinRoomReq

	if err := c.ShouldBindUri(&joinRoomReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "VALIDATEERR",
		})
		return 
	}

	ws, err := g.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("JoinRoomErr: ", err.Error())
		return 
	}
	defer ws.Close()

	chatChan := make(chan string)

	if err := g.chatService.Subscribe(c, joinRoomReq.RoomId, chatChan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return 
	}
	
	go func ()  {
		for {
			chatMsg := <-chatChan
			fmt.Println("chatMsg: ", chatMsg)
			if err := ws.WriteMessage(1, []byte(chatMsg)); err != nil {
				panic(fmt.Sprintf("goroutine WriteMessage Err: %s", err.Error()))
			}
		}
	}()

	for {
		mt, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("JoinRoomReadMessageErr: ", err.Error())
			break 
		}

		fmt.Println(mt, msg)


	}
}