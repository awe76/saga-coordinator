package workflow

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/awe76/saga-coordinator/cache"
)

type state struct {
	ID         int
	IsRollback bool
	Completed  bool
	Done       map[string]Operation
	InProgress map[string]Operation
	Data       map[string]map[string]interface{}
}

func (s *state) getCacheKey() string {
	return fmt.Sprintf("workflow:state:%v", s.ID)
}

func (s *state) init(cache cache.Cache) error {
	ctx := context.Background()

	s.IsRollback = false
	s.Completed = false
	s.Done = make(map[string]Operation)
	s.InProgress = make(map[string]Operation)
	s.Data = make(map[string]map[string]interface{})

	key := s.getCacheKey()
	cache.Set(ctx, key, s)
	return nil
}

func (s *state) update(cache cache.Cache, update func(*state)) error {
	ctx := context.Background()

	key := s.getCacheKey()
	rawState, err := cache.Get(ctx, key)

	if err != nil {
		return err
	}

	json.Unmarshal([]byte(rawState), s)
	if err != nil {
		return err
	}

	update(s)
	err = cache.Set(ctx, key, s)
	return err
}
