package workflow

import (
	"github.com/awe76/saga-coordinator/cache"
	"github.com/awe76/saga-coordinator/producer"
)

type RouteMap struct {
	Route map[string][]Operation
}

type processor struct {
	cache    cache.Cache
	producer producer.Producer
	workflow Workflow
	state    State
	from     route
	to       route
}

type Processor interface {
	Init(workflow Workflow) error
}

func (proc *processor) Init(workflow Workflow) error {
	id, err := reserveID(proc.cache)

	if err != nil {
		return err
	}

	state, err := createState(id, proc.cache)

	if err != nil {
		return err
	}

	proc.initRoutes(workflow.Operations)

	proc.workflow = workflow
	proc.state = state

	return nil
}

func (proc *processor) initRoutes(operations []Operation) {
	getFrom := func(op Operation) string {
		return op.From
	}

	getTo := func(op Operation) string {
		return op.To
	}

	proc.from = createRoute(operations, getFrom)
	proc.to = createRoute(operations, getTo)
}

func (proc *processor) handleOperationComplete(op Operation) {
	update := func(state State) {
		key := op.getKey()
		state.Done[key] = op
	}

	updateState(proc.state.ID, proc.cache, update)
}
