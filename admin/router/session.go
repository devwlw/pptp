package router

import (
	"log"
	"sync"
	"time"
)

var SessionMapInstance SessionMap

type SessionMap map[string]*Session

func (s SessionMap) Add(se *Session) {
	locker := &sync.Mutex{}
	locker.Lock()
	s[se.Id] = se
	locker.Unlock()
}

func (s SessionMap) Get(id string) *Session {
	return s[id]
}

func (s SessionMap) Del(id string) {
	delete(s, id)
}

type Session struct {
	Id      string
	Value   string
	Expired time.Time
}

func (s *Session) IsValid(value string) bool {
	log.Printf("current value:%s, check value:%s", s.Value, value)
	return value == s.Value
}

func (s *Session) IsExpired() bool {
	return time.Now().Unix() > s.Expired.Unix()
}
