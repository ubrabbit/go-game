package main

import (
	"fmt"
	. "server/common"
	db "server/db/mongodb"
	"server/db/saveobject"
	"sync"
	"time"
)

type SavePlayer struct {
	Pid    int
	uuid   string
	update bool
	saved  int
}

func (p *SavePlayer) UUID() string {
	return p.uuid
}

func (p *SavePlayer) Repr() string {
	return fmt.Sprintf("Player(%d)", p.Pid)
}

func (p *SavePlayer) Update() {
	p.update = true
}

func (p *SavePlayer) IsUpdate() bool {
	return p.update
}

func (p *SavePlayer) Save() {
	LogInfo("%s Save", p.Repr())
	p.saved++
}

func (p *SavePlayer) Load() {
	LogInfo("%s Load", p.Repr())
}

type SaveNpc struct {
	Nid    int
	uuid   string
	update bool
	saved  int
}

func (p *SaveNpc) UUID() string {
	return p.uuid
}

func (p *SaveNpc) Repr() string {
	return fmt.Sprintf("Npc(%d)", p.Nid)
}

func (p *SaveNpc) Update() {
	p.update = true
	LogInfo("%s Update", p.Repr())
}

func (p *SaveNpc) IsUpdate() bool {
	return p.update
}

func (p *SaveNpc) Save() {
	LogInfo("%s Save", p.Repr())
	p.saved++
}

func (p *SaveNpc) Load() {
	LogInfo("%s Load", p.Repr())
}

func test_1() {
	totalPlayer := 5
	totalNpc := 5
	playerPool := make(map[int]*SavePlayer, 0)
	npcPool := make(map[int]*SaveNpc, 0)

	deleteUUID1 := ""
	deleteUUID2 := ""
	var deletePlayer1 *SavePlayer = nil
	var deletePlayer2 *SavePlayer = nil
	var deleteNpc1 *SaveNpc = nil
	var deleteNpc2 *SaveNpc = nil

	for i := 1; i < totalPlayer+1; i++ {
		p := new(SavePlayer)
		p.Pid = i
		p.uuid = db.NewUUIDString()
		p.update = false
		p.saved = 0
		playerPool[i] = p
		saveobject.AddObject(p)
		if len(deleteUUID1) == 0 {
			deleteUUID1 = p.uuid
		} else if deletePlayer1 == nil {
			deletePlayer1 = p
		} else {
			deletePlayer2 = p
		}
	}
	for i := 1; i < totalNpc+1; i++ {
		p := new(SaveNpc)
		p.Nid = i
		p.uuid = db.NewUUIDString()
		p.update = false
		p.saved = 0
		npcPool[i] = p
		saveobject.AddObject(p)
		if len(deleteUUID2) == 0 {
			deleteUUID2 = p.uuid
		} else if deleteNpc1 == nil {
			deleteNpc1 = p
		} else {
			deleteNpc2 = p
		}
	}
	LogInfo("Saveobject Length== %d", saveobject.Length())
	LogInfo("Delete: %v %v %v %v %v %v", deleteUUID1, deleteUUID2, deletePlayer1, deletePlayer2, deleteNpc1, deleteNpc2)
	saveobject.RemoveUUID(deleteUUID1)
	saveobject.RemoveUUID(deleteUUID2)
	saveobject.RemoveObj(deletePlayer1)
	saveobject.RemoveObj(deletePlayer2)
	saveobject.RemoveObj(deleteNpc1)
	saveobject.RemoveObj(deleteNpc2)
	saveobject.RemoveUUID("")
	saveobject.RemoveObj(nil)
	LogInfo("Saveobject Length== %d", saveobject.Length())
	if saveobject.Length() != (totalPlayer + totalNpc - 6) {
		LogPanic("test_1 fail!")
	}
	LogInfo("test_1 success")
}

func test_2() {
	testCount := 10000
	totalPlayer := testCount
	totalNpc := testCount
	playerPool := make(map[int]*SavePlayer, 0)
	npcPool := make(map[int]*SaveNpc, 0)

	for i := 1; i < totalPlayer+1; i++ {
		p := new(SavePlayer)
		p.Pid = i
		p.uuid = db.NewUUIDString()
		p.update = false
		p.saved = 0
		playerPool[i] = p
		saveobject.AddObject(p)
	}
	for i := 1; i < totalNpc+1; i++ {
		p := new(SaveNpc)
		p.Nid = i
		p.uuid = db.NewUUIDString()
		p.update = false
		p.saved = 0
		npcPool[i] = p
		saveobject.AddObject(p)
	}

	go saveobject.InitSaveObject()

	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func() {
		for k := 0; k < 6; k++ {
			for _, p := range playerPool {
				p.Update()
			}
			time.Sleep(5 * time.Second)
		}
		wg.Done()
	}()
	go func() {
		for k := 0; k < 6; k++ {
			for _, p := range npcPool {
				p.Update()
			}
			time.Sleep(10 * time.Second)
		}
		wg.Done()
	}()
	wg.Wait()

	time.Sleep(10 * time.Second)
	saveobject.SaveAll()

	for _, p := range playerPool {
		if p.saved <= 0 {
			LogPanic("%s not saved!", p.Repr())
		}
	}
	for _, p := range npcPool {
		if p.saved <= 0 {
			LogPanic("%s not saved!", p.Repr())
		}
	}
	LogInfo("test_2 success")
}

func main() {
	LogInfo("start")
	db.InitDB()

	test_1()
	LogInfo("\n\n")
	test_2()

	LogInfo("test success!")
}
