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
	chatService *service.ChatService
	gameService *service.GameService
}

func NewGameHandler(upgrader *websocket.Upgrader, chatService *service.ChatService, gameService *service.GameService) *GameHandler {
	return &GameHandler{
		upgrader: upgrader,
		chatService: chatService,
		gameService: gameService,
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type JoinRoomReq struct {
	RoomId string `uri:"roomid" binding:"required"`
}

// Type이 Bet이냐 Chat이냐에 따라 
// 요구 필드가 달라짐 
type GameReq struct {
	RoomId string `json:"room_id" binding:"required"`
	Type string `json:"type" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
	Message string `json:"message"`
	BetAmount  uint64 `json:"bet_amount"`
	IsDead     bool `json:"is_dead"`
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
	
	// 다른유저들로부터 채팅이 오면 받아서 전달 
	go func ()  {
		for {
			chatMsg := <-chatChan
			
			if err := ws.WriteMessage(1, []byte(chatMsg)); err != nil {
				panic(fmt.Sprintf("goroutine WriteMessage Err: %s", err.Error())) // panic은 임시용 (나중에 다른걸로 변경)
			}
		}
	}()

	for {
		var gameReq GameReq

		if err := ws.ReadJSON(&gameReq); err != nil {
			fmt.Println("ReadJsonErr: ", err.Error())
			break 
		}

		switch gameReq.Type {
		case "chat":
			err := g.chatService.PublishMessage(c, gameReq.RoomId, gameReq.Nickname, gameReq.Message)
			if err != nil {
				fmt.Println("publishMsgErr: ", err.Error())
				break 
			}
		case "start":
			res, err := g.gameService.StartGame(c, gameReq.RoomId, gameReq.Nickname); if err != nil {
				errorResponse := ErrorResponse{Error: err.Error()}
				if err := ws.WriteJSON(errorResponse); err != nil {
					fmt.Println("GameStartWriteJsonErr1: ", err.Error())
					break 
				}
			}
			
			if err := ws.WriteJSON(res); err != nil {
				fmt.Println("GameStartWriteJsonErr2: ", err.Error())
				break 
			}
		}
		

	}
}