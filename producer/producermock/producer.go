package producermock

import (
	"encoding/json"
	"fmt"
)

type ProducerMock struct {
	messages map[string][]string
}

func New() *ProducerMock {
	return &ProducerMock{
		messages: make(map[string][]string),
	}
}

func (p *ProducerMock) SendMessage(topic string, message interface{}) error {
	raw, err := json.Marshal(message)
	if err != nil {
		return err
	}

	messages, found := p.messages[topic]
	if found {
		p.messages[topic] = append(messages, string(raw))
	} else {
		p.messages[topic] = append([]string{}, string(raw))
	}

	return nil
}

func (p *ProducerMock) Has(topic string, message interface{}) bool {
	raw, err := json.Marshal(message)
	if err != nil {
		return false
	}

	pattern := string(raw)
	fmt.Printf("%v\n", pattern)
	if messages, found := p.messages[topic]; found {
		for _, message := range messages {
			fmt.Printf("%v\n", message)
			if message == pattern {
				return true
			}
		}
	}

	return false
}
