package test

import (
	"time"

	"github.com/therne/lrmr"
	"github.com/therne/lrmr/lrdd"
)

var _ = lrmr.RegisterTypes(HaltForMasterFailure{})

type HaltForMasterFailure struct{}

func (f HaltForMasterFailure) Transform(ctx lrmr.Context, in chan *lrdd.Row, emit func(*lrdd.Row)) error {
	time.Sleep(5 * time.Second)
	for range in {
		ctx.AddMetric("Input", 1)
	}
	return nil
}

func RunWithoutMaster(sess *lrmr.Session) *lrmr.Dataset {
	return sess.Parallelize([]int{1, 2, 3, 4, 5}).
		Do(HaltForMasterFailure{})
}