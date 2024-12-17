package main

import "fmt"

type SafeMap struct {
	value chan map[string]int
}

func NewExample() *SafeMap {
	initialMap := make(map[string]int)
	sm := &SafeMap{
		value: make(chan map[string]int, 1),
	}
	go func() {
		sm.value <- initialMap // 初始化通道
	}()
	return sm
	}

func (sm *SafeMap) Set(key string, value int) {
	newMap := make(map[string]int)
	for k, v := range <-sm.value {
		newMap[k] = v
	}
	newMap[key] = value
	sm.value <- newMap
}

func (sm *SafeMap) Get(key string) (int, bool) {
	m := <-sm.value
	val, ok := m[key]
	sm.value <- m
	return val, ok
}

func main() {
	example := NewExample()
	example.Set("xiongjie", 19)
	val, ok := example.Get("xiongjie")
	if ok {
		fmt.Println(val)
	}
}