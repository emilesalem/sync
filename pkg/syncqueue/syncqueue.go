package syncqueue

import (
	"context"
	"time"
)

// Syncqueue is a queue wrapper which enables synced operations on the underlying queue
type Syncqueue[V any] struct {
	getHeadChan       chan chan V
	setHeadChan       chan setHeadRequest[V]
	flushChan         chan chan V
	sizeChan          chan chan int
	queueCapacity     int
	timeBeforeRequeue time.Duration
}

type setHeadRequest[V any] struct {
	msg   V
	rChan chan bool
}

type Options struct {
	Capacity int
}

// NewSyncqueue returns a queue wrapper which enables synced operations on the underlying queue
func NewSyncqueue[V any](ctx context.Context, opts Options) *Syncqueue[V] {

	s := &Syncqueue[V]{
		getHeadChan:   make(chan chan V),
		setHeadChan:   make(chan setHeadRequest[V]),
		flushChan:     make(chan chan V),
		sizeChan:      make(chan chan int),
		queueCapacity: opts.Capacity,
	}
	go s.serveQueueRequests(ctx)
	return s
}

// Read returns a channel in which the queue head will be written once the queue starts filling up;
func (s *Syncqueue[V]) Read() chan V {
	r := make(chan V)
	s.getHeadChan <- r
	return r
}

// Add appends an item to the queue; if queue capacity is reached, the queue head is popped before appending the new item;
// it returns true if queue capacity is not exceeded, false otherwise
func (s *Syncqueue[V]) Add(m V) bool {
	r := make(chan bool)
	s.setHeadChan <- setHeadRequest[V]{
		msg:   m,
		rChan: r,
	}
	return <-r
}

// Flush returns a channel to which all remaining queue messages will be written
func (s *Syncqueue[V]) Flush() chan V {
	r := make(chan V)
	s.flushChan <- r
	return r
}

// Size returns the number of messages left in the queue
func (s *Syncqueue[V]) Size() int {
	r := make(chan int)
	s.sizeChan <- r
	return <-r
}

func (s *Syncqueue[V]) serveQueueRequests(ctx context.Context) {
	queue := make([]V, 0, 1000)
	blockedGets := make([]chan V, 0, 1000)
	for {
		select {
		case <-ctx.Done():
			return
		case r := <-s.getHeadChan:
			if len(queue) > 0 {
				msg := queue[0]
				queue = queue[1:]
				r <- msg
			} else {
				blockedGets = append(blockedGets, r)
			}
		case r := <-s.setHeadChan:
			ok := true
			if len(blockedGets) > 0 {
				c := blockedGets[0]
				blockedGets = blockedGets[1:]
				c <- r.msg
			} else {
				if len(queue) >= s.queueCapacity {
					queue = queue[1:]
					ok = false
				}
				queue = append(queue, r.msg)
			}
			r.rChan <- ok
		case r := <-s.flushChan:
			for _, x := range queue {
				r <- x
			}
			queue = queue[:0]
			close(r)
		case r := <-s.sizeChan:
			r <- len(queue)
		}
	}
}
