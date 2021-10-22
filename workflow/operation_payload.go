package workflow

import (
	"fmt"

	"github.com/awe76/saga-coordinator/producer"
)

type OperationPayload struct {
	ID        int
	Name      string
	Operation Operation
	Context   map[string]map[string]interface{}
}

func (op *OperationPayload) getTopicKey() string {
	return fmt.Sprintf("operation")
}

func (op *OperationPayload) push(producer producer.Producer) error {
	topic := op.getTopicKey()
	return pushData(producer, topic, op)
}
