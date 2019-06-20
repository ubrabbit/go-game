package player

import (
	"server/game/player/container"
	"sync"
)

type PlayerContainer struct {
	sync.Mutex
	playerList map[int]*Player
}

type Player struct {
	dbLock    sync.Mutex
	ServerNum int    `bson:"servernum" json:"servernum",int`
	Account   string `bson:"account" json:"account",string`
	Pid       int    `bson:"pid" json:"pid",int`
	Name      string `bson:"name" json:"name",string`
	Grade     int    `bson:"grade" json:"grade",int`
	data      map[string]interface{}
	container map[string]container.ContainerInterface
	loaded    bool
	update    bool
	clientID  int
}

var g_Container *PlayerContainer
