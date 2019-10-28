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

func New() (ctr containerd.Container, err error) {
	log.Println("Welcome!")
	log.Println("Launching a quake3 server in a container!")
	defer log.Println("Done!")

	// Establish a connection with the daemon
	ctx := namespaces.WithNamespace(context.Background(), "default")

	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		log.Println(err)
		return
	}
	defer client.Close()

	// Load the base image
	image, err := client.GetImage(ctx, "quake3s")
	if err != nil {
		log.Println(image)
		return
	}

	// Create the container
	ctr, err = client.NewContainer(ctx, "quake3s",
		containerd.WithNewSnapshot("quake3s-snapshot", image),
		containerd.WithNewSpec(oci.WithImageConfig(image)))
	if err != nil {
		log.Println(err)
		return
	}

	// Create a running task
	task, err := ctr.NewTask(ctx, cio.NewCreator(cio.WithStdio))
	if err != nil {
		log.Println(err)
		return
	}

	// make sure we wait before calling start
    _, err = task.Wait(ctx)
	if err != nil {
		fmt.Println(err)
        return
	}

	// call start on the task to execute the quake server
	if err = task.Start(ctx); err != nil {
		log.Println(err)
		return
	}

    return
}
