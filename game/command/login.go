package command

import (
	"server/game/net"
	"server/game/player"
	"server/leaf/gate"
	loginclient "server/login/client"
	"server/msg/protocol"
)

import (
	. "server/common"
)

func HandleC2GSHello(i interface{}, agent gate.Agent) {
	m := i.(*protocol.C2GSHello)
	LogDebug("C2GSHello %v", m.Seed)
	net.OnHello(agent, m)
}

func HandleC2GSIdentity(i interface{}, agent gate.Agent) {
	m := i.(*protocol.C2GSIdentity)
	LogDebug("C2GSIdentity %v", m.Identity)
	net.OnIdentity(agent, m)
}

func HandleC2GSLoadRole(i interface{}, agent gate.Agent) {
	m := i.(*protocol.C2GSLoadRole)
	pid := m.Pid
	LogDebug("C2GSLoadRole %v", pid)

	addr := agent.RemoteAddr().String()
	client := loginclient.GetLoginClient(addr)
	if client == nil || !client.IsAlive() {
		LogError("%s load %d but client is dead", addr, pid)
		return
	}
	//登陆成功才允许加载角色
	if !client.IsLogin() {
		LogError("%s load %d but login failure", addr, pid)
		return
	}
	//验证该连接是否拥有这些角色
	if !client.ValidPlayerEnter(pid) {
		LogError("%s load %d but has no player", addr, pid)
		return
	}

	id := net.LoginSuccess(pid, agent)
	if id <= 0 {
		LogPanic("player %d login but clientID error: %v", pid, id)
		return
	}
	p, err := player.Login(pid, id)
	if err != nil {
		LogPanic("player %d login error: %v", pid, err)
	}
	agent.WriteMsg(&protocol.GS2CLoadRole{
		NLen: len(p.Name),
		Name: p.Name,
	})
	agent.WriteMsg(&protocol.GS2CLoginFinished{})
}

func HandleC2GSLoginFinished(i interface{}, agent gate.Agent) {
	LogDebug("C2GSLoginFinished")
}
