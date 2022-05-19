package cacheserver

import "github.com/go-redis/redis/v8"

const (
	address = "localhost:6380"
	testAddress = "localhost:6381"
	password = ""
	testPassword = ""
)

var redisConfig = &redis.Options{
	Addr: address,
	Password: password,
	DB: 0,
}

var testRedisConfig = &redis.Options{
	Addr: testAddress,
	Password: testPassword,
	DB: 0,
}

func NewRedis() *redis.Client {
	return redis.NewClient(redisConfig)
}

func NewTestRedis() *redis.Client {
	return redis.NewClient(testRedisConfig)
}