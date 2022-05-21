package persistence

import (
	"context"
	"fmt"
	"testing"

	"github.com/PudgeKim/go-holdem/db"
	"github.com/PudgeKim/go-holdem/domain/entity"
	"github.com/PudgeKim/go-holdem/domain/repository"
	"github.com/jmoiron/sqlx"
)

var createTableSchema = `
CREATE TABLE IF NOT EXISTS users (
	id serial PRIMARY KEY,
    nickname text,
    email text,
    balance bigint,
	password text
);
`

var dropTableSchema = `
DROP TABLE IF EXISTS users
`

var DB *sqlx.DB

func newDB() *sqlx.DB {
	if DB != nil {
		return DB 
	}

	conn, err := db.NewPostgresDB(db.TestDBConfig)
	if err != nil {
		fmt.Println("init db failed!")
		panic(err)
	}

	DB = conn 
	return DB 
}

func createTable(db *sqlx.DB) {
	db.MustExec(createTableSchema)
}

func dropTable(db *sqlx.DB) {
	db.MustExec(dropTableSchema)
}

func newUserRepo(conn *sqlx.DB) repository.UserRepository {
	userRepo := NewUserRepository(conn)
	return userRepo
}

func TestSaveUser(t *testing.T) {
	db := newDB()
	userRepo := newUserRepo(db)
	createTable(db)

	user := entity.NewUser("kim", "kim@gmail.com", "mypassword")
	if err := userRepo.Save(context.Background(), user); err != nil {
		t.Error("user was not saved\n", err.Error())
	}

	dropTable(db)
}

func TestUpdateMultipleBalance(t *testing.T) {
	db := newDB()
	userRepo := newUserRepo(db)
	createTable(db)

	user1 := entity.NewUser("kim", "kim@gmail.com",
"mypassword")
	userRepo.Save(context.Background(), user1)
	user2 := entity.NewUser("han", "han@gmail.com", "mypassword")
	userRepo.Save(context.Background(), user2)

	var userIdWithBalances []repository.UserIdWithBalance
	
	user1Id, user2Id := 1, 2
	b1 := repository.UserIdWithBalance{int64(user1Id), 100}
	b2 := repository.UserIdWithBalance{int64(user2Id), 200}
	userIdWithBalances = append(userIdWithBalances, b1, b2)

	if err := userRepo.UpdateMultipleBalance(context.Background(), userIdWithBalances); err != nil {
		t.Error("UpdateMultipleBalance errorrrr\n", err.Error())
	}

	foundUser1, err := userRepo.FindOne(context.Background(), int64(user1Id))
	if err != nil {
		t.Error("FindErr: ", err.Error())
	}
	foundUser2, _ := userRepo.FindOne(context.Background(), int64(user2Id))

	if foundUser1.Balance != 100 && foundUser2.Balance != 200 {
		t.Error("balance is not updated")
	}

	dropTable(db)
}