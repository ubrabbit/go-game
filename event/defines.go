package event

import (
	. "server/common"
	"server/leaf/chanrpc"
	"sync"
)

type Event interface {
	ID() int
	Name() string
	Init(args ...interface{})
	Args() []interface{}
}

type Listener struct {
	sync.Mutex
	module     string
	listenList map[int]map[int]*Functor
}

const (
	EVENT_MODULE_GAME  = "EVENT_GAME"
	EVENT_MODULE_DB    = "EVENT_DB"
	EVENT_MODULE_LOGIN = "EVENT_LOGIN"
	EVENT_MODULE_SCENE = "EVENT_SCENE"
)

var listenerModule map[string]*chanrpc.Server
