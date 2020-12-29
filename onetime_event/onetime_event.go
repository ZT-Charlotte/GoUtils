package syncevent

import (
	"sync"
	"sync/atomic"
)

type Event struct {
	// atomic bool
	fired int32
	c     chan struct{}
	o     sync.Once
}

func NewEvent() *Event {
	return &Event{c: make(chan struct{})}
}

/*
 *@func: caused e to complete, safe to call multiple times and concurrently
 *@ret: returns true if this call caused the signaling channel returned by Done to close
 */
func (e *Event) Fire() bool {
	ret := false

	e.o.Do(func() {
		atomic.StoreInt32(&e.fired, 1)
		close(e.c)
		ret = true
	})

	return ret
}

/*
 *@func: returns a channel that will be closed when Fire is called
 */
func (e *Event) Done() <-chan struct{} {
	return e.c
}

/*
 *@func: returns true if fire has been called
 */
func (e *Event) HasFired() bool {
	return (atomic.LoadInt32(&e.fired) == 1)
}
