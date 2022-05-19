package repository

import (
	"context"

	"github.com/PudgeKim/go-holdem/domain/entity"
)

type UserRepository interface {
	FindOne(ctx context.Context, id int64) (*entity.User, error)
	FindByNickname(ctx context.Context, nickname string) (*entity.User, error)
	Save(ctx context.Context, user *entity.User) error
	UpdateBalance(ctx context.Context, userId int64, balance uint64) (totalBalance uint64, err error)
	UpdateMultipleBalance(ctx context.Context, userIdWithBalances []UserIdWithBalance) error 
}

type UserIdWithBalance struct {
	UserId int64 
	Balance uint64 
}

func NewUserIdWithBalance(userId int64, balance uint64) UserIdWithBalance {
	return UserIdWithBalance{userId, balance}
}