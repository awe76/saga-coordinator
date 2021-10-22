package workflow

type Workflow struct {
	Name       string
	Start      string
	End        string
	Operations []Operation
}
