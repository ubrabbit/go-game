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
		obj.timerList = make(map[string]*TimerItem, 0)
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
	for _, i := range c.timerList {
		i.Stop()
	}
	delete(g_Container, id)
	return true
}

func startTimer(mod *TimerModule, c *TimerObject, key string, timeout int, f *Functor) {
	defer c.Unlock()
	c.Lock()
	i2, ok := c.timerList[key]
	//旧定时器被顶掉了
	if ok {
		i2.Stop()
	}
	if timeout <= 0 {
		timeout = 1
	}
	i := &TimerItem{
		ID:       NewObjectID(),
		Key:      key,
		event:    mod.event,
		chanRPC:  mod.chanRPC,
		callback: f,
		stopped:  false,
	}
	i.timer = mod.wheel.AfterFunc(time.Duration(timeout)*time.Millisecond, i.TimerCallback)
	c.timerList[key] = i
}

func StartTimer(name string, id int, key string, timeout int, f *Functor) {
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

func StartTimerMs(name string, id int, key string, timeout int, f *Functor) {
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
	i, ok := c.timerList[key]
	if !ok {
		return false
	}
	i.Stop()
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

func RegistModule(name string, ch *chanrpc.Server) {
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
		event:   name,
		wheel:   wheel.NewTimingWheel(time.Millisecond, timerQueueDefaultLen),
		chanRPC: ch,
	}
	tw := g_Module[name].wheel
	go tw.Start()
	LogInfo("RegistModule %s", name)
}

func init() {
	LogInfo("init timer")
	g_Module = make(map[string]*TimerModule, 0)
	g_Container = make(map[int]*TimerObject, 0)
}
