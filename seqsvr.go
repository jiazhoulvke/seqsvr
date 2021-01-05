package seqsvr

//生成规则:
//最高位 1bit  值为0
//时间戳 31bit 最大日期为2038年1月19号
//机器ID 9bit 最大支持1024台序列号服务器
//序列号 23bit 每台服务器每秒最多可以生成1<<23-1(8388607)个ID

import (
	"fmt"
	"sync"
	"time"
)

const (
	//Max 一秒内最大可生成的id数量
	Max int64 = 1 << 23 //8388607
)

var (
	//MaxMachineNumber 最大序列号服务器数量
	MaxMachineNumber int64 = 1 << 9
	defaultMaker     *Maker
)

func init() {
	defaultMaker = NewMaker()
}

//SetMachineID 设置机器ID
func SetMachineID(machineID int64) error {
	return defaultMaker.SetMachineID(machineID)
}

//SequenceID 获取一个序列号
func SequenceID() int64 {
	return defaultMaker.SequenceID()
}

//NewMaker 获得一个生成器
func NewMaker(s ...string) *Maker {
	u := Maker{
		curTime: time.Now().Unix(),
		count:   1,
	}
	return &u
}

//Maker unique id 生成器
type Maker struct {
	curTime   int64
	count     int64
	machineID int64
	mutex     sync.RWMutex
}

func (m *Maker) SetMachineID(machineID int64) error {
	if machineID > MaxMachineNumber {
		return fmt.Errorf("机器ID不能大于%d", MaxMachineNumber)
	}
	m.mutex.Lock()
	m.machineID = machineID
	m.mutex.Unlock()
	return nil
}

//SequenceID 获取序列号
func (s *Maker) SequenceID() (sequenceID int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	now := time.Now().Unix()
	if s.count > Max {
		for ; s.curTime == now; now = time.Now().Unix() {
			s.curTime = now
			time.Sleep(time.Duration(100) * time.Millisecond)
		}
		s.count = 1
	}
	if s.curTime != now {
		s.curTime = now
		s.count = 1
	}
	sequenceID = s.curTime<<32 + s.machineID<<23 + s.count
	s.count++
	return
}
