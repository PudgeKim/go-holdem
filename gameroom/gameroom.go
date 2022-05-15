package gameroom

import (
	"context"

	"github.com/PudgeKim/go-holdem/domain/repository"
	"github.com/PudgeKim/go-holdem/game"
	"github.com/PudgeKim/go-holdem/player"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

const RoomLimit = 7

type GameRoom struct {
	ctx context.Context

	userRepo repository.UserRepository

	Id       uuid.UUID  `json:"id"`        // roomId
	Name     string     `json:"name"`      // 방 이름
	HostName string     `json:"host_name"` // 방장 닉네임
	Limit    uint       `json:"limit"`     // 방 하나에 최대 플레이어 수
	Game     *game.Game `json:"-"`
	
}

func NewGameRoom(ctx context.Context, name, hostName string, userRepo repository.UserRepository, redisClient *redis.Client) (*GameRoom, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &GameRoom{
		ctx: ctx,
		userRepo: userRepo,
		Id:       id,
		Name:     name,
		HostName: hostName,
		Limit:    RoomLimit,
		Game:     game.New(ctx, userRepo, redisClient, id),
		
	}, nil
}

func (g *GameRoom) StartGame() (*game.GameStartResponse, error) {
	return g.Game.Start()
}

func (g *GameRoom) Bet(betInfo game.BetInfo) (*game.BetResponse, error) {
	return g.Game.Bet(betInfo)
}

func (g *GameRoom) AddPlayer(player *player.Player) error {
	return g.Game.AddPlayer(player)
}

func (g *GameRoom) FindPlayer(nickname string) (*player.Player, error) {
	return g.Game.FindPlayer(nickname)
}

func (g *GameRoom) LeavePlayer(nickname string) error {
	p, err := g.FindPlayer(nickname)
	if err != nil {
		return err
	}

	if !g.Game.IsStarted {
		if err = g.removePlayer(p.Nickname); err != nil {
			return err
		}
		return nil
	}

	// 이미 게임 중이라면 bool 값만 바꿔두고
	// 게임이 종료되면 IsLeft가 true인 플레이어들 처리함
	p.IsLeft = true
	return nil
}

func (g *GameRoom) removePlayer(nickname string) error {
	removeIndex := -1
	for i := 0; i < len(g.Game.Players); i++ {
		if nickname == g.Game.Players[i].Nickname {
			removeIndex = i
			break
		}
	}

	g.Game.Players = append(g.Game.Players[:removeIndex], g.Game.Players[removeIndex+1:]...)

	return nil
}
