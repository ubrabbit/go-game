package table

import (
	"gopkg.in/mgo.v2"
)

type DBTableServer struct {
}

func (tb *DBTableServer) TableName() string {
	return DBNameServer
}

func (tb *DBTableServer) Init(c *mgo.Collection) {
}

func (tb *DBTableServer) CheckKey() error {
	return nil
}

func (tb *DBTableServer) Insert(c *mgo.Collection) error {
	err := c.Insert(tb)
	return err
}

func (tb *DBTableServer) Save(c *mgo.Collection) error {
	return nil
}
func (tb *DBTableServer) Load(c *mgo.Collection) error {
	return nil
}

func (tb *DBTableServer) Delete(c *mgo.Collection) error {
	return nil
}
