package consumer

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
	"github.com/awe76/saga-coordinator/cache"
	"github.com/awe76/saga-coordinator/client"
	"github.com/awe76/saga-coordinator/producer"
)

type MessageHandler = func(msg *sarama.ConsumerMessage, cache cache.Cache, producer producer.Producer) error
type ErrorHandler = func(err error)
type query struct {
	topic         string
	handleMessage MessageHandler
	handleError   ErrorHandler
}

type consumer struct {
	client   client.Client
	cache    cache.Cache
	producer producer.Producer
	topicCh  chan query
	doneCh   chan struct{}
	conn     sarama.Consumer
}

type Consumer interface {
	Start()
	HandleTopic(topic string, handleMessage MessageHandler, handleError ErrorHandler)
	WaitForInterrupt()
}

func (c *consumer) HandleTopic(topic string, handleMessage MessageHandler, handleError ErrorHandler) {
	query := query{
		topic:         topic,
		handleMessage: handleMessage,
		handleError:   handleError,
	}

	c.topicCh <- query
}

func (c *consumer) Start() {
	c.topicCh = make(chan query)
	// Get signal for finish
	c.doneCh = make(chan struct{})

	var err error

	c.conn, err = c.client.NewKafkaConsumer()
	if err != nil {
		panic(err)
	}

	fmt.Println("Consumer started ")
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case q := <-c.topicCh:
				// Calling ConsumePartition. It will open one connection per broker
				// and share it for all partitions that live on it.
				consumer, err := c.conn.ConsumePartition(q.topic, 0, sarama.OffsetNewest)
				if err != nil {
					fmt.Printf("%v\n", q.topic)
					panic(err)
				}

				go func() {
					for {
						select {
						case err := <-consumer.Errors():
							q.handleError(err)
						case msg := <-consumer.Messages():
							q.handleMessage(msg, c.cache, c.producer)
						}
					}
				}()
			case <-sigchan:
				fmt.Println("Interrupt is detected")
				c.doneCh <- struct{}{}
			}
		}
	}()
}

func (cr *consumer) WaitForInterrupt() {
	<-cr.doneCh
	fmt.Println("Consumer is stopped")

	if err := cr.conn.Close(); err != nil {
		panic(err)
	}
}
