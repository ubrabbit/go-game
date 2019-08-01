package net

import (
	"errors"
	"server/game/player"
	"server/leaf/gate"
	loginclient "server/login/client"
	"server/timer"
)

import (
	. "server/common"
	. "server/msg/protocol"
)

func (n *NetContainer) getAgentClientID(a gate.Agent) int {
	addr := a.RemoteAddr().String()
	c, ok := n.clients[addr]
	if ok {
		return c.ID()
	}
	return -1
}

func (n *NetContainer) getClient(a gate.Agent) *NetClient {
	addr := a.RemoteAddr().String()
	c, ok := n.clients[addr]
	if ok {
		return c
	}
	return nil
}

func (n *NetContainer) getClientByID(id int) *NetClient {
	c, ok := n.clientsID[id]
	if ok {
		return c
	}
	return nil
}

func (n *NetContainer) getClientByPid(pid int) *NetClient {
	c, ok := n.playersID[pid]
	if ok {
		return c
	}
	return nil
}

func (n *NetContainer) clean(c *NetClient) {
	addr := c.RemoteAddr()
	delete(n.clients, addr)
	delete(n.clientsID, c.ID())
	delete(n.playersID, c.PlayerID())
	c.Delete()
	c.Close()
	LogDebug("clean NetClient: %s", c.Repr())
}

// 主动断开连接
func (n *NetContainer) Disconnect(c *NetClient) {
	defer func() {
		n.Unlock()
		r := recover()
		if r != nil {
			LogError("%s Disconnect Error: %v", c.Repr(), r)
		}
	}()
	n.Lock()

	LogInfo("%s Disconnect", c.Repr())
	c.OnDisconnect()
	n.clean(c)
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

	c := newClient(a)
	addr := c.RemoteAddr()
	// 除非出现故障，否则不可能会出现两个相同外部地址连进来的情况。 如果真的发生，说明旧连接没处理掉，需要主动处理。
	c2, exists := n.clients[addr]
	if exists {
		LogInfo("Close Old Address: %s", c2.Repr())
		n.clean(c2)
	}
	id := c.ID()
	LogInfo("%s Connect ", c.Repr())
	n.clients[addr] = c
	n.clientsID[id] = c
	c.OnConnect()

	f2 := NewFunctor("LoginTimeout", n.loginTimeout, c.ID())
	timer.StartTimer(timer.TIMER_MODULE_GAME, c.ID(), "LoginTimeout", connectLoginTimeout, f2)
}

func (n *NetContainer) helloTimeout(args ...interface{}) {
	id := args[0].(int)
	c := n.getClientByID(id)
	if c == nil {
		return
	}
	LogInfo("%s helloTimeout", c.Repr())
	n.Disconnect(c)
}

func (n *NetContainer) loginTimeout(args ...interface{}) {
	id := args[0].(int)
	c := n.getClientByID(id)
	if c == nil {
		return
	}
	LogInfo("%s loginTimeout", c.Repr())
	n.Disconnect(c)
}

func (n *NetContainer) OnHello(a gate.Agent) {
	c := n.getClient(a)
	if c != nil {
		f1 := NewFunctor("HelloTimeout", n.helloTimeout, c.ID())
		timer.StartTimer(timer.TIMER_MODULE_GAME, c.ID(), "HelloTimeout", connectHelloTimeout, f1)
		c.PacketSend(&GS2CHello{Seed: 1234567890})
	} else {
		LogError("%s hello but client not exists", a.RemoteAddr().String())
		a.Destroy()
	}
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

	c := n.getClient(a)
	if c != nil {
		LogInfo("%s OnDisconnect", c.Repr())
		c.OnDisconnect()
		n.clean(c)
	} else {
		LogInfo("%s(-1) OnDisconnect", a.RemoteAddr().String())
	}
}

// 登陆完成
func (n *NetContainer) LoginSuccess(pid int, a gate.Agent) (id int) {
	defer func() {
		n.Unlock()
		r := recover()
		if r != nil {
			LogError("%d LoginSuccess error: %v", n.getAgentClientID(a), r)
			id = 0
		}
	}()
	n.Lock()

	c := n.getClient(a)
	if c == nil {
		LogPanic("%d LoginSuccess, but agent not exists in clients!", pid)
	}
	c2 := n.getClientByPid(pid)
	//顶掉旧的客户端
	if c2 != nil && c2.ID() != c.ID() {
		LogInfo("Close Old Client: %s", c2.Repr())
		n.clean(c2)
	}
	timer.RemoveTimer(timer.TIMER_MODULE_GAME, c.ID(), "HelloTimeout")
	timer.RemoveTimer(timer.TIMER_MODULE_GAME, c.ID(), "LoginTimeout")
	c.LoginSuccess(pid)
	return c.ID()
}

//网络连接入口
func RpcNewAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	g_Container.OnConnect(a)
	err := loginclient.OnClientConnect(a)
	if err != nil {
		LogError("OnConnect error: %v", err)
		g_Container.OnDisconnect(a)
		loginclient.OnClientDisconnect(a)
	} else {
		g_Container.OnHello(a)
		loginclient.OnHello(a)
	}
}

//网络断开连接入口
func RpcCloseAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	g_Container.OnDisconnect(a)
	loginclient.OnClientDisconnect(a)
}

func LoginSuccess(pid int, a gate.Agent) int {
	return g_Container.LoginSuccess(pid, a)
}

func SendToPlayer(p *player.Player, i interface{}) (err error) {
	defer func() {
		g_Container.Unlock()
		r := recover()
		if r != nil {
			err = r.(error)
		}
	}()
	g_Container.Lock()

	if p == nil {
		return errors.New("SendToPlayer but player is nil")
	}
	clientID := p.GetClientID()
	c := g_Container.getClientByID(clientID)
	if c == nil {
		return errors.New(FormatString("SendToPlayer %d but client not exists, clientID: %d", p.Pid, clientID))
	}

	err = c.PacketSend(i)
	if err != nil {
		return errors.New(FormatString("SendToPlayer %d %v error: %v", p.Pid, i, err))
	}
	return nil
}

func init() {
	g_Container = new(NetContainer)
	g_Container.clients = make(map[string]*NetClient, 0)
	g_Container.clientsID = make(map[int]*NetClient, 0)
	g_Container.playersID = make(map[int]*NetClient, 0)
}
