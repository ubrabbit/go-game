package common

import (
	. "server/common"
	"server/leaf/network/packet"
	. "server/msg/protocol"
	"sync"
)

type LoginClient struct {
	sync.Mutex
	Client *Client
}

const (
	helloSeed   = 1234567890
	identityKey = "demo_identity"
)

func NewLoginClient() *LoginClient {
	c := LoginClient{
		Client: NewClient(),
	}
	return &c
}

func (c *LoginClient) PacketSend(p packet.Packet) {
	client := c.Client
	proto, msg := p.PacketData()
	client.SendData = msg
	client.SendProto(proto)
}

func (c *LoginClient) Login(user string, password string) int {
	c.C2GSHelo()
	c.GS2CHelo()
	c.C2GSIdentity()
	c.GS2CIdentity()
	c.C2GSLogin(user, password)
	pid := c.GS2CLogin()
	c.C2GSLoadRole(pid)
	name := c.GS2LoadRole()
	if user != name {
		LogPanic("GS2LoadRole Failure, '%s' != '%s'", user, name)
	}
	c.GS2CLoginFinished()
	return pid
}

//未按流程发包的登陆，会被服务器直接断线
func (c *LoginClient) LoginNoHello(user string, password string) int {
	defer func() {
		r := recover()
		if r != nil {
			err := r.(error)
			if err.Error() != "EOF" {
				LogFatal("LoginNoHello Fail!: %v", err)
			}
		}
	}()
	c.C2GSIdentity()
	c.GS2CIdentity()
	c.C2GSLogin(user, password)
	pid := c.GS2CLogin()
	c.C2GSLoadRole(pid)
	name := c.GS2LoadRole()
	if user != name {
		LogPanic("GS2LoadRole Failure, '%s' != '%s'", user, name)
	}
	c.GS2CLoginFinished()
	return pid
}

//未按流程发包的登陆，会被服务器直接断线
func (c *LoginClient) LoginNoIdentity(user string, password string) int {
	defer func() {
		r := recover()
		if r != nil {
			err := r.(error)
			if err.Error() != "EOF" {
				LogFatal("LoginNoIdentity Fail!: %v", err)
			}
		}
	}()
	c.C2GSHelo()
	c.GS2CHelo()
	c.C2GSLogin(user, password)
	pid := c.GS2CLogin()
	c.C2GSLoadRole(pid)
	name := c.GS2LoadRole()
	if user != name {
		LogPanic("GS2LoadRole Failure, '%s' != '%s'", user, name)
	}
	c.GS2CLoginFinished()
	return pid
}

func (c *LoginClient) C2GSHelo() {
	p := &C2GSHello{
		Seed: helloSeed,
	}
	c.PacketSend(p)
}

func (c *LoginClient) GS2CHelo() {
	proto, msg, err := c.Client.UnpackProto()
	CheckPanic(err)

	p := &GS2CHello{}
	p.UnpackData(msg)
	if proto != p.Protocol() || p.Seed != helloSeed {
		LogPanic("GS2CHelo Fail: %d %v %v", proto, msg, err)
	}
}

func (c *LoginClient) C2GSIdentity() {
	p := &C2GSIdentity{
		Identity: identityKey,
	}
	c.PacketSend(p)
}

func (c *LoginClient) GS2CIdentity() {
	proto, msg, err := c.Client.UnpackProto()
	CheckPanic(err)

	p := &GS2CIdentity{}
	p.UnpackData(msg)
	if proto != p.Protocol() {
		LogPanic("GS2CIdentity Fail: %d %v %v", proto, msg, err)
	}
}

func (c *LoginClient) C2GSLogin(user string, password string) {
	p := &C2GSLogin{
		User:     user,
		Password: password,
	}
	c.PacketSend(p)
}

func (c *LoginClient) GS2CLogin() int {
	proto, msg, err := c.Client.UnpackProto()
	CheckPanic(err)

	p := &GS2CLogin{}
	p.UnpackData(msg)
	if proto != p.Protocol() {
		LogPanic("GS2CLogin Fail: %d %v %v", proto, p, err)
	}
	if p.Pid <= 0 {
		LogPanic("GS2CLogin Fail, p.Pid <= 0: %d %v %v", proto, p, err)
	}
	return p.Pid
}

func (c *LoginClient) C2GSLoadRole(pid int) {
	p := &C2GSLoadRole{
		Pid: pid,
	}
	c.PacketSend(p)
}

func (c *LoginClient) GS2LoadRole() string {
	proto, msg, err := c.Client.UnpackProto()
	CheckPanic(err)

	p := &GS2CLoadRole{}
	p.UnpackData(msg)
	if proto != p.Protocol() {
		LogPanic("GS2LoadRole Fail: %d %v %v", proto, msg, err)
	}
	return p.Name
}

func (c *LoginClient) GS2CLoginFinished() {
	proto, msg, err := c.Client.UnpackProto()
	CheckPanic(err)

	p := &GS2CLoginFinished{}
	if proto != p.Protocol() {
		LogPanic("GS2CLoginFinished Fail: %d %v %v", proto, msg, err)
	}
}
