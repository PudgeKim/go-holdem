package persistence

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/PudgeKim/go-holdem/domain/repository"
	"github.com/go-redis/redis/v8"
)

type ChatRepository struct {
	redisClient *redis.Client
	pubsubMap map[string]*redis.PubSub
}

func NewChatRepository(redisClient *redis.Client) repository.ChatRepository {
	return &ChatRepository{
		redisClient: redisClient,
		pubsubMap: make(map[string]*redis.PubSub),
	}
}

func (c *ChatRepository) Subscribe(ctx context.Context, subscribeChan string, chatChan chan string) error {
	pubsub := c.redisClient.Subscribe(ctx, subscribeChan)
	c.pubsubMap[subscribeChan] = pubsub
	go func() {
		ch := pubsub.Channel()
		for msg := range ch {
			chatChan<- msg.Payload
		}
	}()
	return nil 
}

func (c *ChatRepository) UnSubscribe(ctx context.Context, subscribeChan string) error {
	pubsub := c.pubsubMap[subscribeChan]
	if pubsub == nil {
		return errors.New("Invalid subscribe channel")
	}
	pubsub.Close()
	return nil 
}

func (c *ChatRepository) PublishMessage(ctx context.Context, subscribeChan string, nickname string, message string) error {
	pubsubMsg := ChatMessage{Nickname: nickname, Message: message}
	if err := c.redisClient.Publish(ctx, subscribeChan, pubsubMsg).Err(); err != nil {
		return err 
	}
	return nil 
}

func (c *ChatRepository) IsSubscribed(subscribeChan string) (bool, error) {
	if c.pubsubMap[subscribeChan] == nil {
		return false, nil 
	}
	return true, nil 
}

func (c *ChatRepository) GetAllSubscribeChannel() ([]string, error) {
	subscribeChannels := make([]string, len(c.pubsubMap))

	i := 0
	for key := range c.pubsubMap {
		subscribeChannels[i] = key 
		i++
	}

	return subscribeChannels, nil 
}

type ChatMessage struct {
	Nickname string 
	Message string 
}

// redis에 struct를 저장/가져오기위해 구현해야함
func (c ChatMessage) MarshalBinary() ([]byte, error) {
    return json.Marshal(c)
}

func (c *ChatMessage) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &c); err != nil {
		return err
	}
 
	return nil
}