package main

import (
	"os"

	"github.com/awe76/saga-coordinator/client"
	"github.com/awe76/saga-coordinator/consumer"
	"github.com/awe76/saga-coordinator/gateway"
	"github.com/awe76/saga-coordinator/handler"
)

func main() {
	client := client.NewClient()
	args := os.Args[1:]

	if len(args) == 2 && args[1] == "consumer" {
		consumer := consumer.NewConsumer(client)
		consumer.Start()
		consumer.HandleTopic("comments", handler.HandleComment, handler.HandleError)
		consumer.WaitForInterrupt()

	} else {
		gateway := gateway.NewGateway(client)
		gateway.Run()
	}
}
