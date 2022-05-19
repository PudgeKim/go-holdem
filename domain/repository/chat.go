package repository

import (
	"context"
)

type ChatRepository interface {
	Subscribe(ctx context.Context, subscribeChan string, chatChan chan string) error
	UnSubscribe(ctx context.Context, subscribeChan string) error 
	PublishMessage(ctx context.Context, subscribeChan string, nickname string, message string) error
	IsSubscribed(subscribeChan string) (bool, error)
	GetAllSubscribeChannel() ([]string, error)
}