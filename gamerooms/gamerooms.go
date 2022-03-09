package gamerooms

import (
	"github.com/PudgeKim/go-holdem/gameerror"
	"github.com/PudgeKim/go-holdem/gameroom"
	"github.com/google/uuid"
)

type GameRooms map[uuid.UUID]*gameroom.GameRoom

// var Handler = make(GameRooms)

func (g GameRooms) AddRoom(roomName, hostName string) (*gameroom.GameRoom, error) {
	room, err := gameroom.NewGameRoom(roomName, hostName)
	if err != nil {
		return nil, err
	}

	g[room.Id] = room
	return room, nil
}

func (g GameRooms) RemoveRoom(roomId uuid.UUID) {
	delete(g, roomId)
}

func (g GameRooms) GetGameRoom(roomId uuid.UUID) (*gameroom.GameRoom, error) {
	room := g[roomId]
	if room.Name == "" {
		return nil, gameerror.NoRoomExists
	}

	return room, nil
}

func (g GameRooms) GetGameRoomAfterParse(roomIdString string) (*gameroom.GameRoom, error) {
	roomId, err := uuid.Parse(roomIdString)
	if err != nil {
		return nil, err
	}

	room, err := g.GetGameRoom(roomId)
	if err != nil {
		return nil, err
	}

	return room, nil
}
