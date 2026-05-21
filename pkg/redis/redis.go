package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/sumi-devs/canopy-social/canopy/pkg/config"
)

func NewClient(cfg *config.Config) (*redis.Client, error) {
	opts, err := redis.ParseURL(cfg.Redis.URL)
	if err != nil {
		return nil, err
	}
	opts.MaxRetries = cfg.Redis.MaxRetries
	client := redis.NewClient(opts)
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return client, nil
}
