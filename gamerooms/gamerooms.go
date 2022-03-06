package gamerooms

import (
	"github.com/PudgeKim/go-holdem/gameerror"
	"github.com/PudgeKim/go-holdem/gameroom"
	"github.com/google/uuid"
)

type GameRooms map[uuid.UUID]*gameroom.GameRoom

var Handler = make(GameRooms)

func AddRoom(roomName, hostName string) error {
	room, err := gameroom.NewGameRoom(roomName, hostName)
	if err != nil {
		return err
	}

	Handler[room.Id] = room
	return nil
}

func RemoveRoom(roomId uuid.UUID) {
	delete(Handler, roomId)
}

func GetGameRoom(roomId uuid.UUID) (*gameroom.GameRoom, error) {
	room := Handler[roomId]
	if room.Name == "" {
		return nil, gameerror.NoRoomExists
	}

	return room, nil

}
