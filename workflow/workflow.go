package workflow

type Workflow struct {
	Name       string      `json:"name"`
	Start      string      `json:"start"`
	End        string      `json:"end"`
	Operations []Operation `json:"operations"`
	Payload    interface{} `json:"payload"`
}

func (w *Workflow) toPayload(id int, isReversion bool, data map[string]map[string]interface{}) WorkflowPayload {
	return WorkflowPayload{
		ID:         id,
		Name:       w.Name,
		Data:       data,
		IsRollback: isReversion,
	}
}
