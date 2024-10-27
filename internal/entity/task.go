package entity

import (
	"time"

	"github.com/google/uuid"

	pkg "github.com/igor35hh/scheduler/pkg"
)

type Task interface {
	Complete()
	GetId() string
	GetObject() interface{}
}

type TaskWrapper struct {
	log      pkg.Logger
	id       string
	attempts int
	fn       func() (interface{}, error)
	object   interface{}
}

// NewTask returns an instance of new task
func NewTask(log pkg.Logger, fn func() (interface{}, error)) *TaskWrapper {
	return &TaskWrapper{
		log:      log,
		id:       uuid.NewString(),
		attempts: 3,
		fn:       fn,
	}
}

// GetId returns an id of the task
func (t *TaskWrapper) GetId() string {
	return t.id
}

// GetObject returns object of the task
func (t *TaskWrapper) GetObject() interface{} {
	return t.object
}

// Complete executes task, if task fails, it will retry for 3 times
func (t *TaskWrapper) Complete() {
	for i := 0; i < t.attempts; i++ {
		t.attempts--
		j, err := t.fn()
		if err == nil {
			t.object = j
			t.log.Info("task %s executed succesfully", t.id)
			break
		}
		t.log.Warn("task %s attempt %d, executed with error %v", t.id, i, err)
		time.Sleep(time.Duration(1+i) * time.Minute)
	}
}
