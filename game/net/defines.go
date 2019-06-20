package net

import (
	"server/leaf/gate"
	"sync"
)

type NetContainer struct {
	sync.Mutex
	clientIdx int
	clients   map[string]int
	agents    map[int]gate.Agent
	players   map[int]int
}

var g_Container *NetContainer
