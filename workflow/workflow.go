package workflow

type Workflow struct {
	Name       string
	Start      string
	End        string
	Operations []Operation
}

func (w *Workflow) toPayload(id int, isReversion bool) WorkflowPayload {
	return WorkflowPayload{
		ID:          id,
		Name:        w.Name,
		IsReversion: isReversion,
	}
}
