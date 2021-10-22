package consumer

import (
	"github.com/awe76/saga-coordinator/cache"
	"github.com/awe76/saga-coordinator/client"
)

func NewConsumer(client client.Client) Consumer {
	cache := cache.NewCache(client)
	return &consumer{
		client: client,
		cache:  cache,
	}
}
