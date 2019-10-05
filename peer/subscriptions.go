package peer

import (
	"fmt"
	"sync"

	"github.com/vmihailenco/msgpack/v4"
)

type Subscriptions struct {
	mutex sync.Mutex
	list  []subscription
}

type subscription struct {
	eventName string
	handler   func(payload map[string]interface{})
}

func (s *Subscriptions) Add(eventName string, handler func(message map[string]interface{})) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.list = append(s.list, subscription{
		eventName,
		handler,
	})
}

func (s *Subscriptions) Handle(name string, raw []byte) {
	var message map[string]interface{}
	err := msgpack.Unmarshal(raw, &message)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, subscription := range s.list {
		if subscription.eventName == name {
			subscription.handler(message)
		}
	}
}
