package player

func newPlayer(pid int) *Player {
	p := new(Player)
	p.Pid = pid
	p.loaded = false
	p.update = false
	p.clientID = 0
	return p
}

func (p *Player) GetClientID() int {
	return p.clientID
}

func (p *Player) GetContainer(name string) interface{} {
	c, _ := p.container[name]
	return c
}

func GetPlayer(pid int) *Player {
	defer g_Container.Unlock()
	g_Container.Lock()
	p, _ := g_Container.playerList[pid]
	return p
}

func init() {
	g_Container = new(PlayerContainer)
	g_Container.playerList = make(map[int]*Player, 0)
}
