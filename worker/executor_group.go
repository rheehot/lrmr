package worker

import (
	"context"
	"fmt"
	"github.com/therne/lrmr/lrdd"
	"github.com/therne/lrmr/output"
	"github.com/therne/lrmr/transformation"
	"runtime/debug"
	"sync"
)

type executorPool struct {
	tf  transformation.Transformation
	c   *taskContext
	out output.Output

	ctx    context.Context
	Cancel context.CancelFunc

	inputChans  []chan incomingData
	handles     []*executorHandle
	Concurrency int
}

type incomingData struct {
	handle *executorHandle
	data   lrdd.Row
}

func launchExecutorPool(
	tf transformation.Transformation,
	c *taskContext,
	out output.Output,
	concurrency, queueLen int,
) *executorPool {

	execCtx, cancel := context.WithCancel(context.Background())
	eg := &executorPool{
		tf:          tf,
		c:           c,
		out:         out,
		ctx:         execCtx,
		Cancel:      cancel,
		inputChans:  make([]chan incomingData, concurrency),
		Concurrency: concurrency,
	}

	// spawn executors
	for i := range eg.inputChans {
		eg.inputChans[i] = make(chan incomingData, queueLen)
		go eg.startConsume(execCtx, eg.inputChans[i], i)
	}
	return eg
}

func (teg *executorPool) NewExecutorHandle() *executorHandle {
	conn := &executorHandle{
		Errors: make(chan error),
		teg:    teg,
	}
	teg.handles = append(teg.handles, conn)
	return conn
}

func (teg *executorPool) startConsume(ctx context.Context, inputChan chan incomingData, executorID int) {
	var in incomingData
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic: %v\n%s", r, string(debug.Stack()))
			in.handle.Errors <- err
		}
	}()
	for {
		select {
		case in = <-inputChan:
			childCtx := teg.c.forkForExecutor(executorID)
			err := teg.tf.Apply(childCtx, in.data, teg.out)
			in.handle.wg.Done()
			if err != nil {
				in.handle.Errors <- err
				teg.Cancel()
				return
			}

		case <-ctx.Done():
			// eg.Cancel() called
			return
		}
	}
}

func (teg *executorPool) WaitForCompletion() {
	println("Waiting for finish...")
	for _, handle := range teg.handles {
		handle.WaitForCompletion()
		close(handle.Errors)
	}
	teg.Cancel()
}

type executorHandle struct {
	Errors     chan error
	totalTasks int

	teg *executorPool
	wg  sync.WaitGroup
}

// Enqueue runs task asynchronously with load balancing by given concurrency.
func (eh *executorHandle) Enqueue(data lrdd.Row) {
	if eh.teg.ctx.Err() != nil {
		// executor is already terminated
		return
	}
	slot := eh.totalTasks % eh.teg.Concurrency
	dataQueueChan := eh.teg.inputChans[slot]
	eh.wg.Add(1)
	eh.totalTasks += 1
	dataQueueChan <- incomingData{data: data, handle: eh}
}

func (eh *executorHandle) WaitForCompletion() {
	eh.wg.Wait()
}
