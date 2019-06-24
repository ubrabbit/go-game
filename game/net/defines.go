package net

import (
	"server/leaf/gate"
	"sync"
)

type NetClient struct {
	sync.Mutex
	id       int
	agent    gate.Agent
	playerID int
}

type NetContainer struct {
	sync.Mutex
	clients   map[string]*NetClient
	clientsID map[int]*NetClient
	playersID map[int]*NetClient
}

const (
	connectHelloTimeout = 10
	connectLoginTimeout = 600
)

var g_Container *NetContainer
