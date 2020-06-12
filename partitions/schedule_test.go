package partitions

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/therne/lrmr/lrdd"
	"github.com/therne/lrmr/node"
	"testing"
)

func TestScheduler_AffinityRule(t *testing.T) {
	Convey("Given a partition.Scheduler", t, func() {
		Convey("When executors are sufficient", func() {
			nn := []*node.Node{
				{Host: "localhost:1001", Executors: 3, Tag: map[string]string{"CustomTag": "hello"}},
				{Host: "localhost:1002", Executors: 3, Tag: map[string]string{"CustomTag": "world"}},
				{Host: "localhost:1003", Executors: 3, Tag: map[string]string{"CustomTag": "foo"}},
				{Host: "localhost:1004", Executors: 3, Tag: map[string]string{"CustomTag": "bar"}},
			}

			Convey("When an affinity rule is given with an Partitioner", func() {
				_, aa := Schedule(nn, []Plan{
					{Partitioner: partitionerStub{[]Partition{
						{ID: "familiarWithWorld", AssignmentAffinity: map[string]string{"Host": "localhost:1002"}},
						{ID: "familiarWithFoo", AssignmentAffinity: map[string]string{"CustomTag": "foo"}},
						{ID: "familiarWithBar", AssignmentAffinity: map[string]string{"CustomTag": "bar"}},
						{ID: "familiarWithFreest"},
						{ID: "familiarWithWorld2", AssignmentAffinity: map[string]string{"Host": "localhost:1002"}},
						{ID: "familiarWithFoo2", AssignmentAffinity: map[string]string{"CustomTag": "foo"}},
						{ID: "familiarWithBar2", AssignmentAffinity: map[string]string{"CustomTag": "bar"}},
						{ID: "familiarWithFreest2"},
					}}},
					{ /* ignored */ },
				}, WithoutShufflingNodes())
				So(aa, ShouldHaveLength, 2)
				So(aa[1], ShouldHaveLength, 8)

				keyToHostMap := aa[1].ToMap()
				So(keyToHostMap["familiarWithWorld"], ShouldEqual, "localhost:1002")
				So(keyToHostMap["familiarWithFoo"], ShouldEqual, "localhost:1003")
				So(keyToHostMap["familiarWithBar"], ShouldEqual, "localhost:1004")
				So(keyToHostMap["familiarWithFreest"], ShouldEqual, "localhost:1001")
				So(keyToHostMap["familiarWithWorld2"], ShouldEqual, "localhost:1002")
				So(keyToHostMap["familiarWithFoo2"], ShouldEqual, "localhost:1003")
				So(keyToHostMap["familiarWithBar2"], ShouldEqual, "localhost:1004")
				So(keyToHostMap["familiarWithFreest2"], ShouldEqual, "localhost:1001")
			})
		})

		Convey("When executors are scarce", func() {
			nn := []*node.Node{
				{Host: "localhost:1001", Executors: 1, Tag: map[string]string{"CustomTag": "hello"}},
				{Host: "localhost:1002", Executors: 1, Tag: map[string]string{"CustomTag": "world"}},
				{Host: "localhost:1003", Executors: 1, Tag: map[string]string{"CustomTag": "foo"}},
				{Host: "localhost:1004", Executors: 1, Tag: map[string]string{"CustomTag": "bar"}},
			}

			Convey("When an affinity rule is given with an LogicalPlanner", func() {
				_, pp := Schedule(nn, []Plan{
					{Partitioner: partitionerStub{[]Partition{
						{ID: "familiarWithWorld", AssignmentAffinity: map[string]string{"Host": "localhost:1002"}},
						{ID: "familiarWithFoo", AssignmentAffinity: map[string]string{"CustomTag": "foo"}},
						{ID: "familiarWithBar", AssignmentAffinity: map[string]string{"CustomTag": "bar"}},
						{ID: "familiarWithFreest"},
						{ID: "familiarWithWorld2", AssignmentAffinity: map[string]string{"Host": "localhost:1002"}},
						{ID: "familiarWithFoo2", AssignmentAffinity: map[string]string{"CustomTag": "foo"}},
						{ID: "familiarWithBar2", AssignmentAffinity: map[string]string{"CustomTag": "bar"}},
						{ID: "familiarWithFreest2"},
					}}},
					{ /* ignored */ },
				}, WithoutShufflingNodes())
				So(pp, ShouldHaveLength, 2)
				So(pp[1], ShouldHaveLength, 8)

				keyToHostMap := pp[1].ToMap()
				So(keyToHostMap["familiarWithWorld"], ShouldEqual, "localhost:1002")
				So(keyToHostMap["familiarWithFoo"], ShouldEqual, "localhost:1003")
				So(keyToHostMap["familiarWithBar"], ShouldEqual, "localhost:1004")
				So(keyToHostMap["familiarWithFreest"], ShouldEqual, "localhost:1001")
				So(keyToHostMap["familiarWithWorld2"], ShouldEqual, "localhost:1002")
				So(keyToHostMap["familiarWithFoo2"], ShouldEqual, "localhost:1003")
				So(keyToHostMap["familiarWithBar2"], ShouldEqual, "localhost:1004")
				So(keyToHostMap["familiarWithFreest2"], ShouldEqual, "localhost:1001")
			})
		})
	})
}

type partitionerStub struct {
	Partitions []Partition
}

func (p partitionerStub) PlanNext(int) []Partition {
	return p.Partitions
}

func (p partitionerStub) DeterminePartition(c Context, r *lrdd.Row) (id string, err error) {
	return
}
