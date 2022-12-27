package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var Dbctx = context.Background()

func Connect(db int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})
	pong, err := client.Ping(Dbctx).Result()
	fmt.Println(pong, err)
	if err != nil {
		log.Fatal(err)
	}
	return client
}
