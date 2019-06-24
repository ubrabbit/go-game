package mongodb

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"server/conf"
	"server/db/table"
	db "server/leaf/db/mongodb"
)

import (
	. "server/common"
)

func (s *MongoSession) GetSession() *db.Session {
	return s.Context.Ref()
}

func (s *MongoSession) GetServerDB() *mgo.Database {
	return s.Context.Ref().DB(DatabaseName)
}

func (s *MongoSession) GetServerC(c string) *mgo.Collection {
	return s.Context.Ref().DB(DatabaseName).C(c)
}

func GetMongodb() *MongoSession {
	return g_MongoSession
}

func GetServerDB() *mgo.Database {
	return g_MongoSession.GetServerDB()
}

func GetServerC(c string) *mgo.Collection {
	_, ok := table.TableList[c]
	if !ok {
		LogPanic("unkown db collection: %s", c)
	}
	return g_MongoSession.GetServerC(c)
}

func QueryOne(name string, cond map[string]interface{}, result interface{}) error {
	c := GetServerC(name)
	return c.Find(bson.M(cond)).One(result)
}

func QueryAll(name string, cond map[string]interface{}, result interface{}) error {
	c := GetServerC(name)
	return c.Find(bson.M(cond)).All(result)
}

func InitPlayerIDCounter() {
	c := GetServerC(table.DBNameGlobal)
	item := table.DBTableGlobal{}
	key := "MaxPlayerID"
	err := c.Find(bson.M{"name": key}).One(&item)
	if err == mgo.ErrNotFound {
		LogInfo("InitPlayerIDCounter")
		_, err = c.Upsert(
			bson.M{"name": key},
			bson.M{"$set": bson.M{
				"name":  key,
				"value": minPlayerID,
			}},
		)
		if err != nil {
			LogPanic("Insert PlayerID Error: %v", err)
		}
		err = c.Find(bson.M{"name": key}).One(&item)
	}
	if err != nil {
		LogFatal("InitPlayerIDCounter Failure: %v", err)
	}
	LogInfo("MaxPlayerID: %d", item.Value.(int))
}

func NewPlayerID() int {
	c := GetServerC(table.DBNameGlobal)
	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"value": 1}},
		Upsert:    true,
		ReturnNew: true,
	}
	item := table.DBTableGlobal{}
	if _, err := c.Find(bson.M{"name": "MaxPlayerID"}).Apply(change, &item); err != nil {
		LogPanic("NewPlayerID failed: %v", err)
	}
	pid := item.Value.(int)
	if pid > PlayerIDMax {
		LogPanic("PlayerID %d exceed max: %d", pid, PlayerIDMax)
	}
	return pid
}

func NewUUID() bson.ObjectId {
	return bson.NewObjectId()
}

func NewUUIDString() string {
	return bson.NewObjectId().String()
}

func InitNewServer() {
	c := GetServerC(table.DBNameGlobal)
	item := table.DBTableGlobal{}
	key := "ServerInit"
	err := QueryOne(table.DBNameGlobal, bson.M{"name": key}, &item)
	if err == mgo.ErrNotFound {
		LogInfo("InitNewServer")
		_, err = c.Upsert(
			bson.M{"name": key},
			bson.M{"$set": bson.M{
				"name":  key,
				"value": GetTimeString(),
			}},
		)
	}
	if err != nil {
		LogPanic("ServerInit error: %v", err)
	}
}

func InitTables() {
	for name, obj := range table.TableList {
		LogInfo("init table %s", name)
		c := GetServerC(name)
		obj.Init(c)
	}
}

func makeDatabaseName() string {
	serverNum := conf.Server.ServerNum
	return fmt.Sprintf("db_%d", serverNum)
}

func InitDB() {
	DatabaseName = makeDatabaseName()

	url := conf.Server.DatabaseAddr
	if len(StripString(url)) == 0 {
		url = "localhost:27017"
		LogInfo("database not config, use default: %s", url)
	}
	LogInfo("connect: %s database: %s sessionNum: %d", url, DatabaseName, sessionNum)
	context, err := db.Dial(url, sessionNum)
	if err != nil {
		LogFatal("connect db error: %v", err)
	}
	g_MongoSession.Context = context
	LogInfo("connect: %s success", url)

	InitTables()
	InitNewServer()
	InitPlayerIDCounter()
}

func CloseDB() {
	defer func() {
		r := recover()
		if r != nil {
			LogError("CloseDB error: %v", r)
		}
	}()
	g_MongoSession.Context.Close()
	g_MongoSession.Context = nil
}

func init() {
	g_MongoSession = new(MongoSession)
	g_MongoSession.Context = nil
}
