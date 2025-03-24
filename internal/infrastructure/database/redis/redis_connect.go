package database

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"ord_crm/config"
)

func RedisConnect() (*redis.Client, error) {
	token, err := config.GetToken()
	if err != nil {
		fmt.Printf("Error: %s", err)
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: token.RedisPass,
		DB:       0,
	})

	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	fmt.Println("Redis connected")

	return client, nil
}
