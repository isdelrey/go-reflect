package mesh

import (
	"fmt"
	"sync"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/vmihailenco/msgpack/v4"
)

type stream struct {
	channel chan event
	stream  network.Stream
}

type streams struct {
	mutex sync.Mutex
	list  []stream
}

func (s *streams) New(networkStream network.Stream) {
	newStream := stream{
		channel: make(chan event),
		stream:  networkStream,
	}
	s.list = append(s.list, newStream)

	go newStream.writeStream()
	go newStream.readStream()
}

func (s *stream) writeStream() {
	for {
		event := <-s.channel
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

func (s *stream) readStream() {
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

		Subscriptions.Handle(string(name), raw)
	}
}

func (s *stream) streamLost() {
	fmt.Println(s.stream.Conn().RemotePeer().Pretty()[0:5], "disappeared")
}

var Streams streams
