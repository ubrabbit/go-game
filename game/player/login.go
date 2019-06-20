package player

import (
	. "server/common"
)

func Login(pid int, cid int) (p *Player, err error) {
	defer func() {
		g_Container.Unlock()
		r := recover()
		if r != nil {
			err = r.(error)
		}
	}()
	g_Container.Lock()
	p, ok := g_Container.playerList[pid]
	if !ok {
		p = newPlayer(pid)
		p.clientID = cid
		g_Container.playerList[pid] = p
		p.Create()
		p.Load()
		p.OnLogin()
	} else {
		p.checkLoad()
		p.ReEnter(cid)
	}
	p.heartbeat()
	return p, nil
}

func Logout(pid int) {
	defer func() {
		g_Container.Unlock()
		r := recover()
		if r != nil {
			LogError("Player %d Logout Error: %v", pid, r)
		}
	}()
	g_Container.Lock()
	p, ok := g_Container.playerList[pid]
	if ok {
		p.OnDisconnect()
		p.Delete()
		p.clientID = 0
		delete(g_Container.playerList, pid)
	}
}

func (p *Player) ReEnter(cid int) {
	LogInfo("%d ReEnter", p.Pid)
	p.clientID = cid
}

func (p *Player) OnLogin() {
	LogInfo("%d OnLogin", p.Pid)
}

func (p *Player) OnDisconnect() {
	LogInfo("%d OnDisconnect", p.Pid)
}
