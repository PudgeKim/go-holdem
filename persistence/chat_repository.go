package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/PudgeKim/go-holdem/domain/repository"
	"github.com/go-redis/redis/v8"
)

type ChatRepository struct {
	redisClient *redis.Client
	pubsubMap map[string]*redis.PubSub
	chatChanMap map[string][]chan string 
}

func NewChatRepository(redisClient *redis.Client) repository.ChatRepository {
	return &ChatRepository{
		redisClient: redisClient,
		pubsubMap: make(map[string]*redis.PubSub),
		chatChanMap: make(map[string][]chan string),
	}
}

func (c *ChatRepository) Subscribe(ctx context.Context, subscribeChan string, chatChan chan string) error {
	if c.pubsubMap[subscribeChan] == nil {
		pubsub := c.redisClient.Subscribe(ctx, subscribeChan)
		c.pubsubMap[subscribeChan] = pubsub
	}

	c.chatChanMap[subscribeChan] = append(c.chatChanMap[subscribeChan], chatChan)
	go c.handleMessage(subscribeChan)
	return nil 
}

func (c *ChatRepository) handleMessage(subscribeChan string) {
	pubsub := c.pubsubMap[subscribeChan]
	ch := pubsub.Channel()
	for msg := range ch {
		// 같은 방안에 있는 사람들에게 broadcast 해줘야함 
		for _, chatChan := range c.chatChanMap[subscribeChan] {
			fmt.Println(msg.Payload, c.chatChanMap[subscribeChan])
			chatChan <-msg.Payload
		}
		
	}
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
	chatMsg := ChatMessage{Nickname: nickname, Message: message}
	if err := c.redisClient.Publish(ctx, subscribeChan, chatMsg).Err(); err != nil {
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