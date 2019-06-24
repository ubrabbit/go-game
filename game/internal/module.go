package internal

import (
	"server/base"
	"server/game/world"
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

	world.InitWorld()
}

func (m *Module) OnDestroy() {

}
