package common

import (
	"server/base"
	. "server/common"
	"server/timer"
	"sync/atomic"
	"time"
)

type TimerObject struct {
	id       int
	Count_1  int64
	Count_2  int64
	Second_1 int
	Second_2 int
}

func (o *TimerObject) ID() int {
	return o.id
}

func (o *TimerObject) Create() {
	timer.AddObject(o.ID())
}

func (o *TimerObject) Delete() {
	timer.RemoveObject(o.ID())
}

func (o *TimerObject) TimerCallback1(args ...interface{}) {
	//LogInfo("TimerCallback1")
	atomic.AddInt64(&o.Count_1, 1)
	o.Second_1 = GetSecond()
}

func (o *TimerObject) TimerCallback2(args ...interface{}) {
	//LogInfo("TimerCallback2")
	atomic.AddInt64(&o.Count_2, 1)
	o.Second_2 = GetSecond()
}

const (
	constTimerModule   = "test"
	constTimerEvent    = "test_event"
	constTimerModuleCh = "test_ch"
	constTimerEventCh  = "test_ch_event"

	constTimerCount = 1000
)

func createTimerObject(id int) *TimerObject {
	obj := TimerObject{id: id, Count_1: 0, Count_2: 0}
	obj.Create()
	return &obj
}

func TimerCallback(args []interface{}) {
	f := args[0].(Functor)
	f.Call(args[1:]...)
}

func TestTimerWithChannel(secs []int) {
	//注册通过channel通知的定时器事件
	var skeleton = base.NewSkeleton()
	skeleton.RegisterChanRPC(constTimerEventCh, TimerCallback)
	timer.RegistModule(constTimerModuleCh, constTimerEventCh, skeleton.ChanRPCServer)
	go func() {
		for {
			info := <-skeleton.ChanRPCServer.ChanCall
			skeleton.ChanRPCServer.Exec(info)
		}
	}()

	id := NewObjectID()
	LogInfo("NewObjectID: %d", id)
	obj := createTimerObject(id)

	for _, sec := range secs {
		for i := 0; i < constTimerCount; i++ {
			f1 := NewFunctor("TimerCallback1", obj.TimerCallback1)
			f2 := NewFunctor("TimerCallback2", obj.TimerCallback2)
			timer.StartTimer(constTimerModuleCh, id, "TimerCallback1", sec, f1)
			timer.StartTimer(constTimerModuleCh, id, "TimerCallback2", sec, f2)
		}
	}
	timer.RemoveTimer(constTimerModuleCh, id, "TimerCallback2")

	start := GetSecond()
	//等待所有定时器结束
	sum := 0
	for _, sec := range secs {
		if sum <= sec {
			sum = sec
		}
	}
	time.Sleep(time.Duration(sum+1) * time.Second)

	total := int64(len(secs) * constTimerCount)
	if obj.Count_1 != total || obj.Second_1 != start+sum {
		LogPanic("TestTimerWithChannel fail! obj %v total = %d", obj, total)
	}
	if obj.Count_2 != 0 {
		LogPanic("TestTimerWithChannel fail! obj.Count_2 %d != 0", obj.Count_2)
	}
	LogInfo("TestTimerWithChannel success!")
}

func TestTimerNoChannel(secs []int) {
	//注册通过channel通知的定时器事件
	timer.RegistModule(constTimerModule, constTimerEvent, nil)

	id := NewObjectID()
	LogInfo("NewObjectID: %d", id)
	obj := createTimerObject(id)

	for _, sec := range secs {
		for i := 0; i < constTimerCount; i++ {
			f1 := NewFunctor("TimerCallback1", obj.TimerCallback1)
			f2 := NewFunctor("TimerCallback2", obj.TimerCallback2)
			timer.StartTimer(constTimerModule, id, "TimerCallback1", sec, f1)
			timer.StartTimer(constTimerModule, id, "TimerCallback2", sec, f2)
		}
	}
	timer.RemoveTimer(constTimerModule, id, "TimerCallback2")

	start := GetSecond()
	//等待所有定时器结束
	sum := 0
	for _, sec := range secs {
		if sum <= sec {
			sum = sec
		}
	}
	time.Sleep(time.Duration(sum+1) * time.Second)

	total := int64(len(secs) * constTimerCount)
	if obj.Count_1 != total || obj.Second_1 != start+sum {
		LogPanic("TestTimerWithChannel fail! obj %v total = %d", obj, total)
	}
	if obj.Count_2 != 0 {
		LogPanic("TestTimerWithChannel fail! obj.Count_2 %d != 0", obj.Count_2)
	}
	LogInfo("TestTimerNoChannel success!")
}
