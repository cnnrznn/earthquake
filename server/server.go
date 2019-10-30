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
	task        containerd.Task
	id          string
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
	//ckpt, err := s.ctr.Checkpoint(s.ctx, name,
	//        containerd.WithCheckpointImage,
	//	containerd.WithCheckpointTask,
	//	containerd.WithCheckpointRuntime,
	//	containerd.WithCheckpointRW)
	ckpt, err := s.task.Checkpoint(s.ctx,
		containerd.WithCheckpointName(name))
	if err != nil {
		return
	}

	s.checkpoints = append(s.checkpoints, ckpt)

	return
}

func (s *Server) CheckpointRestore() (err error) {
	log.Println("Restoring server in new container")
	checkpointName := "checkpoint"

	statusC, err := s.task.Wait(s.ctx)
	if err != nil {
		return
	}

	checkpoint, err := s.ctr.Checkpoint(s.ctx, checkpointName,
		containerd.WithCheckpointRuntime,
		containerd.WithCheckpointRW,
		containerd.WithCheckpointTask)
	if err != nil {
		return
	}

	<-statusC

	_, err = s.task.Delete(s.ctx)
	if err != nil {
		return
	}

	err = s.ctr.Delete(s.ctx, containerd.WithSnapshotCleanup)
	if err != nil {
		return
	}

	container, err := s.client.Restore(s.ctx, s.id, checkpoint,
		containerd.WithRestoreImage,
		containerd.WithRestoreSpec,
		containerd.WithRestoreRuntime,
		containerd.WithRestoreRW)
	if err != nil {
		return
	}
	s.ctr = container

	task, err := container.NewTask(s.ctx, cio.NewCreator(cio.WithStdio),
		containerd.WithTaskCheckpoint(checkpoint))
	if err != nil {
		return
	}
	s.task = task

	err = task.Start(s.ctx)
	if err != nil {
		return
	}

	return
}

func FromCkpt() (s Server, err error) {
	return
}

func New() (s Server, err error) {
	s = Server{}
	s.id = "quake3s"

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
	s.task = task

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
