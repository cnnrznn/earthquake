package main

import (
	"time"

	"github.com/cnnrznn/earthquake/server"
)

func main() {
	s, err := server.New()
	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		// wait a little bit
		time.Sleep(3 * time.Second)
		s.Checkpoint()
	}

	return
}
