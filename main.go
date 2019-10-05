package main

import (
	"fmt"
	"time"

	"github.com/ivosequeros/reflect/mesh"
)

type message map[string]interface{}

func main() {
	/* Create mesh: */
	go mesh.Create()

	/* Subscribe to a test event: */
	mesh.Subscribe("test", func(message map[string]interface{}) {
		fmt.Println(message["value"])
	})

	time.Sleep(1 * time.Second)

	/* Broadcast test event: */
	mesh.Broadcast("test", message{
		"value": "hello",
	})

	/* Keep app running: */
	select {}
}
