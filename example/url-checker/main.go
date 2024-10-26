package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/igor35hh/scheduler"
)

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

func main() {
	sched := scheduler.NewScheduler(&scheduler.Parameters{
		Ctx:              context.Background(),
		TasksBuffer:      6,
		CountTasksToPick: 2,
		Log:              scheduler.NewLogger(scheduler.LogLevelInfol),
	})

	for i := 0; i < 100; i++ {
		j := Job{name: "check wiki", url: "https://www.wikipedia.org/"}
		sched.Schedule(j.Run)
	}

	for sched.PendingCount() != 0 || sched.RunningCount() != 0 {
		time.Sleep(1 * time.Second)
	}

	fmt.Println(sched.PendingCount(), sched.ReadyCount())

	for sched.ReadyCount() != 0 {
		if task := sched.GetReady(); task != nil {
			if t, ok := task.(*Job); ok {
				fmt.Println(t.name, t.status, t.url, t.err)
			}
		}
	}

	sched.Stop()
}
