package test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/therne/lrmr/test/testutils"
)

func TestCustomPartitioner(t *testing.T) {
	Convey("Given running nodes", t, func(c C) {
		sess, stop := testutils.StartLocalCluster(c, 2)
		defer stop()

		Convey("When running stage with custom partitioner", func() {
			ds := PartitionerWithNodeAffinityTest(sess)

			Convey("It should assign rows with its designated partitions and physical nodes", func() {
				res, err := ds.Collect()
				So(err, ShouldBeNil)
				So(res, ShouldHaveLength, 2)
				So(res["1"], ShouldHaveLength, 2)
				So(res["2"], ShouldHaveLength, 2)

				So(testutils.StringValues(res["1"]), ShouldContain, "key1-1")
				So(testutils.StringValues(res["1"]), ShouldContain, "key1-2")
				So(testutils.StringValues(res["2"]), ShouldContain, "key2-1")
				So(testutils.StringValues(res["2"]), ShouldContain, "key2-2")
			})
		})
	})
}
