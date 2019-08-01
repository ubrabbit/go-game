package internal

import (
	"reflect"
	"server/leaf/gate"
	"server/leaf/network/packet"
	"server/login/command"
	"server/msg/protocol"
)

import (
	. "server/common"
)

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func handleCommand(args []interface{}) {
	defer func() {
		err := recover()
		if err != nil {
			LogError("handleCommand %x error: %v", args, err)
		}
	}()
	proto := int(args[0].(packet.Packet).Protocol())
	f, ok := g_CommandList[proto]
	if ok {
		agent := args[1].(gate.Agent)
		f(args[0], agent)
	} else {
		LogError("Unknow Login Protocol: %d", proto)
	}
}

var g_CommandList map[int]func(interface{}, gate.Agent)

func setHandler(i interface{}, f func(interface{}, gate.Agent)) {
	proto := int(i.(packet.Packet).Protocol())
	g_CommandList[proto] = f
	handler(i, handleCommand)
}

func init() {
	g_CommandList = make(map[int]func(interface{}, gate.Agent), 0)

	setHandler(&protocol.C2GSIdentity{}, command.HandleC2GSIdentity)
	setHandler(&protocol.C2GSLogin{}, command.HandleC2GSLogin)
}
