package producer

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/awe76/saga-coordinator/client"
)

type producer struct {
	client client.Client
}

type Producer interface {
	Push(topic string, message []byte) error
}

func (p *producer) Push(topic string, message []byte) error {
	conn, err := p.client.NewKafkaProducer()

	if err != nil {
		return err
	}

	defer conn.Close()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := conn.SendMessage(msg)
	if err != nil {
		return err
	}

	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)

	return nil
}
