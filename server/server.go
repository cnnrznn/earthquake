package server

import (
	"context"
	"fmt"
	"log"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
)

type Server struct {
	ctr         containerd.Container
	checkpoints []containerd.Image
	ctx         context.Context
	client      *containerd.Client
}

func getClient() (client *containerd.Client, err error) {
	client, err = containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		log.Println(err)
		return
	}

	return
}

func (s *Server) Checkpoint() (err error) {
	name := fmt.Sprintf("quake3s-c%v", len(s.checkpoints)+1)
	ckpt, err := s.ctr.Checkpoint(s.ctx, name,
		containerd.WithCheckpointTask,
		containerd.WithCheckpointRuntime,
		containerd.WithCheckpointRW)
	if err != nil {
		return
	}

	s.checkpoints = append(s.checkpoints, ckpt)

	return
}

func FromCkpt() (s Server, err error) {
	return
}

func New() (s Server, err error) {
	s = Server{}

	log.Println("Welcome!")
	log.Println("Launching a quake3 server in a container!")
	defer log.Println("Done!")

	// Establish a connection with the daemon
	ctx := namespaces.WithNamespace(context.Background(), "default")
	s.ctx = ctx

	client, err := getClient()
	if err != nil {
		return
	}
	s.client = client

	// Load the base image
	image, err := client.GetImage(ctx, "quake3s")
	if err != nil {
		return
	}

	// Create the container
	ctr, err := client.NewContainer(ctx, "quake3s",
		containerd.WithNewSnapshot("quake3s-snapshot", image),
		containerd.WithNewSpec(oci.WithImageConfig(image)))
	if err != nil {
		return
	}

	s.ctr = ctr

	// Create a running task
	task, err := ctr.NewTask(ctx, cio.NewCreator(cio.WithStdio))
	if err != nil {
		return
	}

	// make sure we wait before calling start
	_, err = task.Wait(ctx)
	if err != nil {
		return
	}

	// call start on the task to execute the quake server
	if err = task.Start(ctx); err != nil {
		return
	}

	return
}
