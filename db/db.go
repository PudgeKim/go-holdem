package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	mainPort     = 5433
	testPort = 5434
	user     = "root"
	password = "mypassword"
	dbname   = "goholdem"
)

var DBConfig = fmt.Sprintf("port=%d host=%s user=%s password=%s dbname=%s sslmode=disable",
mainPort, host, user, password, dbname)

var TestDBConfig = fmt.Sprintf("port=%d host=%s user=%s password=%s dbname=%s sslmode=disable",
testPort, host, user, password, dbname)

func NewPostgresDB(psqlInfo string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err 
	}

	return db, nil
}