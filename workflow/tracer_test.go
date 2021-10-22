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

	expected := map[string][]Operation{
		"s1": {
			operations[0],
			operations[1],
		},
		"s3": {
			operations[2],
		},
	}

	getFrom := func(op Operation) string {
		return op.From
	}

	route := createRoute(operations, getFrom)

	assert.Equal(t, route, expected)
}

func TestHandleWorkflow(t *testing.T) {
	getFrom := func(op Operation) string {
		return op.From
	}

	getTo := func(op Operation) string {
		return op.To
	}

	defaultOperations := []Operation{
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

	extendedOperations := []Operation{
		{
			Name: "op1",
			From: "s1",
			To:   "s2",
		},
		{
			Name: "op2",
			From: "s2",
			To:   "s3",
		},
		{
			Name: "op3",
			From: "s1",
			To:   "s3",
		},
		{
			Name: "op4",
			From: "s3",
			To:   "s4",
		},
		{
			Name: "op5",
			From: "s1",
			To:   "s4",
		},
	}

	var tests = map[string]struct {
		current    string
		start      string
		end        string
		operations []Operation
		done       []string
		inProgress []string
		expected   []string
		isFinished bool
	}{
		"should start worflow": {
			operations: defaultOperations,
			current:    "s1",
			start:      "s1",
			end:        "s2",
			done:       []string{},
			inProgress: []string{},
			expected:   []string{"op1", "op2"},
			isFinished: false,
		},
		"should spawn op3 if op1 and op2 are finished": {
			operations: defaultOperations,
			current:    "s1",
			start:      "s1",
			end:        "s2",
			done:       []string{"op1", "op2"},
			inProgress: []string{},
			expected:   []string{"op3"},
			isFinished: false,
		},
		"should compete workflow": {
			operations: defaultOperations,
			current:    "s1",
			start:      "s1",
			end:        "s2",
			done:       []string{"op1", "op2", "op3"},
			inProgress: []string{},
			expected:   []string{},
			isFinished: true,
		},
		"should start extended workflow": {
			operations: extendedOperations,
			current:    "s1",
			start:      "s1",
			end:        "s4",
			done:       []string{},
			inProgress: []string{},
			expected:   []string{"op1", "op3", "op5"},
			isFinished: false,
		},
		"should spawn op4 and op5 if op1 op2 and op3 are finished": {
			operations: extendedOperations,
			current:    "s1",
			start:      "s1",
			end:        "s4",
			done:       []string{"op1", "op2", "op3"},
			inProgress: []string{},
			expected:   []string{"op4", "op5"},
			isFinished: false,
		},
		"should spawn op4 if op1 op2 and op3 are finished and op5 is in progress": {
			operations: extendedOperations,
			current:    "s1",
			start:      "s1",
			end:        "s4",
			done:       []string{"op1", "op2", "op3"},
			inProgress: []string{"op5"},
			expected:   []string{"op4"},
			isFinished: false,
		},
		"should complete extended workflow": {
			operations: extendedOperations,
			current:    "s1",
			start:      "s1",
			end:        "s4",
			done:       []string{"op1", "op2", "op3", "op4", "op5"},
			inProgress: []string{},
			expected:   []string{},
			isFinished: true,
		},
		"should not complete extended workflow if op4 is in progress": {
			operations: extendedOperations,
			current:    "s1",
			start:      "s1",
			end:        "s4",
			done:       []string{"op1", "op2", "op3", "op5"},
			inProgress: []string{"op4"},
			expected:   []string{},
			isFinished: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			from := createRoute(tc.operations, getFrom)
			to := createRoute(tc.operations, getTo)
			done := make(map[string]Operation)
			inProgress := make(map[string]Operation)

			for _, name := range tc.done {
				for _, op := range tc.operations {
					if op.Name == name {
						addOp(done, op)
					}
				}
			}

			for _, name := range tc.inProgress {
				for _, op := range tc.operations {
					if op.Name == name {
						addOp(inProgress, op)
					}
				}
			}

			spawned := []string{}
			spawn := func(op Operation) error {
				spawned = append(spawned, op.Name)
				return nil
			}

			isFinished := false
			endHandler := func() error {
				isFinished = true
				return nil
			}

			handleWorkflow(tc.current, tc.start, tc.end, from, to, done, inProgress, endHandler, spawn)
			assert.Equal(t, tc.expected, spawned)
			assert.Equal(t, tc.isFinished, isFinished)
		})
	}
}
