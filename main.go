package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/ivosequeros/reflect/store"
)

func main() {
	/* Create mesh: */
	s := store.New(store.Options{
		/* A secret key is exchanged when the connection is established to verify that the other peer can join the mesh */
		Key: "SECRET_KEY",
	})

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	n := strconv.Itoa(r.Intn(100))
	s.Set("Number", n)
	fmt.Println("Generated", n)

	for {
		fmt.Println("Read", s.Get("Number"))
		time.Sleep(time.Second * 2)
	}

	/* Keep app running: */
	select {}
}
