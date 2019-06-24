package common

import (
	"server/base"
	. "server/common"
	"server/event"
	gevent "server/gamedata/event"
	"server/leaf/module"
	"sync"
	"time"
)

type EventTestObj struct {
	ID     int
	Count1 int
	Count2 int
}

func (o *EventTestObj) EventCallback1(args ...interface{}) {
	//e := args[0].(event.Event)
	//LogDebug("%d EventCallback1 %v", o.ID, e.Args())
	o.Count1++
}

func (o *EventTestObj) EventCallback2(args ...interface{}) {
	//e := args[0].(event.Event)
	//LogDebug("%d EventCallback2 %v", o.ID, e.Args())
	o.Count2++
}

const (
	constObjectNum = 10000
	constTimeout1  = 100
	constTimeout2  = 50
)

func TestEventCallback(args []interface{}) {
	defer func() {
		r := recover()
		if r != nil {
			LogError("Callback Error: %v %v", args, r)
		}
	}()
	//LogDebug("TestEventCallback: %v", args)
	f := args[0].(*Functor)
	f.Call(args[1])
}

var skeleton *module.Skeleton = nil

func InitEvent() {
	if skeleton != nil {
		return
	}
	skeleton = base.NewSkeleton()
	event.RegistEventModule(event.EVENT_MODULE_GAME, skeleton.ChanRPCServer)
	skeleton.RegisterChanRPC(event.EVENT_MODULE_GAME, TestEventCallback)
	go func() {
		for {
			info := <-skeleton.ChanRPCServer.ChanCall
			skeleton.ChanRPCServer.Exec(info)
		}
	}()
}

func TestEvent_1() {
	InitEvent()

	l := event.CreateListener(event.EVENT_MODULE_GAME)
	objList := make(map[int]*EventTestObj, 0)
	for i := 0; i < constObjectNum; i++ {
		obj := EventTestObj{
			ID:     NewObjectID(),
			Count1: 0,
			Count2: 0,
		}
		e1 := gevent.Event1001{}
		e2 := gevent.Event1002{}
		l.AddListen(&e1, obj.ID, NewFunctor("EventCallback1", obj.EventCallback1))
		l.AddListen(&e2, obj.ID, NewFunctor("EventCallback2", obj.EventCallback2))
		objList[obj.ID] = &obj
	}

	wg := new(sync.WaitGroup)
	wg.Add(1)
	count := 0
	go func() {
		for i := 0; i < constTimeout1; i++ {
			count++
			e1 := gevent.Event1001{}
			e2 := gevent.Event1002{}
			l.TriggerEvent(&e1)
			l.TriggerEvent(&e2)
			if count == constTimeout2 {
				for id, _ := range objList {
					l.RemoveListen(&e2, id)
				}
			}
			time.Sleep(10 * time.Millisecond)
		}
		wg.Done()
	}()
	wg.Wait()
	time.Sleep(10 * time.Millisecond)

	countTotal1 := constTimeout1
	countTotal2 := constTimeout2
	for id, obj := range objList {
		if obj.Count1 != countTotal1 {
			LogPanic("%d Count1 %d != %d", id, obj.Count1, countTotal1)
		}
		if obj.Count2 != countTotal2 {
			LogPanic("%d Count2 %d != %d", id, obj.Count2, countTotal2)
		}
	}

	LogInfo("TestEvent_1 success !")
}

//多个goroutine同时读写触发事件
func TestEvent_2() {
	InitEvent()

	l := event.CreateListener(event.EVENT_MODULE_GAME)
	objList := make(map[int]*EventTestObj, 0)
	for i := 0; i < constObjectNum; i++ {
		obj := EventTestObj{
			ID:     NewObjectID(),
			Count1: 0,
			Count2: 0,
		}
		e1 := gevent.Event1001{}
		e2 := gevent.Event1002{}
		l.AddListen(&e1, obj.ID, NewFunctor("EventCallback1", obj.EventCallback1))
		l.AddListen(&e2, obj.ID, NewFunctor("EventCallback2", obj.EventCallback2))
		objList[obj.ID] = &obj
	}

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		for i := 0; i < constTimeout1; i++ {
			e1 := gevent.Event1001{}
			e2 := gevent.Event1002{}
			l.TriggerEvent(&e1)
			l.TriggerEvent(&e2)
			time.Sleep(10 * time.Millisecond)
		}
		wg.Done()
	}()
	go func() {
		time.Sleep(constTimeout2 * time.Millisecond)
		e2 := gevent.Event1002{}
		for id, _ := range objList {
			l.RemoveListen(&e2, id)
		}
	}()
	wg.Wait()

	countTotal1 := constTimeout1
	countTotal2 := constTimeout2
	for id, obj := range objList {
		if obj.Count1 != countTotal1 {
			LogPanic("%d Count1 %d != %d", id, obj.Count1, countTotal1)
		}
		if obj.Count2 >= countTotal2 {
			LogPanic("%d Count2 %d >= %d", id, obj.Count2, countTotal2)
		}
	}
	LogInfo("TestEvent_2 success !")
}

//测试事件的读写性能
func TestEvent_3() {
	l := event.CreateListener(event.EVENT_MODULE_GAME)
	objList := make(map[int]*EventTestObj, 0)
	for i := 0; i < constObjectNum; i++ {
		obj := EventTestObj{
			ID:     NewObjectID(),
			Count1: 0,
			Count2: 0,
		}
		e1 := gevent.Event1001{}
		e2 := gevent.Event1002{}
		l.AddListen(&e1, obj.ID, NewFunctor("EventCallback1", obj.EventCallback1))
		l.AddListen(&e2, obj.ID, NewFunctor("EventCallback2", obj.EventCallback2))
		objList[obj.ID] = &obj
	}

	stop := false
	go func() {
		for {
			if stop {
				break
			}
			e1 := gevent.Event1001{}
			e2 := gevent.Event1002{}
			l.TriggerEvent(&e1)
			l.TriggerEvent(&e2)
		}
	}()
	time.Sleep(1 * time.Second)
	stop = true

	countTotal1 := 0
	countTotal2 := 0
	for _, obj := range objList {
		countTotal1 += obj.Count1
		countTotal2 += obj.Count2
	}
	LogInfo("1秒内事件1执行了 %d 次", countTotal1)
	LogInfo("1秒内事件2执行了 %d 次", countTotal2)
	LogInfo("TestEvent_3 success !")

}
