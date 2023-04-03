package syncmap

import "context"

// Syncmap is a thread safe generic map
type Syncmap[K comparable, V any] struct {
	values map[K]V
	rc     chan readRequest[K, V]
	wc     chan writeRequest[K, V]
	dc     chan K
	fc     chan chan map[K]V
	sc     chan chan int
}

type readRequest[K comparable, V any] struct {
	Key          K
	ResponseChan chan V
}

type writeRequest[K comparable, V any] struct {
	Key   K
	Value V
}

// NewSyncmap accepts a context and map, returns a thread safe map initalized with the content of the accepted map
// NewSyncmap spawns a go routine which will return when the context is cancelled
func NewSyncmap[K comparable, V any](ctx context.Context, m map[K]V) *Syncmap[K, V] {
	values := m

	if values == nil {
		values = make(map[K]V)
	}
	s := &Syncmap[K, V]{
		values: values,
		rc:     make(chan readRequest[K, V]),
		wc:     make(chan writeRequest[K, V]),
		dc:     make(chan K),
		fc:     make(chan chan map[K]V),
		sc:     make(chan chan int),
	}
	go s.serveRequests(ctx)

	return s
}

func (s *Syncmap[K, V]) serveRequests(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case readRequest := <-s.rc:
			readRequest.ResponseChan <- s.values[readRequest.Key]
		case writeRequest := <-s.wc:
			s.values[writeRequest.Key] = writeRequest.Value
		case k := <-s.dc:
			delete(s.values, k)
		case f := <-s.fc:
			c := make(map[K]V, len(s.values))
			for k, v := range s.values {
				c[k] = v
			}
			f <- c
		case req := <-s.sc:
			req <- len(s.values)
		}
	}
}

// Get accepts a key and thread safely returns the associated value from the map
func (s *Syncmap[K, V]) Get(k K) V {
	v := make(chan V)
	s.rc <- readRequest[K, V]{k, v}
	return <-v
}

// Set accepts a key and value and thread safely assigns the value to the key in the map
func (s *Syncmap[K, V]) Set(k K, v V) {
	s.wc <- writeRequest[K, V]{k, v}
}

// Delete accepts a key and thread safely removes the entry from the map
func (s *Syncmap[K, V]) Delete(k K) {
	s.dc <- k
}

// Flush thread safely returns a copy of the map
func (s *Syncmap[K, V]) Flush() map[K]V {
	v := make(chan map[K]V)
	s.fc <- v
	return <-v
}

// Size thread safely returns the size of the map
func (s *Syncmap[K, V]) Size() int {
	v := make(chan int)
	s.sc <- v
	return <-v
}
