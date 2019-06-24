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
	i := args[0].(*timer.TimerItem)
	i.Execute()
}

func init() {
	LogInfo("init game timer")
	timer.RegistModule(timer.TIMER_MODULE_GAME, ChanRPC)
	skeleton.RegisterChanRPC(timer.TIMER_MODULE_GAME, GameTimerCallback)
}
