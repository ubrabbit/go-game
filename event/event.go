package event

import (
	. "server/common"
	"server/leaf/chanrpc"
)

func RegistEventModule(module string, ch *chanrpc.Server) {
	_, ok := listenerModule[module]
	if ok {
		LogFatal("event module %s regist before!", module)
	}
	listenerModule[module] = ch
	LogInfo("regist event: %s", module)
}

func getModule(module string) *chanrpc.Server {
	ch, ok := listenerModule[module]
	if !ok {
		LogPanic("event module %s not regist!", module)
	}
	return ch
}

func init() {
	LogInfo("init event module")
	listenerModule = make(map[string]*chanrpc.Server, 0)
}
