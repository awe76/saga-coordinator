package workflow

type WorkflowPayload struct {
	ID          int
	IsReversion bool
	Name        string
	Data        map[string]interface{}
}
