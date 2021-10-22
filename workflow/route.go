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
			result[key] = append(make([]Operation, 0), op)
		}
	}

	return result
}

func isBlocked(blockers []Operation, done map[string]Operation) bool {
	for _, blocker := range blockers {
		key := blocker.getKey()
		if _, found := done[key]; !found {
			// some operation is not finished
			return true
		}
	}

	return false
}

func handleFollowers(
	current string,
	start string,
	end string,
	from, to route,
	done map[string]Operation,
	endWorkflow func() error,
	spawnOperation func(op Operation) error,
) {
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

func handleWorkflow(
	current string,
	start string,
	end string,
	from, to route,
	done map[string]Operation,
	endWorkflow func() error,
	spawnOperation func(op Operation) error,
) {
	// get all operations started from the current vertex
	if ops, found := from[current]; found {
		// spawn all operations worklow starts with
		if current == start {
			handleFollowers(current, start, end, from, to, done, endWorkflow, spawnOperation)
		} else {
			// for each operation started from the current vertex
			for _, op := range ops {
				// get all operations should be finshed to translate to op.To vertex
				if blockers, found := to[op.To]; found {
					// detect if all operations are finished to continue handling
					isBlocked := isBlocked(blockers, done)

					if !isBlocked {
						// if all operations are finshed and end vertex is reached
						if op.To == end {
							endWorkflow()
						} else {
							handleFollowers(op.To, start, end, from, to, done, endWorkflow, spawnOperation)
						}
					}
				}
			}
		}
	}
}
