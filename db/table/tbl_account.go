package table

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)
import (
	. "server/common"
)

type DBTableAccount struct {
	Account  string `bson:"account"`
	Password string `bson:"password"`
	Email    string `bson:"email"`
	Players  []int  `bson:"players"`
}

func (tb *DBTableAccount) TableName() string {
	return DBNameAccount
}

func (tb *DBTableAccount) Init(c *mgo.Collection) {
	keys := []string{"account", "email"}
	for _, key := range keys {
		err := c.EnsureIndex(mgo.Index{
			Key:    []string{key},
			Unique: true,
			Sparse: true,
		})
		if err != nil {
			LogPanic("EnsureIndex %s error: %v", tb.TableName(), err)
		}
	}
}

func (tb *DBTableAccount) CheckKey() error {
	if len(tb.Account) <= 0 {
		return errors.New(FormatString("%s key is empty!", tb.TableName()))
	}
	return nil
}

func (tb *DBTableAccount) Insert(c *mgo.Collection) error {
	err := tb.CheckKey()
	if err != nil {
		return err
	}
	err = c.Insert(tb)
	return err
}

func (tb *DBTableAccount) Save(c *mgo.Collection) error {
	err := tb.CheckKey()
	if err != nil {
		return err
	}
	err = c.Update(
		bson.M{"account": tb.Account},
		tb,
	)
	return err
}

func (tb *DBTableAccount) Load(c *mgo.Collection) error {
	err := tb.CheckKey()
	if err != nil {
		return err
	}
	err = c.Find(bson.M{"account": tb.Account}).One(tb)
	return err
}

func (tb *DBTableAccount) Delete(c *mgo.Collection) error {
	err := tb.CheckKey()
	if err != nil {
		return err
	}
	err = c.Remove(bson.M{"account": tb.Account})
	return err
}
