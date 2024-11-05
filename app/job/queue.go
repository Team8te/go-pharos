package job

import (
	"container/list"
	"context"
	"sync"
)

var index int = 0

type task struct {
	id int
	j  *Job
}

type JobQueue struct {
	mx       sync.RWMutex
	queue    list.List
	notifier chan struct{}
}

func NewQueue() *JobQueue {
	return &JobQueue{
		notifier: make(chan struct{}, 1000),
	}
}

func (q *JobQueue) Push(w Worker) int {
	t := &task{
		j: NewJob(0, w),
	}

	q.notifier <- struct{}{}

	q.mx.Lock()
	defer q.mx.Unlock()

	t.id = index
	index++
	q.queue.PushBack(t)

	return t.id
}

func (q *JobQueue) Init() error {
	return nil
}

func (q *JobQueue) Work(ctx context.Context) error {
	for _ = range q.notifier {
	}

	t := q.pop()
	if t == nil {
		return nil
	}

	return t.j.Run(ctx)
}

func (q *JobQueue) pop() *task {
	q.mx.Lock()
	defer q.mx.Unlock()

	el := q.queue.Front()
	if el != nil {
		return el.Value.(*task)
	}

	return nil
}

func (q *JobQueue) Stop() {
	for e := q.queue.Front(); e != nil; e = e.Next() {
		j := e.Value.(*Job)
		j.Stop()
	}
}

func (q *JobQueue) Name() string {
	return "Job queue"
}
