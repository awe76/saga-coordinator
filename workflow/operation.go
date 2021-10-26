package workflow

import "fmt"

type Operation struct {
	Name string `json:"name"`
	From string `json:"from"`
	To   string `json:"to"`
}

func (op *Operation) getKey() string {
	return fmt.Sprintf("%s:%s:%s", op.Name, op.From, op.To)
}

func (op *Operation) toPayload(id int, w Workflow, isRollback bool) OperationPayload {
	return OperationPayload{
		ID:         id,
		Name:       w.Name,
		IsRollback: isRollback,
		Operation:  *op,
	}
}
