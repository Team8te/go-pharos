package photonq

import (
	"container/heap"
)

type Queue struct {
	heap.Interface
}

func MakeProtonQ() *Queue {
	return &Queue{}
}

/*
func (q *Queue) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	tick := time.Duration(0)
	for {
		select {
		case <-ctx.Done():
			return nil

		case <-time.After(tick):

		}
	}
}
*/
