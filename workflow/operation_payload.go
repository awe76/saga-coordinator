package workflow

type OperationPayload struct {
	ID          int
	IsReversion bool
	Name        string
	Operation   Operation
}
