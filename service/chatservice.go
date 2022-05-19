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

func (c *ChatService) Subscribe(ctx context.Context, roomId string, chatChan chan string) error {
	subscribeChan := getSubscribeChan(roomId)
	isSubscribed, err := c.chatRepo.IsSubscribed(subscribeChan)
	if err != nil {
		return err
	}

	if !isSubscribed {
		if err := c.chatRepo.Subscribe(ctx, subscribeChan, chatChan); err != nil {
			return err 
		}
	}
	
	return nil 
	
}

func (c *ChatService) UnSubscribe(ctx context.Context, roomId string) error {
	subscribeChan := getSubscribeChan(roomId)
	if err := c.chatRepo.UnSubscribe(ctx, subscribeChan); err != nil {
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