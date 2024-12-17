package main

import (
	"fmt"
	"sync"
)




type Map struct {
	lock sync.RWMutex
	value map[string]int
}


func NewExamlpe()*Map  {
	return &Map{
		value:make(map[string]int),
	}
}

func (a*Map)Set(key string, value int)  {
	a.lock.Lock()
	a.value[key] = value
	a.lock.Unlock()
}

func(a*Map)Get(key string) (int, bool){
	a.lock.RLock()
	defer a.lock.RUnlock()
	val, ok := a.value[key]
	return val, ok
}

func main() {
	safemap:= NewExamlpe()
	safemap.Set("xiongjie",19)
	val, ok := safemap.Get("xiongjie")
	if ok {
		fmt.Println(val)
	}
}