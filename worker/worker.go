package worker

import (
	"RawanMostafa08/Kraken/task"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

type Worker struct {
	Name      string
	Queue     queue.Queue
	Db        map[uuid.UUID]*task.Task
	TaskCount int
}

// RunTask executes a task from the worker's queue
func (w *Worker) RunTask() {
}

// StartTask starts a task and updates its state to Running
func (w *Worker) StartTask() {
}

// StopTask stops a running task and updates its state to Completed or Failed
func (w *Worker) StopTask() {
}

// CollectStats collects and returns statistics about the worker
func (w *Worker) CollectStats() {
}
