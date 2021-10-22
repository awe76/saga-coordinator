package workflow

import (
	"encoding/json"

	"github.com/awe76/saga-coordinator/producer"
)

func pushData(producer producer.Producer, topic string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	producer.Push(topic, payload)
	return nil
}
