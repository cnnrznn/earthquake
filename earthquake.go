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

	// wait a little bit
	time.Sleep(10 * time.Second)
	err = s.CheckpointRestore()
	//err = s.Checkpoint()
	if err != nil {
		panic(err)
	}

	return
}
