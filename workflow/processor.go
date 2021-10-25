package workflow

import (
	"github.com/awe76/saga-coordinator/cache"
	"github.com/awe76/saga-coordinator/producer"
)

const (
	WORKFLOW_OPERATION_EXECUTE  = "workflow:operation:execute"
	WORKFLOW_OPERATION_ROLLBACK = "worfklow:operation:rollback"
	WORKFLOW_COMPLETED          = "workflow:completed"
	WORKFLOW_ROLLBACKED         = "workflow:rollbacked"
)

type RouteMap struct {
	Route map[string][]Operation
}

type processor struct {
	cache    cache.Cache
	producer producer.Producer
	workflow Workflow
	state    state
}

type Processor interface {
	StartWorkflow(w Workflow) error
	OnComplete(w Workflow, op OperationPayload) error
	OnFailure(w Workflow, op OperationPayload) error
}

func (p *processor) StartWorkflow(w Workflow) error {
	p.workflow = w
	id, err := reserveID("workflow:index", p.cache)
	if err != nil {
		return err
	}

	p.state = state{
		ID: id,
	}
	err = p.state.init(p.cache)
	if err != nil {
		return err
	}

	t := createDirectTracer(w, p.state, p.endWorkflow, p.spawnOperation)
	t.resolveWorkflow(w.Start)

	return nil
}

func (p *processor) OnComplete(w Workflow, op OperationPayload) error {
	p.workflow = w

	p.state = state{
		ID: op.ID,
	}
	err := p.state.update(p.cache, func(s *state) {
		removeOp(s.InProgress, op.Operation)

		if !s.IsRollback {
			addOp(s.Done, op.Operation)
		}
	})
	if err != nil {
		return err
	}

	if p.state.IsRollback {
		t := createReverseTracer(w, p.state, p.endWorkflow, p.spawnOperation)
		return t.resolveWorkflow(w.End)
	} else {
		t := createDirectTracer(w, p.state, p.endWorkflow, p.spawnOperation)
		return t.resolveWorkflow(w.Start)
	}
}

func (p *processor) OnFailure(w Workflow, op OperationPayload) error {
	p.workflow = w

	p.state = state{
		ID: op.ID,
	}
	err := p.state.update(p.cache, func(s *state) {
		removeOp(s.InProgress, op.Operation)

		s.IsRollback = true
	})
	if err != nil {
		return err
	}

	t := createReverseTracer(w, p.state, p.endWorkflow, p.spawnOperation)
	return t.resolveWorkflow(w.End)
}

func (p *processor) spawnOperation(op Operation) error {
	payload := OperationPayload{
		ID:         p.state.ID,
		IsRollback: p.state.IsRollback,
		Name:       p.workflow.Name,
		Operation:  op,
	}

	err := p.state.update(p.cache, func(s *state) {
		if s.IsRollback {
			removeOp(s.Done, op)
		}
	})
	if err != nil {
		return err
	}

	topic := WORKFLOW_OPERATION_EXECUTE
	if p.state.IsRollback {
		topic = WORKFLOW_OPERATION_ROLLBACK
	}

	return p.producer.SendMessage(topic, payload)
}

func (p *processor) endWorkflow() error {
	payload := WorkflowPayload{
		ID:         p.state.ID,
		IsRollback: p.state.IsRollback,
		Name:       p.workflow.Name,
	}

	topic := WORKFLOW_COMPLETED
	if p.state.IsRollback {
		topic = WORKFLOW_ROLLBACKED
	}

	return p.producer.SendMessage(topic, payload)
}
