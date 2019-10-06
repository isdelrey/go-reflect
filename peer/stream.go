package peer

import (
	"fmt"
	"sync"

	"github.com/ivosequeros/reflect/event"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/vmihailenco/msgpack/v4"
)

type stream struct {
	Channel chan event.Event
	PeerId  string
	stream  network.Stream
	OnEnd   *func(PeerId string)
}

type Streams struct {
	mutex         sync.Mutex
	List          []stream
	Subscriptions *Subscriptions
	OnStart       func(PeerId string)
	OnEnd         func(PeerId string)
}

func (s *Streams) New(networkStream network.Stream) {
	newStream := stream{
		PeerId:  networkStream.Conn().RemotePeer().Pretty(),
		Channel: make(chan event.Event),
		stream:  networkStream,
		OnEnd:   &s.OnEnd,
	}
	s.List = append(s.List, newStream)

	go newStream.writeStream()
	go newStream.readStream(s.Subscriptions)

	if s.OnStart != nil {
		s.OnStart(networkStream.Conn().RemotePeer().Pretty())
	}
	fmt.Println(networkStream.Conn().RemotePeer().Pretty()[0:5], "joined")
}

func (s *stream) writeStream() {
	defer s.streamLost()

	for {
		event := <-s.Channel
		raw, err := msgpack.Marshal(event.Message)

		if err != nil {
			fmt.Println(err)
			continue
		}

		s.stream.Write([]byte(event.Name))
		s.stream.Write([]byte{byte(0)})
		s.stream.Write(raw)
		s.stream.Write([]byte{byte(0)})
	}
}

func (s *stream) readValue() ([]byte, error) {
	buffer := make([]byte, 1)
	value := make([]byte, 0)
	for {
		_, err := s.stream.Read(buffer)

		if err != nil {
			return nil, err
		}

		if buffer[0] == 0 {
			break
		}

		value = append(value, buffer[0])
	}
	return value, nil
}

func (s *stream) readStream(subscriptions *Subscriptions) {
	defer s.streamLost()

	for {
		name, err := s.readValue()
		if err != nil {
			if err.Error() == "stream reset" {
				break
			}
			continue
		}

		raw, err := s.readValue()
		if err != nil {
			if err.Error() == "stream reset" {
				break
			}
			continue
		}

		fmt.Println("[" + string(name) + "]")
		subscriptions.Handle(string(name), raw)
	}
}

func (s *stream) streamLost() {
	if *s.OnEnd != nil {
		(*s.OnEnd)(s.stream.Conn().RemotePeer().Pretty())
	}
	fmt.Println(s.stream.Conn().RemotePeer().Pretty()[0:5], "left")
}
