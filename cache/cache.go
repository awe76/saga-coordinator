package cache

import (
	"context"
	"encoding/json"

	"github.com/awe76/saga-coordinator/client"
)

type cache struct {
	client client.Client
}

type Cache interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string) (string, error)
	Remove(ctx context.Context, key string) error
}

func (c *cache) Set(ctx context.Context, key string, value interface{}) error {
	rdb := c.client.GetRedisClient()
	rawValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, key, string(rawValue), 0).Err()
}

func (c *cache) Get(ctx context.Context, key string) (string, error) {
	rdb := c.client.GetRedisClient()
	return rdb.Get(ctx, key).Result()
}

func (c *cache) Remove(ctx context.Context, key string) error {
	rdb := c.client.GetRedisClient()
	return rdb.Del(ctx, key).Err()
}
