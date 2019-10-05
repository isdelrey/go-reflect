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
