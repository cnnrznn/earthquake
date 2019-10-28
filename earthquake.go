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

	ctx := namespaces.WithNamespace(context.Background(), "default")

	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		log.Println(err)
		return
	}
	defer client.Close()

	image, err := client.GetImage(ctx, "quake3s")
	if err != nil {
		log.Println(image)
		return
	}

	container, err := client.NewContainer(ctx, "quake3s",
		containerd.WithNewSnapshot("quake3s-snapshot", image),
		containerd.WithNewSpec(oci.WithImageConfig(image)))
	if err != nil {
		log.Println(err)
		return
	}
	defer container.Delete(ctx, containerd.WithSnapshotCleanup)

	task, err := container.NewTask(ctx, cio.NullIO)
	if err != nil {
		log.Println(err)
		return
	}

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

	// sleep for a lil bit to see the logs
	time.Sleep(3 * time.Second)

	ckpt, err := task.Checkpoint(ctx, containerd.WithCheckpointName("quake3s-init"))
	if err != nil {
		log.Println(err)
		return
	}
	defer client.ImageService().Delete(ctx, "quake3s-init")

	// kill the process and get the exit status
	if err := task.Kill(ctx, syscall.SIGTERM); err != nil {
		log.Println(err)
		return
	}

	status := <-exitStatusC
	code, _, err := status.Result()
	if err != nil {
		log.Println(err)
		return
	}

	task.Delete(ctx)

	fmt.Printf("quake3 server exited with status: %d\n", code)

	task, err = container.NewTask(ctx, cio.NewCreator(cio.WithStdio),
		containerd.WithTaskCheckpoint(ckpt))
	if err != nil {
		log.Println(err)
		return
	}
	defer task.Delete(ctx)

	// make sure we wait before calling start
	exitStatusC, err = task.Wait(ctx)
	if err != nil {
		fmt.Println(err)
	}

	// call start on the task to execute the redis server
	if err := task.Start(ctx); err != nil {
		log.Println(err)
		return
	}

	status = <-exitStatusC
	code, _, err = status.Result()
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("Checkpoint exited with", code)
}
