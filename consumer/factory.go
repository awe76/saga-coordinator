package consumer

import "github.com/awe76/saga-coordinator/client"

func NewConsumer(client client.Client) Consumer {
	return &consumer{
		client: client,
	}
}
