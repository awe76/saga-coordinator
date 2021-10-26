package workflow

type Workflow struct {
	Name       string      `json:"name"`
	Start      string      `json:"start"`
	End        string      `json:"end"`
	Operations []Operation `json:"operations"`
}

func (w *Workflow) toPayload(id int, isReversion bool) WorkflowPayload {
	return WorkflowPayload{
		ID:         id,
		Name:       w.Name,
		IsRollback: isReversion,
	}
}
