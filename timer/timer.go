package timer

import (
	"server/leaf/chanrpc"
	"server/timer/wheel"
	"time"
)

import (
	. "server/common"
)

func AddObject(id int) {
	defer func() {
		g_Lock.Unlock()
	}()
	g_Lock.Lock()

	_, ok := g_Container[id]
	if !ok {
		obj := &TimerObject{
			ID: id,
		}
		obj.timerList = make(map[string][]*TimerItem, 0)
		g_Container[id] = obj
	}
}

func GetObject(id int) *TimerObject {
	defer func() {
		g_Lock.Unlock()
	}()
	g_Lock.Lock()

	obj, ok := g_Container[id]
	if !ok {
		return nil
	}
	return obj
}

func RemoveObject(id int) bool {
	defer func() {
		g_Lock.Unlock()
	}()
	g_Lock.Lock()

	c, ok := g_Container[id]
	if !ok {
		return false
	}
	for _, list := range c.timerList {
		for _, i := range list {
			i.Stop()
		}
	}
	delete(g_Container, id)
	return true
}

func startTimer(mod *TimerModule, c *TimerObject, key string, timeout int, f Functor) {
	defer c.Unlock()
	c.Lock()
	_, ok := c.timerList[key]
	if !ok {
		c.timerList[key] = make([]*TimerItem, 0)
	}
	i := &TimerItem{
		Key:      key,
		event:    mod.event,
		chanRPC:  mod.chanRPC,
		callback: f,
	}
	i.timer = mod.wheel.AfterFunc(time.Duration(timeout)*time.Millisecond, i.Callback)
	c.timerList[key] = append(c.timerList[key], i)
}

func StartTimer(name string, id int, key string, timeout int, f Functor) {
	c := GetObject(id)
	if c == nil {
		return
	}
	mod := getModule(name)
	if mod == nil {
		return
	}
	startTimer(mod, c, key, timeout*1000, f)
}

func StartTimerMs(name string, id int, key string, timeout int, f Functor) {
	c := GetObject(id)
	if c == nil {
		return
	}
	mod := getModule(name)
	if mod == nil {
		return
	}
	startTimer(mod, c, key, timeout, f)
}

func RemoveTimer(name string, id int, key string) bool {
	c := GetObject(id)
	if c == nil {
		return false
	}

	defer c.Unlock()
	c.Lock()
	_, ok := c.timerList[key]
	if !ok {
		return false
	}
	for _, i := range c.timerList[key] {
		i.Stop()
	}
	delete(c.timerList, key)
	return true
}

func getModule(name string) *TimerModule {
	mod, ok := g_Module[name]
	if ok {
		return mod
	}
	LogFatal("timer module %s not exists", name)
	return nil
}

func RegistModule(name string, event string, ch *chanrpc.Server) {
	defer func() {
		g_Lock.Unlock()
		r := recover()
		if r != nil {
			LogFatal("timer module %s regist error: %v", name, r)
		}
	}()
	g_Lock.Lock()

	_, ok := g_Module[name]
	if ok {
		LogFatal("timer module %s has regist before!", name)
	}
	g_Module[name] = &TimerModule{
		Name:    name,
		event:   event,
		wheel:   wheel.NewTimingWheel(time.Millisecond, timerQueueDefaultLen),
		chanRPC: ch,
	}
	tw := g_Module[name].wheel
	go tw.Start()
	LogInfo("RegistModule %s(%s)", name, event)
}

func init() {
	LogInfo("init timer")
	g_Module = make(map[string]*TimerModule, 0)
	g_Container = make(map[int]*TimerObject, 0)
}
