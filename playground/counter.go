package playground

import (
	"fmt"
	"github.com/therne/lrmr/lrdd"
	"github.com/therne/lrmr/output"
	"github.com/therne/lrmr/transformation"
	"sync/atomic"
)

type Counter struct {
	counter int64
}

func (cnt *Counter) DescribeOutput() *transformation.OutputDesc {
	return transformation.DescribingOutput().Nothing()
}

func (cnt *Counter) Setup(c transformation.Context) error {
	return nil
}

func (cnt *Counter) Run(row lrdd.Row, out output.Output) error {
	atomic.AddInt64(&cnt.counter, 1)
	return nil
}

func (cnt *Counter) Teardown(out output.Output) error {
	fmt.Printf("Result Count: %d\n", cnt.counter)
	return nil
}