package workflow

type route = map[string][]Operation
type getOperationKey = func(op Operation) string

func hasOp(m map[string]Operation, op Operation) bool {
	key := op.getKey()
	_, found := m[key]
	return found
}

func addOp(m map[string]Operation, op Operation) {
	key := op.getKey()
	m[key] = op
}

func createRoute(operations []Operation, getKey getOperationKey) route {
	result := make(map[string][]Operation)
	for _, op := range operations {
		key := getKey(op)
		if next, found := result[key]; found {
			result[key] = append(next, op)
		} else {
			result[key] = append([]Operation{}, op)
		}
	}

	return result
}

func isReady(current string, to route, done map[string]Operation) bool {
	// get all operations should be finished to continue execution
	if blockers, found := to[current]; found {
		// for each potential blocker
		for _, blocker := range blockers {
			// if operation is not completed it is a blocker and the vertex is not ready for issuered operations execution
			if !hasOp(done, blocker) {
				return false
			}
		}
	}

	return true
}

func handleWorkflow(
	current string,
	start string,
	end string,
	from, to route,
	done map[string]Operation,
	inProgress map[string]Operation,
	endWorkflow func() error,
	spawnOperation func(op Operation) error,
) {
	if isReady(current, to, done) {
		if current == end {
			endWorkflow()
		} else {
			// get all operations started from the current vertex
			if ops, found := from[current]; found {
				// for each operation started from the current vertex
				for _, op := range ops {
					if hasOp(done, op) {
						// if operation is completed continue handling
						handleWorkflow(
							op.To,
							start,
							end,
							from,
							to,
							done,
							inProgress,
							endWorkflow,
							spawnOperation,
						)
					} else if !hasOp(inProgress, op) {
						// if operations is not completed spawn it and not in progress
						spawnOperation(op)
						// add operation to in progress
						addOp(inProgress, op)
					}
				}
			}
		}
	}
}
