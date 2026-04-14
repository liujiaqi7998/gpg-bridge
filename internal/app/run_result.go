package app

import "sync"

type runResult struct {
	done chan struct{}
	once sync.Once
	mu   sync.Mutex
	err  error
}

func newRunResult() *runResult {
	return &runResult{done: make(chan struct{})}
}

func (r *runResult) finish(err error) {
	r.once.Do(func() {
		r.mu.Lock()
		r.err = err
		r.mu.Unlock()
		close(r.done)
	})
}

func (r *runResult) wait() error {
	<-r.done
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.err
}
