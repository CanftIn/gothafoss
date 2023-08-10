package pool

import (
	"runtime"
	"sync"

	queue "github.com/CanftIn/gothafoss/lib/im/pool/internal"
)

type Queue struct {
	sync.Mutex
	popable *sync.Cond
	buffer  *queue.Queue
	closed  bool
}

func NewQueue() *Queue {
	e := &Queue{
		buffer: queue.New(),
	}
	e.popable = sync.NewCond(&e.Mutex)
	return e
}

func (e *Queue) Push(v interface{}) {
	e.Mutex.Lock()
	defer e.Mutex.Unlock()
	if !e.closed {
		e.buffer.Add(v)
		e.popable.Signal()
	}
}

func (e *Queue) Close() {
	e.Mutex.Lock()
	defer e.Mutex.Unlock()
	if !e.closed {
		e.closed = true
		e.popable.Broadcast()
	}
}

// block mode
func (e *Queue) Pop() (v interface{}) {
	c := e.popable
	buffer := e.buffer

	e.Mutex.Lock()
	defer e.Mutex.Unlock()

	for buffer.Length() == 0 && !e.closed {
		c.Wait()
	}

	if e.closed {
		return
	}

	if buffer.Length() > 0 {
		v = buffer.Peek()
		buffer.Remove()
	}
	return
}

func (e *Queue) TryPop() (v interface{}, ok bool) {
	buffer := e.buffer

	e.Mutex.Lock()
	defer e.Mutex.Unlock()

	if buffer.Length() > 0 {
		v = buffer.Peek()
		buffer.Remove()
		ok = true
	} else if e.closed {
		ok = true
	}

	return
}

func (e *Queue) Len() int {
	return e.buffer.Length()
}

func (e *Queue) Wait() {
	for {
		if e.closed || e.buffer.Length() == 0 {
			break
		}

		runtime.Gosched() // give schedule to other goroutine
	}
}
