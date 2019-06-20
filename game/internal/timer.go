package internal

import (
	"server/timer"
)
import (
	. "server/common"
)

func GameTimerCallback(args []interface{}) {
	defer func() {
		r := recover()
		if r != nil {
			LogError("Callback Error: %v %v", args, r)
		}
	}()
	//LogDebug("GameTimerCallback")
	f := args[0].(Functor)
	f.Call()
}

func init() {
	LogInfo("init game timer")
	timer.RegistModule(timer.TIMER_MODULE_GAME, timer.TIMER_EVENT_GAME, ChanRPC)
	skeleton.RegisterChanRPC(timer.TIMER_EVENT_GAME, GameTimerCallback)
}
