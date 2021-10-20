package client

import (
	"github.com/Shopify/sarama"
	"github.com/go-redis/redis/v8"
)

func NewClient() Client {
	redisOptions := &redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}
	redisClient := redis.NewClient(redisOptions)

	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	kafkaConfig.Producer.Retry.Max = 5
	kafkaConfig.Consumer.Return.Errors = true

	return &client{
		brokerUrls:  []string{"localhost:29092"},
		kafkaConfig: kafkaConfig,
		redisClient: redisClient,
	}
}
