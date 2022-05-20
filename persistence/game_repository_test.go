package persistence

import (
	"context"
	"testing"

	"github.com/PudgeKim/go-holdem/cacheserver"
	"github.com/PudgeKim/go-holdem/domain/repository"
	"github.com/go-redis/redis/v8"
)

var redisClient = cacheserver.NewTestRedis()

func newGameRepo(redisClient *redis.Client) repository.GameRepository {
	return NewGameRepository(redisClient)
}
func TestGetGame(t *testing.T) {
	gameRepo := newGameRepo(redisClient)

	game, err := gameRepo.GetGame(context.Background(), "wrongRoomId")
	if err != nil {
		t.Log("err: ", err.Error())
	}
	if game != nil {
		t.Error("game should not exist")
	}
	t.Log("game: ", game)
}