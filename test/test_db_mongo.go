package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	. "server/common"
	db "server/db/mongodb"
	table "server/db/table"
)

func makePlayerInfo(pid int) (string, string, int) {
	account := FormatString("Account_%d", pid)
	name := FormatString("Name_%d", pid)
	grade := 1
	return account, name, grade
}

func test_insert() int {
	pid := db.NewPlayerID()
	account, name, grade := makePlayerInfo(pid)
	item := table.DBTablePlayer{
		ServerNum: 301,
		Account:   account,
		Pid:       pid,
		Name:      name,
		Grade:     grade,
	}
	c := db.GetServerC(table.DBNamePlayer)
	err := c.Insert(&item)
	if err != nil {
		LogPanic("test_insert_1 fail: %v !", err)
		return 0
	}
	LogInfo("insert pid %d success", pid)
	//重复插入
	item.Account = "Error Insert!"
	err = c.Insert(&item)
	if err == nil {
		LogPanic("test_insert_1 insert duplicate item!")
		return 0
	}
	LogInfo("test_insert success")
	return pid
}

func test_update(pid int) {
	c := db.GetServerC(table.DBNamePlayer)
	item := table.DBTablePlayer{}
	c.Find(bson.M{"pid": pid}).One(&item)

	newName := "New Player Name"
	item.Name = newName
	c.Update(bson.M{"pid": pid}, &item)
	c.Find(bson.M{"pid": pid}).One(&item)
	if item.Name != newName {
		LogPanic("test_update fail! Name '%s' != '%s'", item.Name, newName)
	}
	LogInfo("test_update success")
}

func test_delete(pid int) {
	c := db.GetServerC(table.DBNamePlayer)
	item := table.DBTablePlayer{}
	c.Remove(bson.M{"pid": pid})
	err := c.Find(bson.M{"pid": pid}).One(&item)
	if err != mgo.ErrNotFound {
		LogPanic("test_delete fail! pid %d still exists", pid)
	}
	LogInfo("test_delete success")
}

func main() {
	LogInfo("start")
	db.InitDB()

	pid := test_insert()
	test_update(pid)
	test_delete(pid)

	LogInfo("")
	LogInfo("")
	LogInfo("test success!")
}
