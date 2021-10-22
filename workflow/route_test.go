package workflow

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateRoute(t *testing.T) {
	operations := []Operation{
		{
			Name: "op1",
			From: "s1",
			To:   "s2",
		},
		{
			Name: "op2",
			From: "s1",
			To:   "s3",
		},
		{
			Name: "op3",
			From: "s3",
			To:   "s2",
		},
	}

	expected := make(map[string][]Operation)
	expected["s1"] = append(make([]Operation, 0), operations[0], operations[1])
	expected["s3"] = append(make([]Operation, 0), operations[2])

	getKey := func(op Operation) string {
		return op.From
	}

	route := createRoute(operations, getKey)

	assert.Equal(t, route, expected)
}

func TestHandleWorkflow(t *testing.T) {
	operations := []Operation{
		{
			Name: "op1",
			From: "s1",
			To:   "s2",
		},
		{
			Name: "op2",
			From: "s1",
			To:   "s3",
		},
		{
			Name: "op3",
			From: "s3",
			To:   "s2",
		},
	}

	getFrom := func(op Operation) string {
		return op.From
	}

	getTo := func(op Operation) string {
		return op.To
	}

	from := createRoute(operations, getFrom)
	to := createRoute(operations, getTo)
	done := make(map[string]Operation)

	requested := make([]Operation, 0)

	endWorkflow := func() error {
		t.Error("not expected to end")
		return nil
	}

	spawnOperation := func(op Operation) error {
		requested = append(requested, op)
		return nil
	}

	expected := []Operation{
		operations[0],
		operations[1],
	}

	handleWorkflow("s1", "s1", "s2", from, to, done, endWorkflow, spawnOperation)

	assert.Equal(t, expected, requested)
}

func TestHandleWorkflowIfCompleted(t *testing.T) {
	operations := []Operation{
		{
			Name: "op1",
			From: "s1",
			To:   "s2",
		},
		{
			Name: "op2",
			From: "s1",
			To:   "s3",
		},
		{
			Name: "op3",
			From: "s3",
			To:   "s2",
		},
	}

	getFrom := func(op Operation) string {
		return op.From
	}

	getTo := func(op Operation) string {
		return op.To
	}

	from := createRoute(operations, getFrom)
	to := createRoute(operations, getTo)
	done := make(map[string]Operation)
	done[operations[0].getKey()] = operations[0]
	done[operations[2].getKey()] = operations[2]

	requested := make([]Operation, 0)

	endWorkflow := func() error {
		t.Error("not expected to end")
		return nil
	}

	spawnOperation := func(op Operation) error {
		requested = append(requested, op)
		return nil
	}

	expected := []Operation{
		operations[1],
	}

	handleWorkflow("s1", "s1", "s2", from, to, done, endWorkflow, spawnOperation)

	assert.Equal(t, expected, requested)
}
