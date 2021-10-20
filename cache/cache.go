package cache

import (
	"context"

	"github.com/awe76/saga-coordinator/client"
)

type cache struct {
	client client.Client
}

type Cache interface {
	Set(ctx context.Context, key string, value string) error
}

func (c *cache) Set(ctx context.Context, key string, value string) error {
	rdb := c.client.GetRedisClient()
	return rdb.Set(ctx, key, value, 0).Err()
}
