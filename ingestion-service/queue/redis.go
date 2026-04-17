package queue

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type RedisQueue struct {
	client *redis.Client
}

func NewRedisQueue(redisURL string) *RedisQueue {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatal("failed to parse redis URL:", err)
	}

	client := redis.NewClient(opts)

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatal("redis ping failed:", err)
	}

	log.Println("connected to redis")
	return &RedisQueue{client: client}
}

func (q *RedisQueue) PushEvent(ctx context.Context, eventID string) error {
	return q.client.XAdd(ctx, &redis.XAddArgs{
		Stream: "events",
		Values: map[string]interface{}{
			"event_id": eventID,
		},
	}).Err()
}