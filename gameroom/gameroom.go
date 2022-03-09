package gameroom

import (
	"github.com/PudgeKim/go-holdem/game"
	"github.com/PudgeKim/go-holdem/gameerror"
	"github.com/PudgeKim/go-holdem/player"
	"github.com/google/uuid"
)

const RoomLimit = 7

type GameRoom struct {
	Id       uuid.UUID  `json:"id"`
	Name     string     `json:"name"`      // 방 이름
	HostName string     `json:"host_name"` // 방장 닉네임
	Limit    uint       `json:"limit"`     // 방 하나에 최대 플레이어 수
	Game     *game.Game `json:"-"`
}

func NewGameRoom(name, hostName string) (*GameRoom, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &GameRoom{
		Id:       id,
		Name:     name,
		HostName: hostName,
		Limit:    RoomLimit,
		Game:     game.New(),
	}, nil
}

func (g *GameRoom) LeavePlayer(nickname string) error {
	p, err := g.FindPlayer(nickname)
	if err != nil {
		return err
	}

	if !g.Game.IsStarted {
		if err = g.removePlayer(*p); err != nil {
			return err
		}
		return nil
	}

	// 이미 게임 중이라면 bool 값만 바꿔두고
	// 게임이 종료되면 IsLeft가 true인 플레이어들 처리함
	p.IsLeft = true
	return nil
}

func (g *GameRoom) HandleReady(nickname string, isReady bool) error {
	p, err := g.FindPlayer(nickname)
	if err != nil {
		return err
	}

	if g.Game.IsStarted {
		return gameerror.GameAlreadyStarted
	}

	p.IsReady = isReady
	return nil
}

func (g *GameRoom) AddPlayer(nickname string, totalBalance, gameBalance uint64) error {
	_, err := g.FindPlayer(nickname)
	// nil이면 플레이어가 이미 존재하는데 또 요청이 온 것
	if err == nil {
		return gameerror.PlayerAlreadyExists
	}

	if len(g.Game.Players) > RoomLimit {
		return gameerror.PlayerLimitationError
	}

	p := player.New(nickname, totalBalance, gameBalance)
	g.Game.Players = append(g.Game.Players, &p)
	return nil
}

func (g *GameRoom) removePlayer(p player.Player) error {
	removeIndex := -1
	for i := 0; i < len(g.Game.Players); i++ {
		if p.Nickname == g.Game.Players[i].Nickname {
			removeIndex = i
			break
		}
	}

	if removeIndex == -1 {
		return gameerror.NoPlayerExists
	}

	g.Game.Players = append(g.Game.Players[:removeIndex], g.Game.Players[removeIndex+1:]...)

	return nil
}

func (g *GameRoom) FindPlayer(nickname string) (*player.Player, error) {
	for _, p := range g.Game.Players {
		if p.Nickname == nickname {
			return p, nil
		}
	}
	return nil, gameerror.NoPlayerExists
}
