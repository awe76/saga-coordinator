package consumer

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
	"github.com/awe76/saga-coordinator/client"
)

type MessageHandler = func(msg *sarama.ConsumerMessage, msgCount int)
type ErrorHandler = func(err error)
type query struct {
	topic         string
	handleMessage MessageHandler
	handleError   ErrorHandler
}

type consumer struct {
	client  client.Client
	topicCh chan query
	doneCh  chan struct{}
	conn    sarama.Consumer
}

type Consumer interface {
	Start()
	HandleTopic(topic string, handleMessage MessageHandler, handleError ErrorHandler)
	WaitForInterrupt()
}

func (cr *consumer) HandleTopic(topic string, handleMessage MessageHandler, handleError ErrorHandler) {
	query := query{
		topic:         topic,
		handleMessage: handleMessage,
		handleError:   handleError,
	}

	cr.topicCh <- query

}

func (cr *consumer) Start() {
	cr.topicCh = make(chan query)
	// Get signal for finish
	cr.doneCh = make(chan struct{})

	var err error

	cr.conn, err = cr.client.NewKafkaConsumer()
	if err != nil {
		panic(err)
	}

	fmt.Println("Consumer started ")
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case q := <-cr.topicCh:
				// Calling ConsumePartition. It will open one connection per broker
				// and share it for all partitions that live on it.
				consumer, err := cr.conn.ConsumePartition(q.topic, 0, sarama.OffsetOldest)
				if err != nil {
					panic(err)
				}

				// Count how many message processed
				msgCount := 0
				go func() {
					for {
						select {
						case err := <-consumer.Errors():
							q.handleError(err)
						case msg := <-consumer.Messages():
							msgCount++
							q.handleMessage(msg, msgCount)
						}
					}
				}()
			case <-sigchan:
				fmt.Println("Interrupt is detected")
				cr.doneCh <- struct{}{}
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
