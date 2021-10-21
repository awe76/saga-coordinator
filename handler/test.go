package handler

import (
	"fmt"

	"github.com/Shopify/sarama"
)

func HandleComment(msg *sarama.ConsumerMessage, msgCount int) {
	fmt.Printf("Received message Count %d: | Topic(%s) | Message(%s) \n", msgCount, string(msg.Topic), string(msg.Value))
}

func HandleError(err error) {
	fmt.Println(err)
}
