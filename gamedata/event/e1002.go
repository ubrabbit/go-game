package event

type Event1002 struct {
}

func (e *Event1002) ID() int {
	return 1002
}

func (e *Event1002) Name() string {
	return "Event1002"
}

func (e *Event1002) Init(args ...interface{}) {
}

func (e *Event1002) Args() []interface{} {
	b := make([]interface{}, 0)
	b = append(b, 1)
	b = append(b, 2)
	b = append(b, 3)
	b = append(b, 4)
	b = append(b, 5)
	return b
}
