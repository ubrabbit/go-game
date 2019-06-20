package player

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	db "server/db/mongodb"
	"server/db/table"
)

import (
	. "server/common"
)

func (p *Player) saveData() *table.DBTablePlayer {
	item := table.DBTablePlayer{
		ServerNum: p.ServerNum,
		Account:   p.Account,
		Pid:       p.Pid,
		Name:      p.Name,
		Grade:     p.Grade,
	}
	item.Data = string(JsonEncode(p.data))
	item.Container = string(JsonEncode(p.container))
	return &item
}

func (p *Player) loadData(item *table.DBTablePlayer) {
	p.Pid = item.Pid
	p.ServerNum = item.ServerNum
	p.Account = item.Account
	p.Name = item.Name
	p.Grade = item.Grade

	if len(item.Data) > 0 {
		JsonDecode(item.Data, &p.data)
	}
	if len(item.Container) > 0 {
		JsonDecode(item.Container, &p.container)
	}
}

func (p *Player) doSave() bool {
	item := p.saveData()
	c := db.GetServerC(table.DBNamePlayer)
	err := c.Update(bson.M{"pid": p.Pid}, item)
	if err != nil {
		LogPanic("%s doSave Error: %v", p.Repr(), err)
	}
	p.loaded = true
	p.update = false
	return true
}

func (p *Player) doLoad() {
	item := table.DBTablePlayer{}
	c := db.GetServerC(table.DBNamePlayer)
	err := c.Find(bson.M{"pid": p.Pid}).One(&item)
	if err == mgo.ErrNotFound || item.Pid != p.Pid {
		LogPanic("load player but pid %d not exists !", p.Pid)
	}
	p.loadData(&item)
	p.loaded = true
	p.update = false
}

func (p *Player) checkSave() bool {
	if !p.loaded || !p.update {
		return false
	}
	return p.doSave()
}

func (p *Player) checkLoad() bool {
	if !p.loaded {
		p.doLoad()
	}
	return p.loaded
}

func (p *Player) Query(attr string) interface{} {
	defer p.dbLock.Unlock()
	p.dbLock.Lock()

	v, _ := p.data[attr]
	return v
}

func (p *Player) Set(attr string, v interface{}) {
	defer p.dbLock.Unlock()
	p.dbLock.Lock()

	p.data[attr] = v
	p.Update()
}

/*---------------------------------------------------------------------------
// saveobject 的接口函数定义开始
*/
func (p *Player) UUID() string {
	return IntToString(p.ID())
}

func (p *Player) Repr() string {
	return fmt.Sprintf("Player[%s][%d]", p.Name, p.ID())
}

func (p *Player) Update() {
	p.update = true
}

func (p *Player) IsUpdate() bool {
	return p.update
}

func (p *Player) Save() {
	defer func() {
		defer p.dbLock.Unlock()
		r := recover()
		if r != nil {
			LogError("%s Save error: %v", p.Repr(), r)
		}
	}()
	p.dbLock.Lock()
	p.checkSave()
}

func (p *Player) Load() {
	defer func() {
		defer p.dbLock.Unlock()
		r := recover()
		if r != nil {
			LogError("%s Load error: %v", p.Repr(), r)
		}
	}()
	p.dbLock.Lock()
	p.checkLoad()
}

/*---------------------------------------------------------------------------
// saveobject 的接口函数定义结束
*/
