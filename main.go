package main

import (
	"fmt"
	"time"

	"github.com/ivosequeros/reflect/sync"
)

type message map[string]interface{}

func main() {
	go sync.Create()
	fmt.Println("Waiting for new messages")
	sync.Subscribe("test", func(message map[string]interface{}) {
		fmt.Println(message["value"])
	})
	time.Sleep(2 * time.Second)

	sync.Broadcast("test", message{
		"value": "hello",
	})
	select {}
}
