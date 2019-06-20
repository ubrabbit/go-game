package table

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)
import (
	. "server/common"
)

type DBTablePlayer struct {
	ServerNum int    `bson:"servernum"`
	Account   string `bson:"account"`
	Pid       int    `bson:"pid"`
	Name      string `bson:"name"`
	Grade     int    `bson:"grade"`
	Data      string `bson:"data"`
	Container string `bson:"container"`
}

func (tb *DBTablePlayer) TableName() string {
	return DBNamePlayer
}

func (tb *DBTablePlayer) Init(c *mgo.Collection) {
	err := c.EnsureIndex(mgo.Index{
		Key:    []string{"pid"},
		Unique: true,
		Sparse: true,
	})
	if err != nil {
		LogPanic("EnsureUniqueIndex %s error: %v", tb.TableName(), err)
	}
}

func (tb *DBTablePlayer) CheckKey() error {
	if tb.Pid <= 0 {
		return errors.New(FormatString("%s key is empty!", tb.TableName()))
	}
	return nil
}

func (tb *DBTablePlayer) Insert(c *mgo.Collection) error {
	err := tb.CheckKey()
	if err != nil {
		return err
	}
	err = c.Insert(tb)
	return err
}

func (tb *DBTablePlayer) Save(c *mgo.Collection) error {
	err := tb.CheckKey()
	if err != nil {
		return err
	}
	err = c.Update(
		bson.M{"pid": tb.Pid},
		tb,
	)
	return err
}

func (tb *DBTablePlayer) Load(c *mgo.Collection) error {
	err := tb.CheckKey()
	if err != nil {
		return err
	}
	err = c.Find(bson.M{"pid": tb.Pid}).One(tb)
	return err
}

func (tb *DBTablePlayer) Delete(c *mgo.Collection) error {
	err := tb.CheckKey()
	if err != nil {
		return err
	}
	err = c.Remove(bson.M{"pid": tb.Pid})
	return err
}
