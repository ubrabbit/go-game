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
	Count_3  int64
	Count_4  int64
	Second_1 int
	Second_2 int
	Second_3 int
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

func (o *TimerObject) TimerCallback3(args ...interface{}) {
	//LogInfo("TimerCallback3")
	atomic.AddInt64(&o.Count_3, 1)
	o.Second_3 = GetSecond()
}

const (
	constTimerModule    = "test"
	constTimerModuleCh  = "test_ch"
	constTimerModuleCh2 = "test_ch2"

	constTimerCount = 1000
)

func createTimerObject(id int) *TimerObject {
	obj := TimerObject{id: id, Count_1: 0, Count_2: 0}
	obj.Create()
	return &obj
}

func TimerCallback(args []interface{}) {
	i := args[0].(*timer.TimerItem)
	//LogInfo("TimerCallback ")
	i.Execute()
}

func TestTimerWithChannel(secs []int) {
	//注册通过channel通知的定时器事件
	var skeleton = base.NewSkeleton()
	skeleton.RegisterChanRPC(constTimerModuleCh, TimerCallback)
	timer.RegistModule(constTimerModuleCh, skeleton.ChanRPCServer)
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
			f3 := NewFunctor("TimerCallback2", obj.TimerCallback3)
			key := FormatString("TimerCallback1_%d", i)
			timer.StartTimer(constTimerModuleCh, id, key, sec, f1)
			timer.StartTimer(constTimerModuleCh, id, "TimerCallback2", sec, f2)
			timer.StartTimer(constTimerModuleCh, id, "TimerCallback3", sec, f3)
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

	total := int64(constTimerCount)
	if obj.Count_1 != total || obj.Second_1 != start+sum {
		LogPanic("TestTimerWithChannel fail! obj obj.Count_1 %d != total %d", obj.Count_1, total)
	}
	if obj.Count_2 != 0 {
		LogPanic("TestTimerWithChannel fail! obj.Count_2 %d != 0", obj.Count_2)
	}
	if obj.Count_3 != 1 {
		LogPanic("TestTimerWithChannel fail! obj.Count_3 %d != 1", obj.Count_3)
	}
	LogInfo("TestTimerWithChannel success!")
}

func TestTimerNoChannel(secs []int) {
	//注册通过channel通知的定时器事件
	timer.RegistModule(constTimerModule, nil)

	id := NewObjectID()
	LogInfo("NewObjectID: %d", id)
	obj := createTimerObject(id)

	for _, sec := range secs {
		for i := 0; i < constTimerCount; i++ {
			f1 := NewFunctor("TimerCallback1", obj.TimerCallback1)
			f2 := NewFunctor("TimerCallback2", obj.TimerCallback2)
			f3 := NewFunctor("TimerCallback2", obj.TimerCallback3)
			key := FormatString("TimerCallback1_%d", i)
			timer.StartTimer(constTimerModuleCh, id, key, sec, f1)
			timer.StartTimer(constTimerModuleCh, id, "TimerCallback2", sec, f2)
			timer.StartTimer(constTimerModuleCh, id, "TimerCallback3", sec, f3)
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

	total := int64(constTimerCount)
	if obj.Count_1 != total || obj.Second_1 != start+sum {
		LogPanic("TestTimerWithChannel fail! obj %v total = %d", obj, total)
	}
	if obj.Count_2 != 0 {
		LogPanic("TestTimerWithChannel fail! obj.Count_2 %d != 0", obj.Count_2)
	}
	if obj.Count_3 != 1 {
		LogPanic("TestTimerWithChannel fail! obj.Count_3 %d != 1", obj.Count_3)
	}
	LogInfo("TestTimerNoChannel success!")
}

func (o *TimerObject) TimerCallbackGroutineSafe(args ...interface{}) {
	//LogInfo("TimerCallbackGroutineSafe")
	LogInfo("这个函数不应该被执行")
	atomic.AddInt64(&o.Count_4, 1)
}
func TestTimerGoroutineSafe() {
	//注册通过channel通知的定时器事件
	var skeleton = base.NewSkeleton()
	skeleton.RegisterChanRPC(constTimerModuleCh2, TimerCallback)
	timer.RegistModule(constTimerModuleCh2, skeleton.ChanRPCServer)
	go func() {
		for {
			info := <-skeleton.ChanRPCServer.ChanCall
			LogInfo("定时器回调，睡眠1秒模拟chan延迟处理 GetGoroutineID: %s", GetGoroutineID())
			time.Sleep(1 * time.Second)
			skeleton.ChanRPCServer.Exec(info)
		}
	}()

	id := NewObjectID()
	LogInfo("NewObjectID: %d GetGoroutineID: %s", id, GetGoroutineID())
	obj := createTimerObject(id)

	f1 := NewFunctor("TimerCallback1", obj.TimerCallbackGroutineSafe)
	LogInfo("定时器1.1秒执行")
	key := "TimerCallback1"
	timer.StartTimer(constTimerModuleCh2, id, key, 1, f1)
	time.Sleep(time.Duration(1100) * time.Millisecond)
	LogInfo("1秒到期，模拟在主Goroutine删除了定时器")
	timer.RemoveTimer(constTimerModuleCh2, id, key)
	LogInfo("睡眠2秒等待回调函数执行")
	time.Sleep(time.Duration(2) * time.Second)

	LogInfo("检查定时器是否被执行，正确的是不应该执行")
	if obj.Count_4 != 0 {
		LogPanic("TestTimerGoroutineSafe fail! obj.Count_4 %d != 0", obj.Count_4)
	}
	LogInfo("TestTimerGoroutineSafe success!")
}
