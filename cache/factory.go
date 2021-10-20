package cache

import "github.com/awe76/saga-coordinator/client"

func NewCache(client client.Client) Cache {
	return &cache{
		client: client,
	}
}
