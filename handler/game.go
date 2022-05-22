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
	authService *service.AuthService
}

func NewGameHandler(upgrader *websocket.Upgrader, chatService *service.ChatService, gameService *service.GameService, authService *service.AuthService) *GameHandler {
	return &GameHandler{
		upgrader: upgrader,
		chatService: chatService,
		gameService: gameService,
		authService: authService,
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreateGameReq struct {
	UserId int64 `json:"user_id" binding:"required"`
	GameBalance uint64 `json:"game_balance" binding:"required"`
}

func (g *GameHandler) CreateGameRoom(c *gin.Context) {
	var createGameReq CreateGameReq

	if err := c.ShouldBindJSON(&createGameReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	id, _ := c.Get("userId")
	userId, _ := id.(int64)

	if userId != createGameReq.UserId {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user id in access token doesn't match to request's user id",
		})
		return 
	}

	user, err := g.authService.FindUser(c, createGameReq.UserId); if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return 
	}


	game, err := g.gameService.CreateGame(c, user, createGameReq.GameBalance); if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"room_id": game.RoomId.String(),
		"hostname": game.HostName,
	})
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
	IsReady bool `json:"is_ready"`
}

// room에 들어가는 순간 websocket을 통해
// ready, betting 등을 할 수 있음 
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
			}
		case "start":
			res, err := g.gameService.StartGame(c, gameReq.RoomId, gameReq.Nickname); if err != nil {
				errorResponse := ErrorResponse{Error: err.Error()}
				if err := ws.WriteJSON(errorResponse); err != nil {
					fmt.Println("GameStartWriteJsonErr1: ", err.Error())
				}
			}
			
			if err := ws.WriteJSON(res); err != nil {
				fmt.Println("GameStartWriteJsonErr2: ", err.Error())
			}
		case "ready":
			if err := g.gameService.HandleReady(c, gameReq.RoomId, gameReq.Nickname, gameReq.IsReady); err != nil {
				fmt.Println("Ready: ", err.Error())
			}
		}
		

	}
}