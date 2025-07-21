package gws

import (
	"github.com/zishang520/socket.io/v2/socket"
	"sync"
)

var UserConnCtrl = NewUserConnMap()

type UserConnMap struct {
	mu sync.RWMutex
	m  map[string]*socket.Socket
}

func NewUserConnMap() *UserConnMap {
	return &UserConnMap{m: make(map[string]*socket.Socket)}
}

func (s *UserConnMap) Get(key string) (*socket.Socket, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.m[key]
	return val, ok
}

func (s *UserConnMap) Set(key string, value *socket.Socket) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[key] = value
}

func (s *UserConnMap) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, key)
}
