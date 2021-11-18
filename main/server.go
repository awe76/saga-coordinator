package main

import (
	"os"

	"github.com/awe76/saga-coordinator/client"
	"github.com/awe76/saga-coordinator/consumer"
	"github.com/awe76/saga-coordinator/gateway"
	"github.com/awe76/saga-coordinator/handler"
	"github.com/awe76/saga-coordinator/portfolio"
	"github.com/awe76/saga-coordinator/workflow"
)

func main() {
	client := client.NewClient()
	args := os.Args[1:]

	if len(args) >= 2 && args[1] == "consumer" {
		consumer := consumer.NewConsumer(client)
		consumer.Start()
		consumer.HandleTopic(workflow.WORKFLOW_START, handler.StartWorkflow, handler.HandleError)

		if len(args) == 3 && args[2] == "trace" {
			consumer.HandleTopic(workflow.WORKFLOW_OPERATION_START, handler.TestHandleOperationStart, handler.HandleError)
		} else {
			consumer.HandleTopic(workflow.WORKFLOW_OPERATION_START, handler.HandleOperationStart, handler.HandleError)

			consumer.HandleTopic("create-building-map-portfolio", portfolio.HandleCreateBuildingMapPotrfolio, handler.HandleError)

			consumer.HandleTopic("create-building-map-portfolio-rollback", portfolio.HandleCreateBuildingMapPotrfolioRollback, handler.HandleError)

			consumer.HandleTopic("create-portfolio", portfolio.HandleCreatePotfolio, handler.HandleError)

			consumer.HandleTopic("create-portfolio-rollback", portfolio.HandleCreatePotfolioRollback, handler.HandleError)
		}
		consumer.HandleTopic(workflow.WORKFLOW_OPERATION_COMPLETED, handler.HandleOperationComplete, handler.HandleError)
		consumer.HandleTopic(workflow.WORKFLOW_OPERATION_FAILED, handler.HandleOperationFailure, handler.HandleError)
		consumer.HandleTopic(workflow.WORKFLOW_COMPLETED, handler.HandleWorkflowCompleted, handler.HandleError)
		consumer.HandleTopic(workflow.WORKFLOW_ROLLBACKED, handler.HandleWorkflowRollbacked, handler.HandleError)
		consumer.WaitForInterrupt()

	} else {
		gateway := gateway.NewGateway(client)
		gateway.Run()
	}
}
