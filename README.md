### 🐙 Reflect for Go

**Reflect is an easy-to-use encrypted p2p pub/sub system with peer auto-discovery.**

#### Use case

You have a group of Go apps or a group of instances of a Go app and you want them to be able to comunicate with each other. You don't want to spend money on a Redis instance and you're afraid if you used a centralised system, it would have to be monitored and it could go down.

Reflect solves this by creating a mesh of encrypted TCP connections between your instances/apps. Every time an app wants to broadcast a message, reflect sends that message to all the other apps running in the same network. Each instance keeps an open connection to each other, so communication is very fast and does not rely on third-parties.

#### How do I use it?

Here's an example of use:
```go
package main

import (
	"fmt"
	"time"

	"github.com/ivosequeros/reflect/mesh"
)

type message map[string]interface{}

func main() {
	/* Create mesh: */
	m := mesh.New(mesh.Options{
		/* A secret key is exchanged when the connection is established to verify that the other peer can join the mesh */
		Key: "SECRET_KEY",
	})

	/* Subscribe to a test event: */
	m.Subscribe("test", func(message map[string]interface{}) {
		fmt.Println("Content:", message["value"])
	})

	time.Sleep(1 * time.Second)

	/* Broadcast test event: */
	m.Broadcast("test", message{
		"value": "hello",
	})

	/* Keep app running: */
	select {}
}


```

It's as easy as to call `mesh.New`. This function returns a mesh object with 4 properties:

- `Subscribe` binds a handler to an event
- `Broadcast` sends a message to all the instances in the same local network
- `SubscriptionChannel` returns a channel that receives events with the specified name
- `BroadcastChannel` returns a channel that sends all messages it received under the event name specified


#### License

Apache License 2.0

### Author
[Ivo Sequeros](https://github.com/ivosequeros)
