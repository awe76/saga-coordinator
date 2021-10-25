package producer

import (
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/awe76/saga-coordinator/client"
)

type producer struct {
	client client.Client
}

type Producer interface {
	SendMessage(topic string, message interface{}) error
}

func (p *producer) SendMessage(topic string, message interface{}) error {
	conn, err := p.client.NewKafkaProducer()

	if err != nil {
		return err
	}

	defer conn.Close()

	raw, err := json.Marshal(message)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(raw),
	}

	partition, offset, err := conn.SendMessage(msg)
	if err != nil {
		return err
	}

	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)

	return nil
}
