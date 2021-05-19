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
}

func TestRunAllSuccessTask(t *testing.T) {
	defer goleak.VerifyNone(t)
	t.Run("single task, single worker, no error", func(t *testing.T) {
		makeTest(t, 1, 1)
	})

	t.Run("single task, several worker, no error", func(t *testing.T) {
		makeTest(t, 1, 5)
	})

	t.Run("several task, single worker, no error", func(t *testing.T) {
		makeTest(t, 6, 1)
	})

	t.Run("several task, many worker, no error", func(t *testing.T) {
		makeTest(t, 6, 20)
	})

	t.Run("many task, several worker, no error", func(t *testing.T) {
		makeTest(t, 61, 5)
	})
}

func TestRunSingleTask(t *testing.T) {
	defer goleak.VerifyNone(t)
	t.Run("single task - single error", func(t *testing.T) {
		tasksCount := 1
		tasks := make([]Task, 0, tasksCount)
		var runTasksCount int32
		tasks = append(tasks, makeBadTask(&runTasksCount))

		workersCount := 1
		maxErrorsCount := 1
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
	})

	t.Run("single task - no error", func(t *testing.T) {
		tasksCount := 1
		tasks := make([]Task, 0, tasksCount)
		var runTasksCount int32
		tasks = append(tasks, makeSuccessTask(&runTasksCount))

		workersCount := 1
		maxErrorsCount := 1
		err := Run(tasks, workersCount, maxErrorsCount)

		require.NoError(t, err)
		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
	})
}

func TestRunFirstTaskWithError(t *testing.T) {
	defer goleak.VerifyNone(t)
	t.Run("many tasks - first error, single worker", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, 0, tasksCount)
		var runTasksCount int32
		tasks = append(tasks, makeBadTask(&runTasksCount))

		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, makeSuccessTask(&runTasksCount))
		}

		workersCount := 1
		maxErrorsCount := 1
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.Equal(t, int32(workersCount+1), runTasksCount, "not all tasks were completed")
	})

	t.Run("many tasks - first error, many worker", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, 0)
		var runTasksCount int32

		tasks = append(tasks, makeBadTask(&runTasksCount))
		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, makeSuccessTask(&runTasksCount))
		}

		workersCount := 4
		maxErrorsCount := 1
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, int32(workersCount+1), runTasksCount, "not all tasks were completed")
	})
}

func TestRunMiddleTaskWithError(t *testing.T) {
	defer goleak.VerifyNone(t)
	t.Run("many tasks - middle error, single worker", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, 0)
		var runTasksCount int32
		badWorkerIndex := 5

		for i := 0; i < tasksCount; i++ {
			if i == badWorkerIndex {
				tasks = append(tasks, makeBadTask(&runTasksCount))
			} else {
				tasks = append(tasks, makeSuccessTask(&runTasksCount))
			}
		}
		workersCount := 1
		maxErrorsCount := 1
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.Equal(t, int32(badWorkerIndex+2), runTasksCount, "not all tasks were completed")
	})
}

func TestRunLastTaskWithError(t *testing.T) {
	defer goleak.VerifyNone(t)
	t.Run("many tasks - last error, single worker", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, 0)
		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, makeSuccessTask(&runTasksCount))
		}
		tasks = append(tasks, makeBadTask(&runTasksCount))

		workersCount := 1
		maxErrorsCount := 1
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.Equal(t, int32(tasksCount+1), runTasksCount, "not all tasks were completed")
	})

	t.Run("many tasks - last error, many worker", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, 0)
		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, makeSuccessTask(&runTasksCount))
		}
		tasks = append(tasks, makeBadTask(&runTasksCount))

		workersCount := 4
		maxErrorsCount := 1
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.Equal(t, int32(tasksCount+1), runTasksCount, "not all tasks were completed")
	})
}

func TestRunWithBadParam(t *testing.T) {
	defer goleak.VerifyNone(t)
	t.Run("Empty tasks", func(t *testing.T) {
		tasks := make([]Task, 0)
		workersCount := 1
		maxErrorsCount := 1
		err := Run(tasks, workersCount, maxErrorsCount)

		require.NoError(t, err)
	})

	t.Run("Empty max errors", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, 0)
		var runTasksCount int32
		badWorkerIndex := 5

		for i := 0; i < tasksCount; i++ {
			if i == badWorkerIndex {
				tasks = append(tasks, makeBadTask(&runTasksCount))
			} else {
				tasks = append(tasks, makeSuccessTask(&runTasksCount))
			}
		}

		workersCount := 1
		maxErrorsCount := 0
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.Equal(t, int32(badWorkerIndex+2), runTasksCount, "not all tasks were completed")
	})

	t.Run("Empty workers", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, 0)
		var runTasksCount int32
		badWorkerIndex := 5

		for i := 0; i < tasksCount; i++ {
			if i == badWorkerIndex {
				tasks = append(tasks, makeBadTask(&runTasksCount))
			} else {
				tasks = append(tasks, makeSuccessTask(&runTasksCount))
			}
		}

		workersCount := 0
		maxErrorsCount := 1
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.Equal(t, int32(badWorkerIndex+2), runTasksCount, "not all tasks were completed")
	})
}

// Сохранил на будущее пример из разбора задания.
func TestRunConcurrency(t *testing.T) {
	const workersCount = 5
	tasks := make([]Task, workersCount)

	waitCh := make(chan struct{})
	var runTasksCount int32
	for i := 0; i < len(tasks); i++ {
		tasks[i] = func() error {
			atomic.AddInt32(&runTasksCount, 1)
			<-waitCh
			return nil
		}
	}

	runErrCh := make(chan error, 1)
	go func() {
		runErrCh <- Run(tasks, workersCount, 0)
	}()

	require.Eventually(t, func() bool {
		return atomic.LoadInt32(&runTasksCount) == workersCount
	}, time.Second, time.Millisecond)

	close(waitCh)

	var runErr error
	require.Eventually(t, func() bool {
		select {
		case runErr = <-runErrCh:
			return true
		default:
			return false
		}
	}, time.Second, time.Millisecond)

	require.NoError(t, runErr)
}

func makeTest(t *testing.T, tasksCount int, workersCount int) {
	tasks := make([]Task, 0, tasksCount)

	var runTasksCount int32
	maxErrorsCount := 1
	for i := 0; i < tasksCount; i++ {
		tasks = append(tasks, makeSuccessTask(&runTasksCount))
	}

	err := Run(tasks, workersCount, maxErrorsCount)

	require.NoError(t, err)
	require.Equal(t, int32(tasksCount), runTasksCount, "not all tasks were completed")
}

func makeSuccessTask(runTasksCount *int32) func() error {
	taskSleep := time.Millisecond * 10

	f := func() error {
		time.Sleep(taskSleep)
		atomic.AddInt32(runTasksCount, 1)
		return nil
	}

	return f
}

func makeBadTask(runTasksCount *int32) func() error {
	taskSleep := time.Millisecond * 10

	f := func() error {
		time.Sleep(taskSleep)
		atomic.AddInt32(runTasksCount, 1)
		return errors.New("Bad task here")
	}

	return f
}
