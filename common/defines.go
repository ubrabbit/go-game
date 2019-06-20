package common

const (
	PlayerIDMin = 1000
	PlayerIDMax = 3*10000*10000 - 1

	ObjectIDMin = PlayerIDMax + 1
	ObjectIDMax = 0xFFFFFFFF - 10000
)

type Object interface {
	ID() int
	Create()
	Delete()
}

var g_ObjectIDChan chan int
