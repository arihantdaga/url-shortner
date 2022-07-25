package database

import (
	"github.com/go-redis/redis/v9"
)

func NewRedis(url string, pass string) (*redis.Client, error) {
	// TODO: Check for errors
	rdb := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: pass, // no password set
		DB:       0,    // use default DB
	})
	return rdb, nil
}
