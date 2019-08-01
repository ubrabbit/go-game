package client

import (
	"gopkg.in/mgo.v2/bson"
	db "server/db/mongodb"
	"server/db/table"
	"server/leaf/gate"
	"server/timer"
)

import (
	. "server/common"
	. "server/msg/protocol"
)

func (c *LoginClient) ID() int {
	return c.id
}

func (c *LoginClient) Create() {
	timer.AddObject(c.ID())
}

func (c *LoginClient) Delete() {
	timer.RemoveObject(c.ID())
}

func (c *LoginClient) Repr() string {
	return FormatString("%s(%d)", c.RemoteAddr(), c.ID())
}

func (c *LoginClient) RemoteAddr() string {
	return c.agent.RemoteAddr().String()
}

func (c *LoginClient) IsLogin() bool {
	return c.authSuccess
}

func (c *LoginClient) IsNew() bool {
	return c.isNew
}

func (c *LoginClient) IsAlive() bool {
	return GetSecond()-c.connectTime <= constClientAlive
}

func (c *LoginClient) GetPlayerID() int {
	if len(c.playerList) > 0 {
		for pid, _ := range c.playerList {
			return pid
		}
	}
	return 0
}

func (c *LoginClient) HasPlayer(pid int) bool {
	_, exists := c.playerList[pid]
	return exists
}

func (c *LoginClient) ValidPlayerEnter(pid int) bool {
	if !c.IsAlive() {
		return false
	}
	if !c.HasPlayer(pid) {
		return false
	}
	return true
}

func (c *LoginClient) loginTimeout(args ...interface{}) {
	defer func() {
		g_Lock.Unlock()
		r := recover()
		if r != nil {
			LogError("LoginClient %s timeout err: %v", c.Repr(), r)
		}
	}()
	g_Lock.Lock()

	cleanClient(c)
}

func createPlayer(acct string) int {
	c := db.GetServerC(table.DBNamePlayer)
	pid := db.NewPlayerID()
	itemPlayer := table.DBTablePlayer{
		ServerNum: GetServerNum(),
		Account:   acct,
		Pid:       pid,
		Name:      acct,
		Grade:     1,
	}
	LogInfo("%s create player %d", acct, pid)
	err := itemPlayer.Insert(c)
	if err != nil {
		LogError("%s create player %d error: %v", acct, pid, err)
		return 0
	}
	return pid
}

func OnHello(agent gate.Agent) {
	client := GetLoginClient(agent.RemoteAddr().String())
	if client == nil {
		agent.Destroy()
		return
	}
	client.helloSuccess = true
}

func OnIdentity(agent gate.Agent, p *C2GSIdentity) {
	client := GetLoginClient(agent.RemoteAddr().String())
	if client == nil {
		agent.Destroy()
		return
	}
	client.identitySuccess = true
}

//这个函数禁止加锁
func ClientLogin(acct string, pwd string, agent gate.Agent) *LoginClient {
	client := GetLoginClient(agent.RemoteAddr().String())
	if client == nil {
		LogError("%s try login but client is dead", agent.RemoteAddr().String())
		return nil
	}
	if !client.helloSuccess {
		LogError("%s try login but hello failed", client.Repr())
		return nil
	}
	if !client.identitySuccess {
		LogError("%s try login but identity failed", client.Repr())
		return nil
	}
	//符合登陆请求时才从堆中分配玩家列表容器
	client.playerList = make(map[int]LoginPlayer, 0)

	item := table.DBTableAccount{}
	c := db.GetServerC(table.DBNameAccount)
	c2 := db.GetServerC(table.DBNamePlayer)
	err := c.Find(bson.M{"account": acct}).One(&item)
	if err == db.ErrNotFound { //注册
		client.isNew = true

		item.Account = acct
		item.Password = pwd
		item.Email = acct
		LogInfo("regist %s", acct)
		err = c.Insert(&item)
		if err != nil {
			LogError("regist %s error: %v", acct, err)
			return nil
		}
		client.authSuccess = true
	} else { //验证
		if item.Password != pwd {
			client.authSuccess = false
		} else {
			client.authSuccess = true
		}
	}
	if !client.authSuccess {
		return client
	}

	pid := 0
	//需要创建新角色
	if len(item.Players) <= 0 {
		pid = createPlayer(acct)
		if pid <= 0 {
			return nil
		}
		item.Players = append(item.Players, pid)
		item.Save(c)
	}
	for i := 0; i < 5; i++ {
		for _, pid := range item.Players {
			itemPlayer := table.DBTablePlayer{
				ServerNum: GetServerNum(),
				Account:   acct,
				Pid:       pid,
				Name:      acct,
				Grade:     1,
			}
			err = itemPlayer.Load(c2)
			if err != nil {
				LogError("Account %s Load Player %d error: %v", acct, pid, err)
				continue
			}
			obj := LoginPlayer{
				ServerNum: itemPlayer.ServerNum,
				Pid:       itemPlayer.Pid,
				Name:      itemPlayer.Name,
				Grade:     itemPlayer.Grade,
			}
			client.playerList[pid] = obj
		}
		if len(client.playerList) > 0 {
			//如果有了能登陆的角色，就跳出循环
			break
		} else {
			/* 当player与account表出现数据不同步时，可能会出现有无法登陆的角色的情况。
			此时创建一个新角色给这个帐号。*/
			pid = createPlayer(acct)
			if pid > 0 {
				item.Players = append(item.Players, pid)
				item.Save(c)
			}
		}
	}

	//到这一步，这个帐号都还没有可登陆角色，说明出异常了！
	if len(client.playerList) <= 0 {
		LogError("fail to load account %s players create new one", acct)
		return nil
	}
	LogDebug("ClientLogin: %s %v", acct, client.playerList)
	return client
}

func cleanClient(c *LoginClient) {
	LogDebug("clean LoginClient: %s", c.Repr())
	c.Delete()
	delete(g_LoginClientList, c.RemoteAddr())
}

func OnClientConnect(agent gate.Agent) (err error) {
	defer func() {
		g_Lock.Unlock()
		r := recover()
		if r != nil {
			err = r.(error)
		}
	}()
	g_Lock.Lock()

	//创建连接对象
	c := &LoginClient{
		id:              NewObjectID(),
		Account:         "",
		Password:        "",
		agent:           agent,
		isNew:           false,
		connectTime:     GetSecond(),
		helloSuccess:    false,
		identitySuccess: false,
		authSuccess:     false,
	}
	g_LoginClientList[c.RemoteAddr()] = c
	c.Create()
	f := NewFunctor("LoginTimeout", c.loginTimeout)
	timer.StartTimer(timer.TIMER_MODULE_LOGIN, c.ID(), "LoginTimeout", constClientAlive, f)
	return nil
}

func OnClientDisconnect(agent gate.Agent) (err error) {
	defer func() {
		g_Lock.Unlock()
		r := recover()
		if r != nil {
			err = r.(error)
		}
	}()
	g_Lock.Lock()

	addr := agent.RemoteAddr().String()
	c, ok := g_LoginClientList[addr]
	if ok {
		cleanClient(c)
	}
	return nil
}

func GetLoginClient(addr string) *LoginClient {
	defer g_Lock.Unlock()
	g_Lock.Lock()
	c, _ := g_LoginClientList[addr]
	return c
}

func init() {
	g_LastCleanTime = GetSecond()
	g_LoginClientList = make(map[string]*LoginClient, 0)
}
