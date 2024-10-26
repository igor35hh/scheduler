package scheduler

import (
	"context"

	"github.com/igor35hh/scheduler/src"
)

type Scheduler interface {
	Schedule(func() (interface{}, error)) string
	Stop()
	GetReady() interface{}
	PendingCount() int
	ReadyCount() int
	RunningCount() int
}

type Parameters struct {
	Ctx              context.Context // The context is using to stop jobs
	TasksBuffer      int64           // The count of concurrently running tasks
	CountTasksToPick uint            // The count of task to pick up from the pending queue
	Log              src.Logger      // The logger instance to log steps of tasks execution
}

func NewScheduler(p *Parameters) Scheduler {
	return src.NewService(p.Ctx, p.TasksBuffer, p.CountTasksToPick, p.Log)
}

var (
	LogLevelError = src.LogLevelError
	LogLevelWarn  = src.LogLevelWarn
	LogLevelInfol = src.LogLevelInfo
	LogLevelDebug = src.LogLevelDebug
)

func NewLogger(level src.LogLevel) src.Logger {
	return src.NewLogger(level)
}
