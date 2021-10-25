package workflow

import (
	"testing"

	"github.com/awe76/saga-coordinator/cache/cachemock"
	"github.com/awe76/saga-coordinator/producer/producermock"
	"github.com/stretchr/testify/assert"
)

func TestProcessor(t *testing.T) {

	ops := []Operation{
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

	defaultWorkflow := Workflow{
		Name:       "default workflow",
		Start:      "s1",
		End:        "s2",
		Operations: ops,
	}

	type step struct {
		action   func(t *testing.T, w Workflow, p *processor)
		validate func(t *testing.T, w Workflow, p *producermock.ProducerMock)
	}

	var tests = map[string]struct {
		w     Workflow
		steps []step
	}{
		"default workflow is completed": {
			w: defaultWorkflow,
			steps: []step{
				{
					action: func(t *testing.T, w Workflow, p *processor) {
						assert.NoError(t, p.StartWorkflow(w))
					},
					validate: func(t *testing.T, w Workflow, p *producermock.ProducerMock) {
						op1 := ops[0].toPayload(1, w, false)
						assert.True(t, p.Has(WORKFLOW_OPERATION_EXECUTE, op1))

						op2 := ops[1].toPayload(1, w, false)
						assert.True(t, p.Has(WORKFLOW_OPERATION_EXECUTE, op2))
					},
				},
				{
					action: func(t *testing.T, w Workflow, p *processor) {
						op1 := ops[0].toPayload(1, w, false)
						assert.NoError(t, p.OnComplete(w, op1))

						op2 := ops[1].toPayload(1, w, false)
						assert.NoError(t, p.OnComplete(w, op2))
					},
					validate: func(t *testing.T, w Workflow, p *producermock.ProducerMock) {
						op3 := ops[2].toPayload(1, w, false)
						assert.True(t, p.Has(WORKFLOW_OPERATION_EXECUTE, op3))
					},
				},
				{
					action: func(t *testing.T, w Workflow, p *processor) {
						op3 := ops[2].toPayload(1, w, false)
						assert.NoError(t, p.OnComplete(w, op3))
					},
					validate: func(t *testing.T, w Workflow, p *producermock.ProducerMock) {
						wp := w.toPayload(1, false)
						assert.True(t, p.Has(WORKFLOW_COMPLETED, wp))
					},
				},
			},
		},
		"default workflow is rollbacked": {
			w: defaultWorkflow,
			steps: []step{
				{
					action: func(t *testing.T, w Workflow, p *processor) {
						assert.NoError(t, p.StartWorkflow(w))
					},
					validate: func(t *testing.T, w Workflow, p *producermock.ProducerMock) {
						op1 := ops[0].toPayload(1, w, false)
						assert.True(t, p.Has(WORKFLOW_OPERATION_EXECUTE, op1))

						op2 := ops[1].toPayload(1, w, false)
						assert.True(t, p.Has(WORKFLOW_OPERATION_EXECUTE, op2))
					},
				},
				{
					action: func(t *testing.T, w Workflow, p *processor) {
						op1 := ops[0].toPayload(1, w, false)
						assert.NoError(t, p.OnComplete(w, op1))

						op2 := ops[1].toPayload(1, w, false)
						assert.NoError(t, p.OnComplete(w, op2))
					},
					validate: func(t *testing.T, w Workflow, p *producermock.ProducerMock) {
						op3 := ops[2].toPayload(1, w, false)
						assert.True(t, p.Has(WORKFLOW_OPERATION_EXECUTE, op3))
					},
				},
				{
					action: func(t *testing.T, w Workflow, p *processor) {
						op3 := ops[2].toPayload(1, w, false)
						assert.NoError(t, p.OnFailure(w, op3))
					},
					validate: func(t *testing.T, w Workflow, p *producermock.ProducerMock) {
						op1 := ops[0].toPayload(1, w, true)
						assert.True(t, p.Has(WORKFLOW_OPERATION_ROLLBACK, op1))

						op2 := ops[1].toPayload(1, w, true)
						assert.True(t, p.Has(WORKFLOW_OPERATION_ROLLBACK, op2))
					},
				},
				{
					action: func(t *testing.T, w Workflow, p *processor) {
						op1 := ops[0].toPayload(1, w, true)
						assert.NoError(t, p.OnComplete(w, op1))

						op2 := ops[1].toPayload(1, w, true)
						assert.NoError(t, p.OnComplete(w, op2))
					},
					validate: func(t *testing.T, w Workflow, p *producermock.ProducerMock) {
						wp := w.toPayload(1, true)
						assert.True(t, p.Has(WORKFLOW_ROLLBACKED, wp))
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			cache := cachemock.New()
			producer := producermock.New()

			for _, step := range tc.steps {
				proc := &processor{
					cache:    cache,
					producer: producer,
				}
				step.action(t, tc.w, proc)
				step.validate(t, tc.w, producer)
			}
		})
	}

}
