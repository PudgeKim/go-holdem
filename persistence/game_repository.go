package persistence

import (
	"context"
	"time"

	"github.com/PudgeKim/go-holdem/domain/entity"
	"github.com/PudgeKim/go-holdem/domain/repository"
	"github.com/PudgeKim/go-holdem/errors/gameerror"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

const (
	ROOM_LIMIT = 7
	REDIS_TIME_DURATION = time.Hour * 144
)

type gameRepository struct {
	redisClient *redis.Client
}

func NewGameRepository(redisClient *redis.Client) repository.GameRepository {
	return &gameRepository{
		redisClient: redisClient,
	}
}

func (g *gameRepository) GetGame(ctx context.Context, roomId string) (*entity.Game, error) {
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

func (g *gameRepository) SaveGame(ctx context.Context, roomId string, game *entity.Game) error {
	statusCmd := g.redisClient.Set(ctx, roomId, game, REDIS_TIME_DURATION)
	if statusCmd.Err() != nil {
		return statusCmd.Err()
	}
	return nil 
}

func (g *gameRepository) CreateGame(ctx context.Context, hostPlayer *entity.Player, minBetAmount uint64) (*entity.Game, string, error) {
	roomId, err := uuid.NewRandom()
	if err != nil {
		return nil, "", err
	}

	game := entity.NewGame(roomId, ROOM_LIMIT, hostPlayer, minBetAmount)
	return game, roomId.String(), nil 
}

func (g *gameRepository) DeleteGame(ctx context.Context, roomId string) error {
	cmd := g.redisClient.Del(ctx, roomId)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil 
}

func (g *gameRepository) FindPlayer(ctx context.Context, roomId string, nickname string) (*entity.Player, error) {
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

func (g *gameRepository) AddPlayer(ctx context.Context, roomId string, player *entity.Player) error {
	game, err := g.GetGame(ctx, roomId)
	if err != nil {
		return err
	}

	if game.IsPlayerExist(player.Nickname) {
		return gameerror.PlayerAlreadyExists
	}

	if len(game.Players) > ROOM_LIMIT {
		return gameerror.PlayerLimitationError
	}

	game.Players = append(game.Players, player)

	if err := g.SaveGame(ctx, roomId, game); err != nil {
		return err 
	}

	return nil
}






