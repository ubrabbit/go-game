package internal

import (
	"server/base"
	"server/db/mongodb"
	"server/db/saveobject"
	"server/leaf/module"
)

var (
	skeleton = base.NewSkeleton()
	ChanRPC  = skeleton.ChanRPCServer
)

type Module struct {
	*module.Skeleton
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton

	saveobject.InitSaveObject()
	mongodb.InitDB()
}

func (m *Module) OnDestroy() {
	saveobject.SaveAll()
}
