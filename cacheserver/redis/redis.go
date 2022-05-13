package redisclient

import "github.com/go-redis/redis/v8"

const (
	address = "localhost:6380"
	testAddress = "localhost:6381"
	password = ""
	testPassword = ""
)

var redisOptions = &redis.Options{
	Addr: address,
	Password: password,
	DB: 0,
}

var testRedisOptions = &redis.Options{
	Addr: testAddress,
	Password: testPassword,
	DB: 0,
}

func New(redisOptions *redis.Options) *redis.Client {
	return redis.NewClient(redisOptions)
}