package seqsvr

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSequenceID(t *testing.T) {
	Convey("测试序列号", t, func() {
		var err error
		err = SetMachineID(0)
		So(err, ShouldNotBeNil)
		err = SetMachineID(9999999)
		So(err, ShouldNotBeNil)
		err = SetMachineID(1)
		So(err, ShouldBeNil)
		seqid := SequenceID()
		t.Logf("十进制序列号:%d\n", seqid)
		t.Logf("十六进制序列号:%X\n", seqid)
		So(seqid, ShouldNotEqual, 0)
		var i int64
		maker := NewMaker("default")
		for i = 0; i < Max; i++ {
			maker.SequenceID()
		}
	})
}

func BenchmarkSequenceID(b *testing.B) {
	sm := NewMaker("test")
	for i := 0; i < b.N; i++ {
		sm.SequenceID()
	}
}
