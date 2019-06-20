package table

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

import (
	. "server/common"
)

type DBTableGlobal struct {
	Name  string      `bson:"name"`
	Value interface{} `bson:"value"`
}

func (tb *DBTableGlobal) TableName() string {
	return DBNameGlobal
}

func (tb *DBTableGlobal) Init(c *mgo.Collection) {
	err := c.EnsureIndex(mgo.Index{
		Key:    []string{"name"},
		Unique: true,
		Sparse: true,
	})
	if err != nil {
		LogPanic("EnsureIndex global error: %v", err)
	}
}

func (tb *DBTableGlobal) CheckKey() error {
	if len(tb.Name) <= 0 {
		return errors.New(FormatString("%s key is empty!", tb.TableName()))
	}
	return nil
}

func (tb *DBTableGlobal) Insert(c *mgo.Collection) error {
	err := tb.CheckKey()
	if err != nil {
		return err
	}
	err = c.Insert(tb)
	return err
}

func (tb *DBTableGlobal) Save(c *mgo.Collection) error {
	err := tb.CheckKey()
	if err != nil {
		return err
	}
	err = c.Update(
		bson.M{"name": tb.Name},
		tb,
	)
	return err
}

func (tb *DBTableGlobal) Load(c *mgo.Collection) error {
	err := tb.CheckKey()
	if err != nil {
		return err
	}
	err = c.Find(bson.M{"name": tb.Name}).One(tb)
	return err
}

func (tb *DBTableGlobal) Delete(c *mgo.Collection) error {
	err := tb.CheckKey()
	if err != nil {
		return err
	}
	err = c.Remove(bson.M{"name": tb.Name})
	return err
}
