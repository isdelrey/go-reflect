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
	stream  network.Stream
}

type Streams struct {
	mutex         sync.Mutex
	List          []stream
	Subscriptions *Subscriptions
}

func (s *Streams) New(networkStream network.Stream) {
	newStream := stream{
		Channel: make(chan event.Event),
		stream:  networkStream,
	}
	s.List = append(s.List, newStream)

	go newStream.writeStream()
	go newStream.readStream(s.Subscriptions)
}

func (s *stream) writeStream() {
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

		fmt.Printf("Sent %d bytes: %x\n", len(raw), raw)
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
			continue
		}

		raw, err := s.readValue()
		if err != nil {
			continue
		}

		fmt.Printf("Received %d bytes: %x\n", len(raw), raw)

		fmt.Println("Received", string(name))

		subscriptions.Handle(string(name), raw)
	}
}

func (s *stream) streamLost() {
	fmt.Println(s.stream.Conn().RemotePeer().Pretty()[0:5], "disappeared")
}
