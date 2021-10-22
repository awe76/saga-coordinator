package workflow

import (
	"github.com/awe76/saga-coordinator/cache"
	"github.com/awe76/saga-coordinator/client"
	"github.com/awe76/saga-coordinator/producer"
)

func NewProcessor(client client.Client) Processor {
	producer := producer.NewProducer(client)
	cache := cache.NewCache(client)

	return &processor{
		producer: producer,
		cache:    cache,
	}
}
