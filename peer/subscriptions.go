package peer

import (
	"fmt"
	"sync"
	"time"

	"github.com/ivosequeros/reflect/message"
	"github.com/vmihailenco/msgpack/v4"
)

type Subscriptions struct {
	mutex sync.Mutex
	list  []subscription
}

type subscription struct {
	id        time.Time
	eventName string
	handler   func(payload message.Message)
}

func (s *Subscriptions) Add(eventName string, handler func(message message.Message)) time.Time {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	id := time.Now()

	s.list = append(s.list, subscription{
		id,
		eventName,
		handler,
	})
	return id
}

func (s *Subscriptions) Remove(id time.Time) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	i := 0
	for j, e := range s.list {
		if e.id == id {
			i = j
			break
		}
	}

	s.list = append(s.list[:i], s.list[i+1:]...)
}

func (s *Subscriptions) Handle(name string, raw []byte) {
	var message message.Message
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
