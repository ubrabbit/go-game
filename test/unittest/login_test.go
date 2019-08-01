package unittest

import (
	"math"
	"sync"
	"testing"
	"time"
)

import (
	"gopkg.in/mgo.v2/bson"
	. "server/common"
	db "server/db/mongodb"
	"server/db/table"
	. "server/msg/protocol"
	. "server/test/common"
)

const (
	playerCount  = 10
	echoCount    = 10
	echoInternal = 1 //ms
)

func TestLoginOldAccount(t *testing.T) {
	var loginUsers = map[string]string{
		"demo_1": "123456",
		"demo_2": "123456",
		"demo_3": "123456",
		"xyc":    "xyc",
		"lpx":    "lpx",
	}
	for user, password := range loginUsers {
		c := NewLoginClient()
		pid := c.Login(user, password)
		LogInfo("%s(%s) login success! pid=%d", user, password, pid)
	}
}

func TestLoginNewAccount(t *testing.T) {
	db.InitDB()
	c := db.GetServerC("global")

	playerID := table.DBTableGlobal{}
	err := c.Find(bson.M{"name": "MaxPlayerID"}).One(&playerID)
	if err != nil {
		LogPanic("Query MaxPlayerID Error: %v", err)
	}
	maxPid := playerID.Value.(int)
	wg := new(sync.WaitGroup)
	wg.Add(playerCount)
	for i := 0; i < playerCount; i++ {
		id := maxPid + i
		go func(id int) {
			user := FormatString("User_%d", id)
			password := FormatString("Password_%d", id)
			c := NewLoginClient()
			pid := c.Login(user, password)
			LogInfo("%s(%s) login success! pid=%d", user, password, pid)
			client := c.Client
			for j := 0; j < echoCount; j++ {
				v1, v2, v3, v4, v5, v6 := j, j+1, j*2, j*4, user, []byte{'a', 'b', 'c', 'b', 'a'}
				v1 = v1 % math.MaxUint8
				v2 = v2 % math.MaxUint16

				client.C2GSEcho(v1, v2, v3, v4, v5, v6)
				p0 := client.Protocol.(*TestEcho)
				client.GS2CEcho()
				p := client.Protocol.(*TestEcho)
				if (p.Int1 != v1) || (p.Int2 != v2) || (p.Int3 != v3) || (p.Int4 != v4) || (p.Str != v5) || (string(p.Byte) != string(v6)) {
					LogPanic("Player %d ClientEcho Fail! p: %v p0: %v", pid, p, p0)
				}
				time.Sleep(time.Duration(echoInternal) * time.Millisecond)
			}
			wg.Done()
		}(id)
	}
	wg.Wait()
}
