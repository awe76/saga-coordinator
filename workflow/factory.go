package workflow

import (
	"github.com/awe76/saga-coordinator/cache"
	"github.com/awe76/saga-coordinator/producer"
)

func NewProcessor(cache cache.Cache, producer producer.Producer) Processor {
	return &processor{
		producer: producer,
		cache:    cache,
	}
}
