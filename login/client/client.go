package client

import (
	"gopkg.in/mgo.v2/bson"
	db "server/db/mongodb"
	"server/db/table"
	"server/leaf/gate"
)

import (
	. "server/common"
)

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

func ClientHello(agent gate.Agent) {}

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

//这个函数禁止加锁
func ClientLogin(acct string, pwd string, agent gate.Agent) *LoginClient {
	//创建连接对象
	client := LoginClient{
		Account:     acct,
		Password:    pwd,
		agent:       agent,
		authSuccess: false,
		isNew:       false,
		connectTime: GetSecond(),
	}
	client.authSuccess = false
	client.playerList = make(map[int]LoginPlayer, 0)
	err := OnClientConnect(&client)
	if err != nil {
		return nil
	}
	item := table.DBTableAccount{}
	c := db.GetServerC(table.DBNameAccount)
	c2 := db.GetServerC(table.DBNamePlayer)
	err = c.Find(bson.M{"account": acct}).One(&item)
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
	} else { //验证
		if item.Password != pwd {
			return nil
		}
	}
	client.authSuccess = true

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
	return &client
}

//每隔一段时间在新连接时触发。不使用定时器处理，因为用锁复杂度太高！
func cleanOldClients() {
	defer func() {
		r := recover()
		if r != nil {
			LogError("CleanOldClients error: %v", r)
		}
	}()
	now := GetSecond()
	if now-g_LastCleanTime <= constCleanInternal {
		return
	}
	g_LastCleanTime = now
	for addr, obj := range g_LoginClientList {
		if now-obj.connectTime >= constClientAlive {
			delete(g_LoginClientList, addr)
		}
	}
}

func OnClientConnect(c *LoginClient) (err error) {
	defer func() {
		g_Lock.Unlock()
		r := recover()
		if r != nil {
			err = r.(error)
		}
	}()
	g_Lock.Lock()

	addr := c.agent.RemoteAddr().String()
	g_LoginClientList[addr] = c
	cleanOldClients()
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
