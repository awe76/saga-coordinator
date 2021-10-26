package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/awe76/saga-coordinator/cache"
	"github.com/awe76/saga-coordinator/producer"
	"github.com/awe76/saga-coordinator/workflow"
)

const indexKey = "workflow:index"

func getWorkflowKey(id int) string {
	return fmt.Sprintf("workflow:definition:%d", id)
}

func setWorkflow(id int, w workflow.Workflow, c cache.Cache) error {
	ctx := context.Background()
	key := getWorkflowKey(id)
	return c.Set(ctx, key, w)
}

func getWorkflow(id int, c cache.Cache) (workflow.Workflow, error) {
	ctx := context.Background()
	key := getWorkflowKey(id)

	var w workflow.Workflow
	raw, err := c.Get(ctx, key)
	if err != nil {
		return w, err
	}

	json.Unmarshal([]byte(raw), &w)
	return w, nil
}

func StartWorkflow(msg *sarama.ConsumerMessage, c cache.Cache, p producer.Producer) error {
	var w workflow.Workflow
	err := json.Unmarshal(msg.Value, &w)
	if err != nil {
		return err
	}

	id, err := workflow.ReserveID(indexKey, c)
	if err != nil {
		return err
	}

	err = setWorkflow(id, w, c)
	if err != nil {
		return err
	}

	proc := workflow.NewProcessor(c, p)
	err = proc.StartWorkflow(w, id)
	if err != nil {
		return err
	}

	return nil
}

func HandleOperationStart(msg *sarama.ConsumerMessage, c cache.Cache, p producer.Producer) error {
	var op workflow.OperationPayload
	err := json.Unmarshal(msg.Value, &op)
	if err != nil {
		return err
	}

	fmt.Printf("%s operation is started\n", op.Operation.Name)
	return p.SendMessage(workflow.WORKFLOW_OPERATION_COMPLETED, op)
}

func HandleOperationComplete(msg *sarama.ConsumerMessage, c cache.Cache, p producer.Producer) error {
	var op workflow.OperationPayload
	err := json.Unmarshal(msg.Value, &op)
	if err != nil {
		return err
	}

	fmt.Printf("%s operation is completed\n", op.Operation.Name)
	proc := workflow.NewProcessor(c, p)
	w, err := getWorkflow(op.ID, c)
	if err != nil {
		panic(err)
	}

	return proc.OnComplete(w, op)
}

func HandleOperationFailure(msg *sarama.ConsumerMessage, c cache.Cache, p producer.Producer) error {
	var op workflow.OperationPayload
	err := json.Unmarshal(msg.Value, &op)
	if err != nil {
		return err
	}

	proc := workflow.NewProcessor(c, p)
	w, err := getWorkflow(op.ID, c)
	if err != nil {
		return err
	}
	return proc.OnFailure(w, op)
}

func HandleWorkflowCompleted(msg *sarama.ConsumerMessage, c cache.Cache, p producer.Producer) error {
	var w workflow.WorkflowPayload
	err := json.Unmarshal(msg.Value, &w)
	if err != nil {
		return err
	}

	fmt.Printf("%s %d workflow is completed\n", w.Name, w.ID)
	return nil
}
