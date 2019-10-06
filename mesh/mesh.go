package mesh

import (
	"errors"
	"time"

	"github.com/ivosequeros/reflect/event"
	"github.com/ivosequeros/reflect/message"
	"github.com/ivosequeros/reflect/peer"
)

type Mesh struct {
	Peer *peer.Peer
}

type Options struct {
	Key string
}

func New(opts Options) *Mesh {
	peer := peer.AutoNewPeer(peer.Options{
		Key: opts.Key,
	})

	return &Mesh{
		peer,
	}
}

func (m *Mesh) Broadcast(Name string, Message message.Message) {
	for _, stream := range m.Peer.Streams.List {
		stream.Channel <- event.Event{
			Name,
			Message,
		}
	}
}

func (m *Mesh) Push(PeerId string, Name string, Message message.Message) error {
	for _, stream := range m.Peer.Streams.List {
		if stream.PeerId == PeerId {
			stream.Channel <- event.Event{
				Name,
				Message,
			}
			break
		}
		return errors.New("peer not found")
	}
	return nil
}

func (m *Mesh) Subscribe(name string, handler func(message message.Message)) time.Time {
	return m.Peer.Subscriptions.Add(name, handler)
}

func (m *Mesh) Unsubscribe(id time.Time) {
	m.Peer.Subscriptions.Remove(id)
}
