package main

import (
	"net/http"

	"github.com/PudgeKim/go-holdem/cacheserver"
	"github.com/PudgeKim/go-holdem/handler"
	"github.com/PudgeKim/go-holdem/persistence"
	"github.com/PudgeKim/go-holdem/service"
	"github.com/gorilla/websocket"
)

const (
	ServerAddress = "localhost:7070"
)

func main() {
	redisClient := cacheserver.NewRedis()


	chatRepo := persistence.NewChatRepository(redisClient)


	chatService := service.NewChatService(chatRepo)

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	gameHandler := handler.NewGameHandler(&upgrader, *chatService)

	myHandlers := handler.NewHandlers(gameHandler)

	router := myHandlers.Routes()
	router.Run(ServerAddress)
}
