# Kraken

Kraken is an orchestration system for running tasks across a cluster of worker
machines. A central manager accepts tasks from users, schedules them onto
workers, and tracks their state. Workers run the tasks and report metrics back
to the manager.

## Layout

```
.
├── main.go            # entry point
├── manager/           # accepts tasks, schedules them, tracks workers
├── scheduler/         # feasibility, scoring, picking
├── worker/            # runs tasks, exposes API and metrics
├── node/              # representation of a physical/virtual machine
└── task/              # task model and lifecycle
```

## Documentation

- [Architecture](docs/architecture.md) — components and how they interact
- [Tasks](docs/tasks.md) — task model and lifecycle
