package workerpool

import "sync/atomic"

type WorkerPool struct {
	Max    int
	Curr   *AtomicInt64
	Queued *AtomicInt64
}

func (wp *WorkerPool) Go(fn func()) {
	wp.Queued.Inc()

	for wp.Curr.Get() >= wp.Max {
	}

	wp.Curr.Inc()
	wp.Queued.Dec()

	go func() {
		defer func() { wp.Curr.Dec() }()

		fn()
	}()
}

type AtomicInt64 int64

func (ai *AtomicInt64) Inc() {
	atomic.AddInt64((*int64)(ai), 1)
}

func (ai *AtomicInt64) Dec() {
	atomic.AddInt64((*int64)(ai), -1)
}

func (ai *AtomicInt64) Get() int {
	return int(atomic.LoadInt64((*int64)(ai)))
}
