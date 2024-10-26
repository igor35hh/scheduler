package src

import (
	"time"

	"github.com/google/uuid"
)

type TaskInterface interface {
	Complete()
}

type Task struct {
	log      Logger
	id       string
	attempts int
	fn       func() (interface{}, error)
	object   interface{}
}

// NewTask returns an instance of new task
func NewTask(log Logger, fn func() (interface{}, error)) *Task {
	return &Task{
		log:      log,
		id:       uuid.NewString(),
		attempts: 3,
		fn:       fn,
	}
}

// Complete executes task, if task fails, it will retry for 3 times
func (t *Task) Complete() {
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
