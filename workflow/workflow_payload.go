package workflow

type WorkflowPayload struct {
	ID         int
	IsRollback bool
	Name       string
	Data       map[string]interface{}
}
