package task

import (
	"context"
	"io"
	"log"
	"math"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
)

type State int

const (
	Pending State = iota
	Scheduled
	Running
	Completed
	Failed
)

type Task struct {
	ID            uuid.UUID
	ContainerID   string
	Name          string
	State         State
	Image         string
	Cpu           float64
	Memory        int
	Disk          int
	ExposedPorts  network.PortSet
	PortBindings  map[string]string
	RestartPolicy string
	StartTime     time.Time
	EndTime       time.Time
}

type TaskEvent struct {
	ID        uuid.UUID
	State     State
	Timestamp time.Time
	Task      Task
}

type Config struct {
	Name          string
	Image         string
	Cpu           float64
	Memory        int64
	Disk          int64
	AttachStdin   bool // Whether to attach to stdin
	AttachStdout  bool // Whether to attach to stdout
	AttachStderr  bool // Whether to attach to stderr
	ExposedPorts  network.PortSet
	Cmd           []string // Command to run in the container
	Env           []string // Environment variables for the container
	RestartPolicy string   // Restart policy for the container
	ContainerID   string   // ID of the running container, set after creation
}

type Docker struct {
	Client *client.Client // Docker client for interacting with the Docker API
	Config Config
}

type DockerResult struct {
	Error       error
	Action      string // "create", "start", "stop", "remove"
	ContainerId string // ID of the container associated with the task
	Result      string // Output from Docker API
}

func (d *Docker) Run() DockerResult {
	ctx := context.Background()
	// Pull the image if it doesn't exist locally
	reader, err := d.Client.ImagePull(ctx, d.Config.Image, client.ImagePullOptions{})
	if err != nil {
		log.Printf("Error pulling image %s: %v\n", d.Config.Image, err)
		return DockerResult{Error: err}
	}
	io.Copy(os.Stdout, reader)

	// Create the container
	createOpts := client.ContainerCreateOptions{
		Name: d.Config.Name,
		Config: &container.Config{
			Image:        d.Config.Image,
			Tty:          false,
			Env:          d.Config.Env,
			ExposedPorts: d.Config.ExposedPorts,
		},
		HostConfig: &container.HostConfig{
			RestartPolicy: container.RestartPolicy{Name: container.RestartPolicyMode(d.Config.RestartPolicy)},
			Resources: container.Resources{
				Memory:   d.Config.Memory,
				NanoCPUs: int64(d.Config.Cpu * math.Pow(10, 9)), // Convert CPU cores to NanoCPUs
			},
			PublishAllPorts: true,
		},
	}
	resp, err := d.Client.ContainerCreate(ctx, createOpts)
	if err != nil {
		log.Printf("Error creating container using image %s: %v\n",
			d.Config.Image, err)
		return DockerResult{Error: err}
	}

	// Start the container
	_, err = d.Client.ContainerStart(ctx, resp.ID, client.ContainerStartOptions{})
	if err != nil {
		log.Printf("Error starting container %s: %v\n", resp.ID, err)
		return DockerResult{Error: err}
	}

	d.Config.ContainerID = resp.ID

	out, err := d.Client.ContainerLogs(ctx, resp.ID, client.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		log.Printf("Error getting logs for container %s: %v\n", resp.ID, err)
		return DockerResult{Error: err}
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	return DockerResult{
		ContainerId: resp.ID,
		Action:      "start",
		Result:      "success",
	}
}

func (d *Docker) Stop(id string) DockerResult {
	log.Printf("Attempting to stop container %v", id)
	ctx := context.Background()
	_, err := d.Client.ContainerStop(ctx, id, client.ContainerStopOptions{})
	if err != nil {
		log.Printf("Error stopping container %s: %v\n", id, err)
		return DockerResult{Error: err}
	}
	_, err = d.Client.ContainerRemove(ctx, id, client.ContainerRemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   false,
		Force:         false,
	})
	if err != nil {
		log.Printf("Error removing container %s: %v\n", id, err)
		return DockerResult{Error: err}
	}
	return DockerResult{
		ContainerId: id,
		Action:      "stop",
		Result:      "success",
	}
}
