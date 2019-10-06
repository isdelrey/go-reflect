package store

import (
	"fmt"

	"github.com/ivosequeros/reflect/mesh"
	"github.com/ivosequeros/reflect/message"
)

type Store struct {
	Mesh *mesh.Mesh
	keys map[string]string
}

type Options struct {
	Key string
}

const (
	Empty  = "store.Empty"
	All    = "store.All"
	Set    = "store.Set"
	Get    = "store.Get"
	Delete = "store.Delete"
)

func New(opts Options) *Store {
	mesh := mesh.New(mesh.Options{
		Key: opts.Key,
	})

	store := Store{
		Mesh: mesh,
		keys: make(map[string]string),
	}

	go store.sync()

	return &store
}

func (s *Store) sync() {
	s.Mesh.Subscribe(Set, func(m message.Message) {
		fmt.Println("s")
		s.keys[m["k"].(string)] = m["v"].(string)
	})
	s.Mesh.Subscribe(Delete, func(m message.Message) {
		delete(s.keys, m["k"].(string))
	})
	s.Mesh.Subscribe(All, func(m message.Message) {
		for k, v := range m {
			s.keys[k] = v.(string)
		}
	})

	s.Mesh.Peer.Streams.OnStart = func(PeerId string) {
		m := message.Message{}
		for k, v := range s.keys {
			m[k] = v
		}

		s.Mesh.Push(PeerId, All, m)
	}
}

func (s *Store) Get(key string) string {
	return s.keys[key]
}

func (s *Store) Set(key string, value string) {
	s.keys[key] = value
	s.Mesh.Broadcast(Set, message.Message{
		"k": key,
		"v": value,
	})
}

func (s *Store) Delete(key string, value string) {
	delete(s.keys, key)
	s.Mesh.Broadcast(Delete, message.Message{
		"k": key,
	})
}
