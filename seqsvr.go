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
	makers           map[string]*Maker
	mutex            sync.RWMutex
)

func init() {
	makers = make(map[string]*Maker)
	defaultMaker = NewMaker("default")
}

//Maker unique id 生成器
type Maker struct {
	curTime   int64
	count     int64
	machineID int64
	mutex     sync.RWMutex
}

//SetMachineID 设置机器ID
func (s *Maker) SetMachineID(machineID int64) error {
	if machineID < 1 {
		return fmt.Errorf("机器ID必须为大于0的正整数")
	}
	if machineID > MaxMachineNumber {
		return fmt.Errorf("机器ID不能大于%d", MaxMachineNumber)
	}
	s.machineID = machineID
	return nil
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
func NewMaker(key string) *Maker {
	mutex.Lock()
	defer mutex.Unlock()
	if _, ok := makers[key]; ok {
		return makers[key]
	}
	u := Maker{
		curTime:   time.Now().Unix(),
		count:     1,
		machineID: 1,
	}
	makers[key] = &u
	return &u
}

//SequenceID 获取序列号
func (s *Maker) SequenceID() (sequenceID int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.count > Max {
		for now := time.Now().Unix(); s.curTime == now; now = time.Now().Unix() {
			s.curTime = now
			time.Sleep(time.Duration(100) * time.Millisecond)
		}
		s.count = 1
	}
	sequenceID = s.curTime<<32 + s.machineID<<23 + s.count
	s.count++
	return
}
