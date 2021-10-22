package workflow

import (
	"fmt"

	"github.com/awe76/saga-coordinator/producer"
)

type WorkflowPayload struct {
	ID   int
	Name string
	Data map[string]interface{}
}

func (wf *WorkflowPayload) getTopicKey() string {
	return fmt.Sprintf("workflow")
}

func (wf *WorkflowPayload) push(producer producer.Producer) error {
	topic := wf.getTopicKey()
	return pushData(producer, topic, wf)
}
