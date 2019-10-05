package sync

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"net"
	"strconv"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	host "github.com/libp2p/go-libp2p-host"
	ma "github.com/multiformats/go-multiaddr"

	"github.com/ivosequeros/reflect/mdns"
)

func Create() {
	ctx, _ := context.WithCancel(context.Background())

	// Select available port
	port := 8000
	for {
		_, err := net.Dial("tcp", ":"+strconv.Itoa(port))
		if err != nil {
			break
		}
		port++
	}

	host := new(ctx, port)
	fmt.Println("Running on port", port, "as", host.ID().Pretty()[0:5])
	peerChan := mdns.Initiate(ctx, host, "_host-discovery")

	go awaitPeers(ctx, host, peerChan)
}

func awaitPeers(ctx context.Context, host host.Host, peerChan chan peer.AddrInfo) {
	for {
		peer := <-peerChan // will block until we discover a peer
		peerstore := host.Peerstore()

		present := false
		for _, a := range peerstore.Peers() {
			if a == peer.ID {
				present = true
			}
		}

		if present {
			return
		}

		errConnect := host.Connect(ctx, peer)
		if errConnect != nil {
			fmt.Println(fmt.Sprintf("Error when connecting peers: %v", errConnect))
			return
		}

		stream, errStream := host.NewStream(ctx, peer.ID, protocol.ID("/reflect/1.0.0"))
		if errStream != nil {
			fmt.Println("Stream open failed", errStream)
		}

		go outgoing(stream)
	}
}

func new(ctx context.Context, p int) host.Host {
	hma, _ := ma.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", p))

	// Generate a key
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
	if err != nil {
		panic(err)
	}

	// Create host
	host, err := libp2p.New(ctx, libp2p.ListenAddrs(hma), libp2p.Identity(prvKey))
	if err != nil {
		log.Fatal(err)
	}
	host.SetStreamHandler(protocol.ID("/reflect/1.0.0"), incoming)

	return host
}

func outgoing(stream network.Stream) {
	stream.Write([]byte("hello"))

	fmt.Println("Said hello to", stream.Conn().RemotePeer().Pretty()[0:5])

	Streams.New(stream)
}

func incoming(stream network.Stream) {
	buffer := make([]byte, 5)
	stream.Read(buffer)

	if string(buffer) != "hello" {
		stream.Close()
	}

	fmt.Println(stream.Conn().RemotePeer().Pretty()[0:5] + " said hello")

	Streams.New(stream)
}
