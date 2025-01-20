// pkg/check/worker_pool.go

package check

import (
	"sync"
)

type WorkerPool struct {
	numWorkers int
	jobChan    chan func()
	stopChan   chan struct{}
	wg         sync.WaitGroup
}

func NewWorkerPool(workerCount int) *WorkerPool {
	pool := &WorkerPool{
		numWorkers: workerCount,
		jobChan:    make(chan func(), 1000),
		stopChan:   make(chan struct{}),
	}

	pool.wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go func() {
			defer pool.wg.Done()
			for {
				select {
				case job := <-pool.jobChan:
					job()
				case <-pool.stopChan:
					return
				}
			}
		}()
	}

	return pool
}

func (p *WorkerPool) Submit(job func()) {
	p.jobChan <- job
}

func (p *WorkerPool) Shutdown() {
	close(p.stopChan)
	p.wg.Wait()
}
