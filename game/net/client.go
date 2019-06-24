package net

import (
	"server/game/player"
	"server/leaf/gate"
	"server/timer"
)

import (
	. "server/common"
)

func newClient(a gate.Agent) *NetClient {
	c := NetClient{
		id:    NewObjectID(),
		agent: a,
	}
	c.playerID = 0
	c.Create()
	return &c
}

func (c *NetClient) ID() int {
	return c.id
}

func (c *NetClient) PlayerID() int {
	return c.playerID
}

func (c *NetClient) Create() {
	timer.AddObject(c.ID())
}

func (c *NetClient) Delete() {
	timer.RemoveObject(c.ID())
}

func (c *NetClient) RemoteAddr() string {
	return c.agent.RemoteAddr().String()
}

func (c *NetClient) Repr() string {
	return FormatString("%s(%d)", c.RemoteAddr(), c.ID())
}

func (c *NetClient) Close() {
	c.agent.Close()
}

func (c *NetClient) OnConnect() {}

func (c *NetClient) OnDisconnect() (id int) {
	defer func() {
		c.Unlock()
		r := recover()
		if r != nil {
			LogError("%s OnDisconnect error: %v", c.Repr(), r)
			id = 0
		}
	}()
	c.Lock()

	player.Logout(c.playerID)
	return c.id
}

// 登陆完成
func (c *NetClient) LoginSuccess(pid int) (id int) {
	defer func() {
		c.Unlock()
		r := recover()
		if r != nil {
			LogError("%s LoginSuccess %d error: %v", c.Repr(), pid, r)
			id = 0
		}
	}()
	c.Lock()

	pid2 := c.playerID
	//同一个连接加载了不同的角色，就顶前一个角色下线
	if pid != pid2 {
		player.Logout(pid2)
	}
	c.playerID = pid
	LogInfo("%s LoginSuccess: %d", c.Repr(), pid)
	return c.id
}

func (c *NetClient) PacketSend(i interface{}) (err error) {
	defer func() {
		r := recover()
		if r != nil {
			LogError("%s PacketSend %v error: %v", c.Repr(), i, r)
			err = r.(error)
		}
	}()
	c.agent.WriteMsg(i)
	return nil
}
