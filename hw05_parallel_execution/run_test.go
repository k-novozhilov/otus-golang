package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("concurrent execution without time.Sleep", func(t *testing.T) {
		const (
			taskCount   = 50
			workerCount = 5
			maxErrors   = 1
		)

		var (
			startedTasks   int32
			finishedTasks  int32
			startedSignal  = make(chan struct{})
			finishedSignal = make(chan struct{})
		)

		tasks := make([]Task, 0, taskCount)

		for i := 0; i < taskCount; i++ {
			tasks = append(tasks, func() error {
				atomic.AddInt32(&startedTasks, 1)
				if atomic.LoadInt32(&startedTasks) == 1 {
					close(startedSignal)
				}

				for j := 0; j < 1000; j++ {
					_ = j
				}

				newFinished := atomic.AddInt32(&finishedTasks, 1)
				if newFinished == int32(taskCount) {
					close(finishedSignal)
				}
				return nil
			})
		}

		go func() {
			err := Run(tasks, workerCount, maxErrors)
			require.NoError(t, err)
		}()

		<-startedSignal

		require.Eventually(t, func() bool {
			return atomic.LoadInt32(&startedTasks) > 1
		}, time.Second, 10*time.Millisecond, "tasks are not running concurrently")

		<-finishedSignal

		require.Equal(t, int32(taskCount), atomic.LoadInt32(&finishedTasks), "not all tasks were completed")
	})

	t.Run("m <= 0 should return ErrErrorsLimitExceeded", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		err := Run(tasks, 5, 0)
		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "expected ErrErrorsLimitExceeded for m=0, got %v", err)

		err = Run(tasks, 5, -1)
		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "expected ErrErrorsLimitExceeded for m=-1, got %v", err)
	})
}
