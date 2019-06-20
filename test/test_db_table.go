package main

/*
测试 server/db/table 里面所有结构的增删改查
*/

import (
	. "server/common"
	db "server/db/mongodb"
	table "server/db/table"
)

func test_global() {
	c := db.GetServerC(table.DBNameGlobal)
	item := table.DBTableGlobal{
		Name:  "test_global",
		Value: "test_global_value",
	}
	err := item.Insert(c)
	if err != nil {
		LogPanic("test_global Insert error: %v", err)
	}
	//重复插入
	err = item.Insert(c)
	if err == nil {
		LogPanic("test_global Insert Duplicate Error")
	}

	//读取
	item2 := table.DBTableGlobal{
		Name: "test_global",
	}
	err = item2.Load(c)
	if err != nil {
		LogPanic("test_global Load error: %v", err)
	}
	if item.Value != item2.Value {
		LogPanic("test_global insert fail: %v != %v", item.Value, item2.Value)
	}

	//修改
	item.Value = "Modify"
	err = item.Save(c)
	if err != nil {
		LogPanic("test_global Save Error: %v", err)
	}
	item2.Load(c)
	if item.Value != item2.Value {
		LogPanic("test_global Save fail: %v != %v", item.Value, item2.Value)
	}

	//删除
	err = item.Delete(c)
	if err != nil {
		LogPanic("test_global Delete Error: %v", err)
	}
	err = item2.Load(c)
	if err != db.ErrNotFound {
		LogPanic("test_global Delete fail: %v", item2.Value)
	}

	LogInfo("test_global success")
	LogInfo("")
}

func test_account() {
	c := db.GetServerC(table.DBNameAccount)
	item := table.DBTableAccount{
		Account:  "test_account",
		Password: "123456",
		Email:    "test@demo.com",
		Players:  []int{20001, 20002, 20003},
	}
	err := item.Insert(c)
	if err != nil {
		LogPanic("test_account Insert error: %v", err)
	}
	//重复插入
	err = item.Insert(c)
	if err == nil {
		LogPanic("test_account Insert Duplicate Error")
	}

	//读取
	item2 := table.DBTableAccount{
		Account: "test_account",
	}
	err = item2.Load(c)
	if err != nil {
		LogPanic("test_account Load error: %v", err)
	}
	if item.Email != item2.Email || item.Password != item2.Password {
		LogPanic("test_account insert fail: %v != %v", item, item2)
	}

	//修改
	item.Email = "demo@test.com"
	err = item.Save(c)
	if err != nil {
		LogPanic("test_account Save Error: %v", err)
	}
	item2.Load(c)
	if item.Email != item2.Email {
		LogPanic("test_account Save fail: %v != %v", item, item2)
	}

	//删除
	err = item.Delete(c)
	if err != nil {
		LogPanic("test_account Delete Error: %v", err)
	}
	err = item2.Load(c)
	if err != db.ErrNotFound {
		LogPanic("test_account Delete fail: %v", item2)
	}

	LogInfo("test_account success")
	LogInfo("")
}

func test_player() {
	c := db.GetServerC(table.DBNamePlayer)
	pid := db.NewPlayerID()
	item := table.DBTablePlayer{
		Account:   "test_account",
		Pid:       pid,
		Name:      "ubrabbit",
		Grade:     1,
		ServerNum: 10086,
	}
	err := item.Insert(c)
	if err != nil {
		LogPanic("test_player Insert error: %v", err)
	}
	//重复插入
	err = item.Insert(c)
	if err == nil {
		LogPanic("test_player Insert Duplicate Error")
	}

	//读取
	item2 := table.DBTablePlayer{
		Pid: pid,
	}
	err = item2.Load(c)
	if err != nil {
		LogPanic("test_player Load error: %v", err)
	}
	if item.Account != item2.Account || item.Name != item2.Name || item.Pid != item2.Pid {
		LogPanic("test_player insert fail: %v != %v", item, item2)
	}

	//修改
	item.Name = "xyc"
	err = item.Save(c)
	if err != nil {
		LogPanic("test_player Save Error: %v", err)
	}
	item2.Load(c)
	if item.Name != item2.Name {
		LogPanic("test_player Save fail: %v != %v", item, item2)
	}

	//删除
	err = item.Delete(c)
	if err != nil {
		LogPanic("test_player Delete Error: %v", err)
	}
	err = item2.Load(c)
	if err != db.ErrNotFound {
		LogPanic("test_player Delete fail: %v", item2)
	}

	LogInfo("test_player success")
	LogInfo("")
}

func test_error() {
	c := db.GetServerC(table.DBNameGlobal)
	item := table.DBTableGlobal{}
	err := item.Insert(c)
	LogInfo("insert empty global: %v : %v", item, err)

	err = item.Load(c)
	LogInfo("load empty global: %v : %v", item, err)

	err = item.Save(c)
	LogInfo("save empty global: %v : %v", item, err)

	err = item.Delete(c)
	LogInfo("save empty global: %v : %v", item, err)

	err = item.Load(c)
	LogInfo("load global after delete: %v : %v", item, err)
}

func main() {
	LogInfo("start")
	db.InitDB()

	LogInfo("")
	LogInfo("")
	test_global()
	test_account()
	test_player()
	test_error()

	LogInfo("")
	LogInfo("")
	LogInfo("test success!")
}
