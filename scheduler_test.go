package scheduler_test

import (
	"context"
	"testing"
	"time"

	"github.com/igor35hh/scheduler"
	"github.com/stretchr/testify/assert"
)

type Job struct {
	name   string
	status string
}

func (j *Job) Run() (interface{}, error) {
	j.status = j.name
	return j, nil
}

func TestSchedule(t *testing.T) {
	sched := scheduler.NewScheduler(&scheduler.Parameters{
		Ctx:              context.Background(),
		TasksBuffer:      6,
		CountTasksToPick: 2,
		Log:              scheduler.NewLogger(scheduler.LogLevelInfol),
	})

	for i := 0; i < 100; i++ {
		j := Job{name: "send email"}
		sched.Schedule(j.Run)
		k := Job{name: "check the link"}
		sched.Schedule(k.Run)
	}

	for sched.PendingCount() != 0 || sched.RunningCount() != 0 {
		time.Sleep(1 * time.Second)
	}

	sched.Stop()

	assert.Equal(t, sched.PendingCount(), 0)
	assert.Equal(t, sched.RunningCount(), 0)
	assert.Equal(t, sched.ReadyCount(), 200)

	for sched.ReadyCount() != 0 {
		if task := sched.GetReady(); task != nil {
			if ts, ok := task.(*Job); ok {
				assert.Equal(t, ts.status, ts.name)
			}
		}
	}

	assert.Equal(t, sched.ReadyCount(), 0)
}
