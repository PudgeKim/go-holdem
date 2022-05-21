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

// it should be nil 
func TestGetGame(t *testing.T) {
	gameRepo := newGameRepo(redisClient)

	_, err := gameRepo.GetGame(context.Background(), "wrongRoomId")
	if err != nil {
		if err.Error() != redis.Nil.Error() {
			t.Log("err: ", err.Error())
		}
	}

}