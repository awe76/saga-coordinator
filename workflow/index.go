package workflow

import (
	"context"
	"encoding/json"

	"github.com/awe76/saga-coordinator/cache"
	"github.com/go-redis/redis/v8"
)

type Index struct {
	ID int
}

func reserveID(key string, cache cache.Cache) (int, error) {
	var index Index

	ctx := context.Background()
	rawIndex, err := cache.Get(ctx, key)

	if err == redis.Nil {
		index = Index{
			ID: 0,
		}
	} else if err != nil {
		return index.ID, err
	} else {
		json.Unmarshal([]byte(rawIndex), &index)
	}

	index.ID++

	cache.Set(ctx, key, index)
	return index.ID, nil
}
