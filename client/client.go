package client

import (
	"github.com/Shopify/sarama"
	"github.com/go-redis/redis/v8"
)

type client struct {
	brokerUrls  []string
	kafkaConfig *sarama.Config
	redisClient *redis.Client
}

type Client interface {
	NewKafkaProducer() (sarama.SyncProducer, error)
	NewKafkaConsumer() (sarama.Consumer, error)
	GetRedisClient() *redis.Client
}

func (c *client) NewKafkaProducer() (sarama.SyncProducer, error) {
	return sarama.NewSyncProducer(c.brokerUrls, c.kafkaConfig)
}

func (c *client) NewKafkaConsumer() (sarama.Consumer, error) {
	return sarama.NewConsumer(c.brokerUrls, c.kafkaConfig)
}

func (c *client) GetRedisClient() *redis.Client {
	return c.redisClient
}
