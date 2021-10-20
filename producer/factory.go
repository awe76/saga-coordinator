package producer

import "github.com/awe76/saga-coordinator/client"

func NewProducer(client client.Client) Producer {
	return &producer{
		client: client,
	}
}
