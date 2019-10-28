package main

import (
	"context"
	"fmt"
	"log"
	"syscall"
	"time"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
)

func main() {
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
	container, err := client.NewContainer(ctx, "quake3s",
		containerd.WithNewSnapshot("quake3s-snapshot", image),
		containerd.WithNewSpec(oci.WithImageConfig(image)))
	if err != nil {
		log.Println(err)
		return
	}
	defer container.Delete(ctx, containerd.WithSnapshotCleanup)

	// Create a running task
	task, err := container.NewTask(ctx, cio.NewCreator(cio.WithStdio))
	if err != nil {
		log.Println(err)
		return
	}
	defer task.Delete(ctx)

	// make sure we wait before calling start
	exitStatusC, err := task.Wait(ctx)
	if err != nil {
		fmt.Println(err)
	}

	// call start on the task to execute the redis server
	if err := task.Start(ctx); err != nil {
		log.Println(err)
		return
	}

	time.Sleep(20 * time.Second)

	task.Kill(ctx, syscall.SIGTERM)

	status := <-exitStatusC
	code, _, err := status.Result()
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("Quake3 Server exited with", code)
}
