package peer

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

type Peer struct {
	ctx           context.Context
	host          host.Host
	peerChan      chan peer.AddrInfo
	Subscriptions *Subscriptions
	Streams       *Streams
	key           string
}

type Options struct {
	Key string
}

func AutoNewPeer(options Options) *Peer {
	// Select available port
	port := 8000
	for {
		_, err := net.Dial("tcp", ":"+strconv.Itoa(port))
		if err != nil {
			break
		}
		port++
	}

	peer := NewPeer(port, options.Key)
	fmt.Println("Running on port", port, "as", peer.host.ID().Pretty()[0:5])

	peer.peerChan = mdns.Initiate(peer.ctx, peer.host, "_host-discovery")
	go peer.awaitPeers()

	return peer
}

func (p *Peer) awaitPeers() {
	for {
		peer := <-p.peerChan // will block until we discover a peer
		peerstore := p.host.Peerstore()

		present := false
		for _, a := range peerstore.Peers() {
			if a == peer.ID {
				present = true
			}
		}

		if present {
			return
		}

		errConnect := p.host.Connect(p.ctx, peer)
		if errConnect != nil {
			fmt.Println(fmt.Sprintf("Error when connecting peers: %v", errConnect))
			return
		}

		stream, errStream := p.host.NewStream(p.ctx, peer.ID, protocol.ID("/reflect/1.0.0"))
		if errStream != nil {
			fmt.Println("Stream open failed", errStream)
		}

		go p.newOutgoingStream(stream)
	}
}

func NewPeer(port int, key string) *Peer {
	ctx, _ := context.WithCancel(context.Background())
	hma, _ := ma.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port))

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

	var subscriptions Subscriptions
	var streams Streams
	streams.Subscriptions = &subscriptions

	peer := &Peer{
		host:          host,
		ctx:           ctx,
		peerChan:      make(chan peer.AddrInfo),
		Subscriptions: &subscriptions,
		Streams:       &streams,
		key:           key,
	}

	peer.host.SetStreamHandler(protocol.ID("/reflect/1.0.0"), peer.newIncomingStream)

	return peer
}

func (p *Peer) newOutgoingStream(stream network.Stream) {
	stream.Write([]byte(p.key))

	fmt.Println("Sent key ->", stream.Conn().RemotePeer().Pretty()[0:5])

	p.Streams.New(stream)
}

func (p *Peer) newIncomingStream(stream network.Stream) {
	buffer := make([]byte, len(p.key))
	stream.Read(buffer)

	if string(buffer) != p.key {
		stream.Close()
	}

	fmt.Println("Received key <-", stream.Conn().RemotePeer().Pretty()[0:5])

	p.Streams.New(stream)
}
