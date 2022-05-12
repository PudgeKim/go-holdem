package main

import (
	"github.com/PudgeKim/go-holdem/gamerooms"
	"github.com/PudgeKim/go-holdem/handlers"
)

const (
	GrpcAddress   = "localhost:6060"
	ServerAddress = "localhost:7070"
)

func main() {

	gameroomsMap := make(gamerooms.GameRooms)

	gameroomsHandler := handlers.NewGameRoomsHandler(gameroomsMap)
	gameroomHandler := handlers.NewGameRoomHandler(gameroomsMap, grpcHandler)
	gameHandler := handlers.NewGameHandler(gameroomsMap, grpcHandler)

	myHandlers := handlers.NewHandlers(gameroomsHandler, gameroomHandler, gameHandler)

	router := myHandlers.Routes()
	router.Run(ServerAddress)
}
