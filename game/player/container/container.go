package container

type ContainerInterface interface {
	ItemList() []interface{}
	AddItem(interface{}) bool
	RemoveItem(interface{}) interface{}
}

type DictContainer struct {
	Name string                 `json:name,string`
	Data map[string]interface{} `json:data,string`
}

type ListContainer struct {
	Name string        `json:name,string`
	Data []interface{} `json:data,[]string`
}
