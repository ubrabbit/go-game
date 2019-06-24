package timer

import (
	. "server/common"
)

func (i *TimerItem) Repr() string {
	return FormatString("[Timer][%s][%s]%d", i.event, i.Key, i.ID)
}

func (i *TimerItem) Stop() {
	defer func() {
		i.Unlock()
		r := recover()
		if r != nil {
			LogError("timer %v stop error: %v", i, r)
		}
	}()
	i.Lock()

	if i.stopped {
		return
	}
	i.stopped = true
	i.timer.Stop()
}

/*
	Timer回调时会有这种极端情况：
	1: 定时器到期，推送入执行模块的channel等待执行
	2：与此同时，定时器所属模块执行了RemoveTimer操作
	3：因为在执行2时，1已经在channel等待执行，所以接下里定时器会被模块chanel取出执行
	这里会造成明明已经删除定时器，但仍然被执行了的BUG。
	解决方式是在回调时重新判断定时器可执行条件。
*/
func (i *TimerItem) TimerCallback() {
	defer func() {
		r := recover()
		if r != nil {
			LogError("timer %v callback error: %v", i, r)
		}
	}()

	//由其他goroutine回调的函数，在执行时需要重新判断是否可执行
	if i.stopped {
		LogInfo("%s callback but stopped", i.Repr())
		return
	}

	//加锁不能包含 chanRPC.Go 或 callback.Call
	if i.chanRPC != nil {
		i.chanRPC.Go(i.event, i)
	} else {
		go i.callback.Call()
	}
}

func (i *TimerItem) IsStop() bool {
	return i.stopped
}

func (i *TimerItem) Execute() {
	//LogDebug("Timer Execute GoroutineID: %s", GetGoroutineID())
	//由其他goroutine回调的函数，在执行时需要重新判断是否可执行
	i.Lock()
	if i.stopped {
		i.Unlock()
		LogInfo("%s execute but stopped", i.Repr())
		return
	}
	i.stopped = true
	i.Unlock()

	//Call不能加锁，因为Call可能会调用timer的函数导致死锁！
	i.callback.Call()
}
