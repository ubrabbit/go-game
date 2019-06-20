package client

import (
	"server/leaf/gate"
	"sync"
)

type LoginPlayer struct {
	ServerNum int
	Pid       int
	Name      string
	Grade     int
}

type LoginClient struct {
	Account     string
	Password    string
	agent       gate.Agent
	authSuccess bool
	isNew       bool
	playerList  map[int]LoginPlayer
	connectTime int
	errorCode   int
	errorMsg    string
}

const (
	constCleanInternal = 5 * 60
	constClientAlive   = 10 * 60
)

var g_Lock sync.Mutex
var g_LastCleanTime int
var g_LoginClientList map[string]*LoginClient
