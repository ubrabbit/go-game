package timer

import (
	. "server/common"
)

func (i *TimerItem) Stop() {
	defer func() {
		r := recover()
		if r != nil {
			LogError("timer %v stop error: %v", i, r)
		}
	}()
	i.timer.Stop()
}

func (i *TimerItem) Callback() {
	defer func() {
		r := recover()
		if r != nil {
			LogError("timer %v callback error: %v", i, r)
		}
	}()

	if i.chanRPC != nil {
		i.chanRPC.Go(i.event, i.callback)
	} else {
		go i.callback.Call()
	}
}
