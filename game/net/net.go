package net

import (
	"errors"
	"fmt"
	"server/game/player"
	"server/leaf/gate"
	"server/leaf/network"
	"server/leaf/network/packet"
)

import (
	. "server/common"
)

func getAgentClientID(a gate.Agent) int {
	addr := a.RemoteAddr().String()
	id, ok := g_Container.clients[addr]
	//LogInfo("getAgentClientID: %s %v", addr, id)
	if ok {
		return id
	}
	return -1
}

func (n *NetContainer) clean(a gate.Agent) {
	addr := a.RemoteAddr().String()
	id, ok := g_Container.clients[addr]
	if ok {
		delete(n.agents, id)
		delete(n.players, id)
	}
	delete(g_Container.clients, addr)
}

func (n *NetContainer) doDisconnect(a gate.Agent) {
	id := getAgentClientID(a)
	LogInfo("doDisconnect: %d", id)
	n.clean(a)
	a.Close()
}

// 主动断开连接
func (n *NetContainer) Disconnect(a gate.Agent) {
	defer func() {
		n.Unlock()
		r := recover()
		if r != nil {
			LogError("%v OnDisconnect Error: %v", a, r)
		}
	}()
	n.Lock()
	n.doDisconnect(a)
}

func (n *NetContainer) OnConnect(a gate.Agent) {
	defer func() {
		n.Unlock()
		r := recover()
		if r != nil {
			LogError("%v OnConnect Error: %v", a, r)
		}
	}()
	n.Lock()
	g_Container.clientIdx++
	id := g_Container.clientIdx
	addr := a.RemoteAddr().String()
	// 除非出现故障，否则不可能会出现两个相同外部地址连进来的情况。 如果真的发生，说明旧连接没处理掉，需要主动处理。
	id2, exists := n.clients[addr]
	if exists {
		a2, exists := n.agents[id2]
		if exists {
			LogInfo("Close Old Address: %s %d", addr, id2)
			n.clean(a2)
			a2.Close()
		}
	}
	LogInfo("Connect: %s %d", addr, id)
	n.agents[id] = a
	n.clients[addr] = id
}

func (n *NetContainer) OnDisconnect(a gate.Agent) {
	defer func() {
		n.Unlock()
		r := recover()
		if r != nil {
			LogError("%v OnDisconnect Error: %v", a, r)
		}
	}()
	n.Lock()
	id := getAgentClientID(a)
	pid, ok := n.players[id]
	if ok {
		LogInfo("Disconnect: %d(%d)", id, pid)
		player.Logout(pid)
	} else {
		LogInfo("Disconnect: %d", id)
	}
	n.clean(a)
}

// 登陆完成
func (n *NetContainer) OnPlayerLogin(pid int, a gate.Agent) (id int) {
	defer func() {
		n.Unlock()
		r := recover()
		if r != nil {
			LogError("%d OnPlayerLogin error: %v", getAgentClientID(a), r)
			id = 0
		}
	}()
	n.Lock()
	id = getAgentClientID(a)
	_, ok := n.agents[id]
	if !ok {
		LogPanic("%d OnPlayerLogin, but agent not exists in clients!", pid)
	}
	pid2, ok := n.players[id]
	//同一个连接加载了不同的角色，就顶前一个角色下线
	if ok && pid != pid2 {
		player.Logout(pid2)
	}
	n.players[id] = pid
	return id
}

func RpcNewAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	g_Container.OnConnect(a)
}

func RpcCloseAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	g_Container.OnDisconnect(a)
}

func LoginSuccess(pid int, a gate.Agent) int {
	LogInfo("LoginSuccess: %d", pid)
	return g_Container.OnPlayerLogin(pid, a)
}

func PacketSend(id int, proto uint8, msgData []byte) (err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = r.(error)
		}
	}()

	g_Container.Lock()
	agent, ok := g_Container.agents[id]
	g_Container.Unlock()
	if !ok {
		return errors.New(fmt.Sprintf("client %d has no agent!", id))
	}
	msg := network.PacketProto(proto, msgData)
	agent.WriteMsg(msg)
	return nil
}

func SendToPlayer(p *player.Player, i interface{}) (err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = r.(error)
		}
	}()
	if p == nil {
		return errors.New("SendToPlayer but player is nil")
	}
	clientID := p.GetClientID()
	if p.GetClientID() <= 0 {
		return errors.New(fmt.Sprintf("player %d has no clientID: %d", p.Pid, clientID))
	}
	proto, msgData := i.(packet.Packet).PacketData()
	err = PacketSend(clientID, proto, msgData)
	if err != nil {
		return errors.New(fmt.Sprintf("SendToPlayer %d %d error: %v", p.Pid, proto, err))
	}
	return nil
}

func init() {
	g_Container = new(NetContainer)
	g_Container.clientIdx = 1000
	g_Container.clients = make(map[string]int, 0)
	g_Container.agents = make(map[int]gate.Agent, 0)
	g_Container.players = make(map[int]int, 0)
}
