package scheduler

import (
	"context"

	src "github.com/igor35hh/scheduler/internal"
	pkg "github.com/igor35hh/scheduler/pkg"
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
	Log              pkg.Logger      // The logger instance to log steps of tasks execution
}

func NewScheduler(p *Parameters) Scheduler {
	return src.NewService(p.Ctx, p.TasksBuffer, p.CountTasksToPick, p.Log)
}

var (
	LogLevelError = pkg.LogLevelError
	LogLevelWarn  = pkg.LogLevelWarn
	LogLevelInfol = pkg.LogLevelInfo
	LogLevelDebug = pkg.LogLevelDebug
)

func NewLogger(level pkg.LogLevel) pkg.Logger {
	return pkg.NewLogger(level)
}
