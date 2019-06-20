package saveobject

import (
	. "server/common"
	"server/leaf/timer"
	"time"
)

func AddObject(obj SaveObject) {
	defer func() {
		g_SaveObjectList.Unlock()
		r := recover()
		if r != nil {
			LogPanic("AddObject %s(%s) error: %v", obj.UUID(), obj.Repr(), r)
		}
	}()
	g_SaveObjectList.Lock()
	uuid := obj.UUID()
	_, exists := g_SaveObjectList.saveList[uuid]
	if exists {
		LogPanic("uuid %s has been in saveList!", uuid)
	}
	g_SaveObjectList.saveList[uuid] = obj
}

func doRemove(uuid string) {
	_, exists := g_SaveObjectList.saveList[uuid]
	if exists {
		delete(g_SaveObjectList.saveList, uuid)
	}
}

func RemoveUUID(uuid string) {
	defer func() {
		g_SaveObjectList.Unlock()
		r := recover()
		if r != nil {
			LogPanic("RemoveUUID %s error: %v", uuid, r)
		}
	}()
	g_SaveObjectList.Lock()
	doRemove(uuid)
}

func RemoveObj(obj SaveObject) {
	defer func() {
		g_SaveObjectList.Unlock()
		r := recover()
		if r != nil {
			LogPanic("RemoveObj %s(%s) error: %v", obj.UUID(), obj.Repr(), r)
		}
	}()
	g_SaveObjectList.Lock()
	if obj == nil {
		return
	}
	doRemove(obj.UUID())
}

func Length() int {
	return len(g_SaveObjectList.saveList)
}

func checkSaveOne(obj SaveObject) (saved bool, err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = r.(error)
		}
	}()
	if obj.IsUpdate() {
		LogDebug("save %s", obj.Repr())
		obj.Save()
		return true, nil
	}
	return false, nil
}

func SaveAll() {
	defer func() {
		g_SaveObjectList.Unlock()
		r := recover()
		if r != nil {
			LogPanic("SaveAll error: %v", r)
		}
	}()
	g_SaveObjectList.Lock()
	LogInfo("SaveAll")
	for _, obj := range g_SaveObjectList.saveList {
		_, err := checkSaveOne(obj)
		if err != nil {
			LogError("save %s error: %v", obj.Repr(), err)
			continue
		}
	}
	LogInfo("SaveAll finished")
}

func (s *timerSave) checkSaveAll() {
	defer func() {
		g_SaveObjectList.Unlock()
		r := recover()
		if r != nil {
			LogError("checkSaveAll error: %v", r)
		}
	}()
	g_SaveObjectList.Lock()

	s.idx++
	if s.idx%dbSaveCycle == 0 {
		s.idx = 0
	}
	total := len(g_SaveObjectList.saveList)
	length := total / dbSaveCycle
	for _, obj := range g_SaveObjectList.saveList {
		saved, err := checkSaveOne(obj)
		if err != nil {
			LogError("save %s error: %v", obj.Repr(), err)
			continue
		}
		if saved {
			length--
		}
		//对象过多时，分批次存储，只在一个循环周期结束时全部存一次
		if s.idx != 0 && total >= 300 && length <= 0 {
			break
		}
	}
	LogInfo("checkSaveAll")
}

func InitSaveObject() {
	LogInfo("InitSaveObject")

	go func() {
		LogInfo("start db timer")
		s := timerSave{idx: 0}
		for {
			g_DBTimer.AfterFunc(time.Duration(dbSaveInternal)*time.Second, s.checkSaveAll)
			(<-g_DBTimer.ChanTimer).Cb()
		}
	}()
}

func init() {
	g_SaveObjectList = new(SaveObjectList)
	g_SaveObjectList.saveList = make(map[string]SaveObject, 0)

	g_DBTimer = timer.NewDispatcher(maxTimerLength)
}
