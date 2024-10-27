# scheduler: A Golang Job Scheduling Package

The scheduler is a task scheduling package that lets you run Go functions in FIFO order.
It is designed to handle high-volume processing with manageable parameters.

## Quick Start
```
go get github.com/igor35hh/scheduler
```

```golang
For instance, you can use the following in your application that checks urls:

  type Job struct {
    name   string
    status int
    url    string
    err    error
  }

  func (j *Job) Run() (interface{}, error) {
    resp, err := http.Get(j.url)
    j.status = resp.StatusCode
    j.err = err
    fmt.Println(j)
    return j, err
  }

  sched := scheduler.NewScheduler(&scheduler.Parameters{
    Ctx:              context.Background(),
    TasksBuffer:      6,
    CountTasksToPick: 2,
    Log:              scheduler.NewLogger(2),
  })

  j := Job{name: "check wiki", url: "https://www.wikipedia.org/"}
  sched.Schedule(j.Run)

  for sched.PendingCount() != 0 || sched.RunningCount() != 0 {
    time.Sleep(1 * time.Second)
  }

  task := sched.GetReady()
  if t, ok := task.(*Job); ok {
    fmt.Println(t.name, t.status, t.url, t.err)
  }
```

## Examples

- [Examples directory](example)

## Concepts

- **Scheduler**: The scheduler is the interface for the service interaction. It provides a way to initialise the service, 
  to add the new task, and get the task after execution.
- **Service**: The service provides a pool of workers to execute tasks concurrently.
- **Queue**: The queue is the linked list structure that provides temporary storage for tasks.
- **Task**: The task wraps a "function", which will be added to the queue and then run by the service.

### Logging
Logs can be enabled.
The Logger interface can be implemented with your desired logging library.
The provided NewLogger uses the standard library's log package.