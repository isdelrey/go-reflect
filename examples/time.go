package examples

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ivosequeros/reflect/mesh"
	"github.com/ivosequeros/reflect/message"
)

func Run() {
	/* Create mesh: */
	m := mesh.New(mesh.Options{
		/* A secret key is exchanged when the connection is established to verify that the other peer can join the mesh */
		Key: "SECRET_KEY",
	})

	/* Subscribe to a test event: */
	m.Subscribe("test", func(message message.Message) {
		fmt.Println(float64(time.Now().UnixNano()-message["start"].(int64))/1000, "µs", "("+strconv.Itoa(int(1000000000/int64(time.Now().UnixNano()-message["start"].(int64))))+"/s)")

	})

	for {
		/* Broadcast test event: */
		m.Broadcast("test", message.Message{
			"start": time.Now().UnixNano(),
		})
		time.Sleep(time.Second * 1)
	}

	/* Keep app running: */
	select {}
}
