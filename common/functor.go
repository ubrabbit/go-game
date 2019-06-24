package common

type Functor struct {
	Name     string
	Args     []interface{}
	callfunc func(...interface{})
}

type FuncInterface interface {
	Call(arg ...interface{})
}

func (f *Functor) Call(args ...interface{}) {
	defer func() {
		err := recover()
		if err != nil {
			LogError("functor '%s(%v)' Call Error: '%s'", f.Name, f.Args, err)
			return
		}
	}()
	f.callfunc(append(f.Args, args...)...)
}

func NewFunctor(name string, f func(...interface{}), args ...interface{}) *Functor {
	return &Functor{
		Name:     name,
		Args:     args,
		callfunc: f}
}
