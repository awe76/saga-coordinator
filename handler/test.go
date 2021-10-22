package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/awe76/saga-coordinator/cache"
	"github.com/go-redis/redis/v8"
)

type State struct {
	Count int
}

func updateState(cache cache.Cache) State {
	var state State

	ctx := context.Background()
	stateStr, err := cache.Get(ctx, "state")

	if err == redis.Nil {
		state = State{
			Count: 0,
		}
	} else if err != nil {
		panic(err)
	} else {
		json.Unmarshal([]byte(stateStr), &state)
	}

	state.Count++
	nextState, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}

	cache.Set(ctx, "state", string(nextState))

	return state
}

func HandleComment(msg *sarama.ConsumerMessage, cache cache.Cache) {
	state := updateState(cache)
	fmt.Printf("Received message Count %d: | Topic(%s) | Message(%s) \n", state.Count, string(msg.Topic), string(msg.Value))
}

func HandleError(err error) {
	fmt.Println(err)
}
