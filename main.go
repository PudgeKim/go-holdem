package main

import (
	pb "github.com/PudgeKim/go-holdem-protos/protos"
	"github.com/PudgeKim/go-holdem/gamerooms"
	"github.com/PudgeKim/go-holdem/grpc_client"
	"github.com/PudgeKim/go-holdem/handlers"
	"google.golang.org/grpc"
)

const (
	GrpcAddress   = "localhost:6060"
	ServerAddress = "localhost:7070"
)

func main() {
	conn, err := grpc.Dial(GrpcAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic("client did not connect: ")
	}
	defer conn.Close()

	grpcClient := pb.NewAuthClient(conn)
	gameroomsMap := make(gamerooms.GameRooms)

	grpcHandler := grpc_client.NewGrpcHandler(grpcClient)
	gameroomsHandler := handlers.NewGameRoomsHandler(gameroomsMap)
	gameroomHandler := handlers.NewGameRoomHandler(gameroomsMap, grpcHandler)
	gameHandler := handlers.NewGameHandler(gameroomsMap, grpcHandler)

	handlers := handlers.NewHandlers(gameroomsHandler, gameroomHandler, gameHandler)

	router := handlers.Routes()
	router.Run(ServerAddress)

}
