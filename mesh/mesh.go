package mesh

import (
	"github.com/ivosequeros/reflect/event"
	"github.com/ivosequeros/reflect/peer"
)

type Mesh struct {
	Peer *peer.Peer
}

func New() *Mesh {
	peer := peer.AutoNewPeer()

	return &Mesh{
		peer,
	}
}

func (m *Mesh) Broadcast(Name string, Message interface{}) {
	for _, stream := range m.Peer.Streams.List {
		stream.Channel <- event.Event{
			Name,
			Message,
		}
	}
}

func (m *Mesh) Subscribe(name string, handler func(message map[string]interface{})) {
	m.Peer.Subscriptions.Add(name, handler)
}

func (m *Mesh) SubscriptionChannel(name string) chan map[string]interface{} {
	channel := make(chan map[string]interface{})

	m.Peer.Subscriptions.Add(name, func(message map[string]interface{}) {
		channel <- message
	})

	return channel
}

func (m *Mesh) BroadcastChannel(name string) chan map[string]interface{} {
	channel := make(chan map[string]interface{})

	go func() {
		for {
			message := <-channel
			m.Broadcast(name, message)
		}
	}()

	return channel
}
