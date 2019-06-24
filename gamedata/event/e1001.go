package event

type Event1001 struct {
}

func (e *Event1001) ID() int {
	return 1001
}

func (e *Event1001) Name() string {
	return "Event1001"
}

func (e *Event1001) Init(args ...interface{}) {
}

func (e *Event1001) Args() []interface{} {
	b := make([]interface{}, 0)
	b = append(b, "a")
	b = append(b, "b")
	b = append(b, "c")
	b = append(b, "d")
	b = append(b, "e")
	return b
}
