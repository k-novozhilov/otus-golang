package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	if len(tasks) == 0 {
		return nil
	}

	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	var errorsCount int32
	tasksCh := make(chan Task, len(tasks))
	doneCh := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(n)

	startWorkers(&wg, n, tasksCh, doneCh, &errorsCount, m)
	sendTasks(tasks, tasksCh, &errorsCount, m)
	close(tasksCh)

	if int(atomic.LoadInt32(&errorsCount)) >= m {
		close(doneCh)
	}

	wg.Wait()

	if doneCh != nil && !isClosed(doneCh) {
		close(doneCh)
	}

	if int(atomic.LoadInt32(&errorsCount)) >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func startWorkers(
	wg *sync.WaitGroup,
	count int,
	tasksCh <-chan Task,
	doneCh <-chan struct{},
	errorsCount *int32,
	errLimit int,
) {
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case task, ok := <-tasksCh:
					if !ok {
						return
					}

					if int(atomic.LoadInt32(errorsCount)) >= errLimit {
						return
					}

					if err := task(); err != nil {
						atomic.AddInt32(errorsCount, 1)
					}
				case <-doneCh:
					return
				}
			}
		}()
	}
}

func sendTasks(tasks []Task, tasksCh chan<- Task, errorsCount *int32, errLimit int) {
	for _, task := range tasks {
		if int(atomic.LoadInt32(errorsCount)) >= errLimit {
			break
		}
		tasksCh <- task
	}
}

func isClosed(ch <-chan struct{}) bool {
	select {
	case <-ch:
		return true
	default:
		return false
	}
}
