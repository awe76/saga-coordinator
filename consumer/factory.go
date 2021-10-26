package consumer

import (
	"github.com/awe76/saga-coordinator/cache"
	"github.com/awe76/saga-coordinator/client"
	"github.com/awe76/saga-coordinator/producer"
)

func NewConsumer(client client.Client) Consumer {
	cache := cache.NewCache(client)
	producer := producer.NewProducer(client)
	return &consumer{
		client:   client,
		cache:    cache,
		producer: producer,
	}
}
