package persistence

import (
	"context"

	"github.com/PudgeKim/go-holdem/domain/entity"
	"github.com/PudgeKim/go-holdem/domain/repository"
	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) repository.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) FindOne(ctx context.Context, id int64) (*entity.User, error) {
	var user entity.User

	row := r.db.QueryRowx(`SELECT * FROM users WHERE id=$1`, id)
	err := row.StructScan(&user); if err != nil {
		return nil, err
	}

	return &user, nil 
}

func (r *userRepository) FindByNickname(ctx context.Context, nickname string) (*entity.User, error) {
	var user entity.User

	row := r.db.QueryRowx(`SELECT * FROM users WHERE nickname=$1`, nickname)
	err := row.StructScan(&user); if err != nil {
		return nil, err 
	}
	return &user, nil 
}

func (r *userRepository) Save(ctx context.Context, user *entity.User) error {
	if _, err := r.db.Exec(`INSERT INTO users (nickname, email, balance) VALUES ($1, $2, $3)`, user.Nickname, user.Email, 0); err != nil {
		return err 
	}
	return nil
}

func (r *userRepository) UpdateBalance(ctx context.Context, userId int64, balance uint64) (totalBalance uint64, err error) {
	var user entity.User

	row := r.db.QueryRowx("UPDATE users SET balance=$1 WHERE id=$2", balance, userId)
	if err := row.Scan(&user); err != nil {
		return 0, err 
	}
	return user.Balance, nil 
}

func (r *userRepository) UpdateMultipleBalance(ctx context.Context, userIdWithBalances []repository.UserIdWithBalance) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err 
	}

	// TODO
	// 나중에 batch update로 바꾸자 
	for i:=0; i<len(userIdWithBalances); i++ {
		_, err := tx.Exec(`UPDATE users SET balance=$1 WHERE id=$2`, userIdWithBalances[i].Balance, userIdWithBalances[i].UserId)
		if err != nil {
			tx.Rollback()
			return err 
		}
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err 
	}

	return nil
}