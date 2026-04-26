package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	RDB *redis.Client
}

func New(redisURL string) *Client {
	opt, _ := redis.ParseURL(redisURL)

	rdb := redis.NewClient(opt)

	return &Client{RDB: rdb}
}

func (c *Client) Ping(ctx context.Context) error {
	return c.RDB.Ping(ctx).Err()
}
