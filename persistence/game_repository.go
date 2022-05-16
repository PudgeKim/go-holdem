package persistence

import (
	"context"
	"time"

	"github.com/PudgeKim/go-holdem/domain/entity"
	"github.com/PudgeKim/go-holdem/gameerror"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

const (
	ROOM_LIMIT = 7
	REDIS_TIME_DURATION = time.Hour * 144
)

type GameRepository struct {
	redisClient *redis.Client
}

func NewGameRepository(redisClient *redis.Client) *GameRepository {
	return &GameRepository{
		redisClient: redisClient,
	}
}

func (g *GameRepository) GetGame(ctx context.Context, roomId string) (*entity.Game, error) {
	stringCmd := g.redisClient.Get(ctx, roomId)
	if stringCmd.Err() != nil {
		return nil, stringCmd.Err()
	}

	var game entity.Game 

	if err := stringCmd.Scan(&game); err != nil {
		return nil, err 
	}

	return &game, nil 
}

func (g *GameRepository) SaveGame(ctx context.Context, roomId string, game *entity.Game) error {
	statusCmd := g.redisClient.Set(ctx, roomId, game, REDIS_TIME_DURATION)
	if statusCmd.Err() != nil {
		game.Undo()
		return statusCmd.Err()
	}
	return nil 
}

func (g *GameRepository) CreateGame(ctx context.Context) (*entity.Game, error) {
	roomId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	game := entity.NewGame(roomId, ROOM_LIMIT)
	return game, nil 
}

func (g *GameRepository) DeleteGame(ctx context.Context, roomId string) error {
	cmd := g.redisClient.Del(ctx, roomId)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil 
}

func (g *GameRepository) FindPlayer(ctx context.Context, roomId string, nickname string) (*entity.Player, error) {
	game, err := g.GetGame(ctx, roomId); if err != nil {
		return nil, err 
	}

	for _, p := range game.Players {
		if p.Nickname == nickname {
			return p, nil
		}
	}
	return nil, gameerror.NoPlayerExists
}

func (g *GameRepository) AddPlayer(ctx context.Context, roomId string, player *entity.Player) error {
	game, err := g.GetGame(ctx, roomId)
	// nil이면 플레이어가 이미 존재하는데 또 요청이 온 것
	if err == nil {
		return gameerror.PlayerAlreadyExists
	}

	if len(game.Players) > ROOM_LIMIT {
		return gameerror.PlayerLimitationError
	}

	game.Players = append(game.Players, player)

	if err := g.SaveGame(ctx, roomId, game); err != nil {
		return err 
	}

	game.SetMemento()

	return nil
}




