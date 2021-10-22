package workflow

type route = map[string][]Operation
type getOperationKey = func(op Operation) string

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
			key := blocker.getKey()
			// if operation is not completed it is a blocker and the vertex is not ready
			if _, found := done[key]; !found {
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
					key := op.getKey()
					if _, found := done[key]; found {
						// if operation is completed continue handling
						handleWorkflow(op.To, start, end, from, to, done, endWorkflow, spawnOperation)
					} else {
						// if operations is not completed spawn it
						spawnOperation(op)
					}
				}
			}
		}
	}
}
