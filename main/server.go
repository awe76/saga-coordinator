package main

import (
	"os"

	"github.com/awe76/saga-coordinator/client"
	"github.com/awe76/saga-coordinator/consumer"
	"github.com/awe76/saga-coordinator/gateway"
)

func main() {
	client := client.NewClient()
	args := os.Args[1:]

	if len(args) == 2 && args[1] == "consumer" {
		consumer := consumer.NewConsumer(client)
		consumer.Run()

	} else {
		gateway := gateway.NewGateway(client)
		gateway.Run()
	}
}
