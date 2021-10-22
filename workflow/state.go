package workflow

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/awe76/saga-coordinator/cache"
)

type State struct {
	ID         int
	Done       map[string]Operation
	Rollbacked map[string]Operation
	Data       map[string]map[string]interface{}
}

type updater = func(state State)

func getStateKey(id int) string {
	return fmt.Sprintf("worflow:state:%v", id)
}

func createState(id int, cache cache.Cache) (State, error) {
	ctx := context.Background()

	state := State{
		ID:         id,
		Done:       make(map[string]Operation),
		Rollbacked: make(map[string]Operation),
		Data:       make(map[string]map[string]interface{}),
	}

	key := getStateKey(id)
	rawState, err := json.Marshal(state)
	if err != nil {
		return state, err
	}

	cache.Set(ctx, key, string(rawState))

	return state, nil
}

func updateState(id int, cache cache.Cache, update updater) (State, error) {
	var state State
	ctx := context.Background()

	key := getStateKey(id)
	rawState, err := cache.Get(ctx, key)

	if err != nil {
		return state, err
	}

	json.Unmarshal([]byte(rawState), &state)
	if err != nil {
		return state, err
	}

	update(state)

	rawNextState, err := json.Marshal(state)
	if err != nil {
		return state, err
	}

	cache.Set(ctx, key, string(rawNextState))
	return state, nil
}
