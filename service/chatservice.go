package service

import (
	"context"
	"fmt"

	"github.com/PudgeKim/go-holdem/domain/repository"
)

type ChatService struct {
	chatRepo repository.ChatRepository
}

func NewChatService(chatRepo repository.ChatRepository) *ChatService {
	return &ChatService{
		chatRepo: chatRepo,
	}
}

func (c *ChatService) Subscribe(ctx context.Context, roomId string, userId int64, chatChan chan string) error {
	subscribeChan := getSubscribeChan(roomId)
	if err := c.chatRepo.Subscribe(ctx, subscribeChan, userId, chatChan); err != nil {
		return err 
	}
	
	return nil 
}

func (c *ChatService) UnSubscribe(ctx context.Context, roomId string, userId int64) error {
	subscribeChan := getSubscribeChan(roomId)
	if err := c.chatRepo.UnSubscribe(ctx, subscribeChan, userId); err != nil {
		return err 
	}
	return nil 
}

func (c *ChatService) PublishMessage(ctx context.Context, roomId, nickname, message string) error {
	subscribeChan := getSubscribeChan(roomId)
	if err := c.chatRepo.PublishMessage(ctx, subscribeChan, nickname, message); err != nil {
		return err 
	}
	return nil 
}

func getSubscribeChan(roomId string) string {
	return fmt.Sprintf("chat-%s", roomId)
}