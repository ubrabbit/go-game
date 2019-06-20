package saveobject

import (
	"server/leaf/timer"
	"sync"
)

const (
	maxTimerLength = 64
	dbSaveInternal = 30
	dbSaveCycle    = 10
)

type SaveObject interface {
	UUID() string
	Repr() string
	Update()
	IsUpdate() bool
	Save()
	Load()
}

type SaveObjectList struct {
	sync.Mutex
	saveList map[string]SaveObject
}

type timerSave struct {
	idx int
}

var g_DBTimer *timer.Dispatcher
var g_SaveObjectList *SaveObjectList
