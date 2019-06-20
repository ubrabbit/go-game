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

func init() {
	g_CommandList = make(map[int]func(interface{}, gate.Agent), 0)
	g_CommandList[protocol.C2GS_HELLO] = command.HandleC2GSHello
	g_CommandList[protocol.C2GS_IDENTIFY] = command.HandleC2GSIdentity
	g_CommandList[protocol.C2GS_LOGIN] = command.HandleC2GSLogin

	handler(&protocol.C2GSHello{}, handleCommand)
	handler(&protocol.C2GSIdentity{}, handleCommand)
	handler(&protocol.C2GSLogin{}, handleCommand)
}
