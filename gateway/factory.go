package gateway

import (
	"github.com/awe76/saga-coordinator/client"
	"github.com/awe76/saga-coordinator/producer"
)

func NewGateway(client client.Client) Gateway {
	producer := producer.NewProducer(client)
	handler := &handler{
		producer: producer,
	}

	return &gateway{
		handler: handler,
	}
}
