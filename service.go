package src

import (
	"context"
	"sync/atomic"
	"time"
)

type Scheduler interface {
	Schedule(func() (interface{}, error)) string
	Stop()
	GetReady() interface{}
	PendingCount() int
	ReadyCount() int
	RunningCount() int
}

// Scheduler ...
type TaskScheduler struct {
	ctx               context.Context
	ctxCancel         context.CancelFunc
	pendingQueue      Queue
	readyQueue        Queue
	tasksBuffer       int64
	chanBuffer        chan struct{}
	countRunningTasks int64
	countTasksToPick  uint
	log               Logger
}

func NewService(
	ctx context.Context,
	tasksBuffer int64,
	countTasksToPick uint,
	log Logger,
) *TaskScheduler {
	ctxScheduler, cancel := context.WithCancel(ctx)
	sc := &TaskScheduler{
		ctx:              ctxScheduler,
		ctxCancel:        cancel,
		tasksBuffer:      tasksBuffer,
		chanBuffer:       make(chan struct{}, tasksBuffer),
		countTasksToPick: countTasksToPick,
		pendingQueue:     NewQueue(),
		readyQueue:       NewQueue(),
		log:              log,
	}

	go sc.start()

	return sc
}

func (s *TaskScheduler) GetReady() interface{} {
	task := s.readyQueue.Pop()
	return task.GetObject()
}

// Schedule method put task into the pending queue and return task id
func (s *TaskScheduler) Schedule(fn func() (interface{}, error)) string {
	task := NewTask(s.log, fn)
	s.pendingQueue.Add(task)
	s.log.Info("the task %s was added to pending queue", task.id)
	return task.id
}

// Stop method cancel context of service, it will stop all running jobs
func (s *TaskScheduler) Stop() {
	s.ctxCancel()
}

// PendingLenght method returns lenght of pending queue
func (s *TaskScheduler) PendingCount() int {
	return s.pendingQueue.Len()
}

// ReadyLenght method returns lenght of ready queue
func (s *TaskScheduler) ReadyCount() int {
	return s.readyQueue.Len()
}

func (s *TaskScheduler) RunningCount() int {
	return int(atomic.LoadInt64(&s.countRunningTasks))
}

func (s *TaskScheduler) get() Task {
	task := s.pendingQueue.Pop()
	if task != nil {
		s.log.Info("the task %s was taken into the proccess", task.GetId())
	}

	return task
}

func (s *TaskScheduler) runWorkers(pool uint) {
	for i := pool; i > 0; i-- {
		if task := s.get(); task != nil {
			select {
			case s.chanBuffer <- struct{}{}:
				atomic.AddInt64(&s.countRunningTasks, 1)
			case <-s.ctx.Done():
				return
			}

			go func() {
				select {
				case <-s.chanBuffer:
					task.Complete()
					s.readyQueue.Add(task)
					atomic.AddInt64(&s.countRunningTasks, -1)
				case <-s.ctx.Done():
					return
				}
			}()
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}

func (s *TaskScheduler) start() {
	defer s.ctxCancel()
	defer close(s.chanBuffer)
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			currentRunning := atomic.LoadInt64(&s.countRunningTasks)
			if currentRunning == 0 {
				s.runWorkers(uint(s.tasksBuffer))
			} else if currentRunning == (s.tasksBuffer - int64(s.countTasksToPick)) {
				s.runWorkers(s.countTasksToPick)
			}
		}
	}
}
