package xnet

import "sync"

type Router struct {
	l sync.RWMutex
	m map[int]interface{}
}

func NewRouter() *Router {
	r := &Router{
		m: make(map[int]interface{}),
	}
	return r
}
