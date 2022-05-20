package main

import (
	"net/http"

	"github.com/PudgeKim/go-holdem/cacheserver"
	"github.com/PudgeKim/go-holdem/db"
	"github.com/PudgeKim/go-holdem/handler"
	"github.com/PudgeKim/go-holdem/persistence"
	"github.com/PudgeKim/go-holdem/service"
	"github.com/gorilla/websocket"
)

const (
	ServerAddress = "localhost:7070"
)

var createTableSchema = `
CREATE TABLE IF NOT EXISTS users (
	id serial PRIMARY KEY,
    nickname text,
    email text,
    balance bigint
);
`

var dropTableSchema = `
DROP TABLE IF EXISTS users
`

func main() {
	db, err := db.NewPostgresDB(db.DBConfig)
	if err != nil {
		panic(err)
	}

	db.MustExec(createTableSchema)
	db.MustExec("INSERT INTO users (nickname, email, balance) VALUES ($1, $2, $3)", "john", "john@gmail.com", 10000)
	db.MustExec("INSERT INTO users (nickname, email, balance) VALUES ($1, $2, $3)", "sarah", "sarah@gmail.com", 10000)

	redisClient := cacheserver.NewRedis()


	chatRepo := persistence.NewChatRepository(redisClient)
	gameRepo := persistence.NewGameRepository(redisClient)
	userRepo := persistence.NewUserRepository(db)

	chatService := service.NewChatService(chatRepo)
	gameService := service.NewGameService(userRepo, gameRepo)

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	gameHandler := handler.NewGameHandler(&upgrader, chatService, gameService)

	myHandlers := handler.NewHandlers(gameHandler)

	router := myHandlers.Routes()
	router.Run(ServerAddress)
}
