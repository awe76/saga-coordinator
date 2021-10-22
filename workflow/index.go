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

func reserveID(cache cache.Cache) (int, error) {
	var index Index
	key := "worflow:index"

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
	rawNextIndex, err := json.Marshal(index)
	if err != nil {
		return index.ID, err
	}

	cache.Set(ctx, key, string(rawNextIndex))
	return index.ID, nil
}
