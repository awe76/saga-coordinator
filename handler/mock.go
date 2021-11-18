package handler

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/Shopify/sarama"
	"github.com/awe76/saga-coordinator/cache"
	"github.com/awe76/saga-coordinator/producer"
	"github.com/awe76/saga-coordinator/workflow"
)

func TestHandleOperationStart(msg *sarama.ConsumerMessage, c cache.Cache, p producer.Producer) error {
	var op workflow.OperationPayload
	err := json.Unmarshal(msg.Value, &op)
	if err != nil {
		return err
	}

	if op.IsRollback {
		fmt.Printf("%s operation rollback is started\n", op.Operation.Name)
	} else {
		fmt.Printf("%s operation is started\n", op.Operation.Name)
	}

	rand.Seed(time.Now().UnixNano())

	pause := rand.Intn(500)
	// sleep for some random time
	time.Sleep(time.Duration(pause) * time.Millisecond)

	op.Payload = rand.Float32()

	// randomly complete or fault the operation
	if op.IsRollback || rand.Float32() < 0.8 {
		return p.SendMessage(workflow.WORKFLOW_OPERATION_COMPLETED, op)
	} else {
		return p.SendMessage(workflow.WORKFLOW_OPERATION_FAILED, op)
	}
}
