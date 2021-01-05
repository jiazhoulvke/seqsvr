package seqsvr

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSequenceID(t *testing.T) {
	Convey("测试序列号", t, func() {
		var err error
		err = SetMachineID(9999999)
		So(err, ShouldNotBeNil)
		err = SetMachineID(0)
		So(err, ShouldBeNil)
		err = SetMachineID(0)
		So(err, ShouldBeNil)
		seqid := SequenceID()
		So(seqid, ShouldNotEqual, 0)
		var i int64
		maker := NewMaker()
		m := make(map[int64]bool)
		var idNum int64 = 1_000_000
		for i = 0; i < idNum; i++ {
			m[maker.SequenceID()] = true
		}
		So(len(m), ShouldEqual, idNum)
	})
}

func BenchmarkSequenceID(b *testing.B) {
	sm := NewMaker()
	for i := 0; i < b.N; i++ {
		sm.SequenceID()
	}
}
