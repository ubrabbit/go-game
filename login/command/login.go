package command

import (
	"server/leaf/gate"
	loginclient "server/login/client"
	"server/msg/protocol"
)

import (
	. "server/common"
)

func HandleC2GSIdentity(i interface{}, agent gate.Agent) {
	m := i.(*protocol.C2GSIdentity)
	LogDebug("C2GSIdentity %v", m.Identity)
	loginclient.OnIdentity(agent, m)
	agent.WriteMsg(&protocol.GS2CIdentity{})
}

//login模块只验证帐号密码和返回角色列表，角色登陆进游戏在game模块（属于另一个goroutine）
func HandleC2GSLogin(i interface{}, agent gate.Agent) {
	m := i.(*protocol.C2GSLogin)
	LogDebug("HandleC2GSLogin: %s : %s", m.User, m.Password)

	client := loginclient.ClientLogin(m.User, m.Password, agent)
	if client == nil {
		agent.Destroy()
		return
	}
	loginT, pid := 0, 0
	if client.IsLogin() {
		loginT = 1
		pid = client.GetPlayerID()
	}
	agent.WriteMsg(&protocol.GS2CLogin{Type: loginT, Pid: pid})
}
