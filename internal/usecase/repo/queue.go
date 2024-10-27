package repo

import (
	"sync"

	entity "github.com/igor35hh/scheduler/internal/entity"
)

type Queue interface {
	Add(entity.Task)
	Pop() entity.Task
	Delete(string) bool
	Len() int
}

type TaskQueue struct {
	mu         sync.Mutex
	start, end *node
	lenght     int
}

type node struct {
	value entity.Task
	next  *node
}

func NewQueue() *TaskQueue {
	return &TaskQueue{start: nil, end: nil, lenght: 0}
}

func (q *TaskQueue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.lenght
}

// Add method adds task to the queue
func (q *TaskQueue) Add(task entity.Task) {
	q.mu.Lock()
	defer q.mu.Unlock()
	n := &node{value: task, next: nil}
	if q.lenght == 0 {
		q.start = n
		q.end = n
	} else {
		q.end.next = n
		q.end = n
	}
	q.lenght++
}

// Pop method retrieves task from the queue
func (q *TaskQueue) Pop() entity.Task {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.lenght == 0 {
		return nil
	}
	n := q.start
	if q.lenght == 1 {
		q.start = nil
		q.end = nil
	} else {
		q.start = q.start.next
	}
	q.lenght--
	return n.value
}

// Delete method removes task from the queue
func (q *TaskQueue) Delete(id string) bool {
	if q.lenght == 0 {
		return false
	}

	if q.start.value.GetId() == id {
		q.start = q.start.next
		return true
	}

	current := q.start
	for current.next != nil {
		if current.next.value.GetId() == id {
			current.next = current.next.next
			q.lenght--
			return true
		}
		current = current.next
	}

	return false
}
