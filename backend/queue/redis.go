package queue

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var (
	Ctx = context.Background()
	RDB *redis.Client
)

func ConnectRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	if err := RDB.Ping(Ctx).Err(); err != nil {
		log.Fatal("failed to connect redis: ", err)
	}

	log.Println("REDIS connected")
}
