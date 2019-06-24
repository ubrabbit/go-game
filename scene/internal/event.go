package internal

import (
	"server/event"
)
import (
	. "server/common"
)

func SceneEventCallback(args []interface{}) {
	defer func() {
		r := recover()
		if r != nil {
			LogError("Callback Error: %v %v", args, r)
		}
	}()
	//LogDebug("SceneEventCallback")
	f := args[0].(*Functor)
	f.Call(args[1])
}

func init() {
	LogInfo("init scene event")

	event.RegistEventModule(event.EVENT_MODULE_SCENE, ChanRPC)
	skeleton.RegisterChanRPC(event.EVENT_MODULE_SCENE, SceneEventCallback)
}
