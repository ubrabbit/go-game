package table

import (
	"gopkg.in/mgo.v2"
)

type DBTable interface {
	Init(c *mgo.Collection)
	TableName() string
	CheckKey() error
	Insert(c *mgo.Collection) error
	Save(c *mgo.Collection) error
	Load(c *mgo.Collection) error
	Delete(c *mgo.Collection) error
}

//必须要注册了集合后才能存储，避免写代码时放飞自我
const (
	DBNameGlobal  = "global"
	DBNameServer  = "server"
	DBNamePlayer  = "player"
	DBNameAccount = "account"
)

var TableList = map[string]DBTable{
	DBNameGlobal:  &DBTableGlobal{},
	DBNameServer:  &DBTableServer{},
	DBNamePlayer:  &DBTablePlayer{},
	DBNameAccount: &DBTableAccount{},
}
