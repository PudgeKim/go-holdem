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
	// postgres, err := db.NewPostgresDB(db.DBConfig)
	// if err != nil {
	// 	panic(err)
	// }
	redisClient := cacheserver.NewRedis()


	chatRepo := persistence.NewChatRepository(redisClient)
	// gameRepo := persistence.NewGameRepository(redisClient)
	// userRepo := persistence.NewUserRepository(postgres)

	chatService := service.NewChatService(chatRepo)
	//gameService := service.NewGameService(userRepo, gameRepo)

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	gameHandler := handler.NewGameHandler(&upgrader, chatService)

	myHandlers := handler.NewHandlers(gameHandler)

	router := myHandlers.Routes()
	router.Run(ServerAddress)
}
