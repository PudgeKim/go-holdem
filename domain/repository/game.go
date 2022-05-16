package repository

import (
	"context"

	"github.com/PudgeKim/go-holdem/domain/entity"
)

type GameRepository interface {
	GetGame(ctx context.Context, roomId string) (*entity.Game, error)
	SaveGame(ctx context.Context, roomId string, game *entity.Game) error
	CreateGame(ctx context.Context) (*entity.Game, error)
	DeleteGame(ctx context.Context) error 
	FindPlayer(ctx context.Context, roomId string, nickname string) (*entity.Player, error)
	AddPlayer(ctx context.Context, roomId string, player *entity.Player) error
}