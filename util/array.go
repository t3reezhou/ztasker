package util

import "sync"

type Int64Set struct {
	m map[int64]struct{}
	sync.RWMutex
}

func NewInt64Set(origin []int64) *Int64Set {
	set := make(map[int64]struct{}, 0)
	for _, i := range origin {
		set[i] = struct{}{}
	}
	return &Int64Set{m: set}
}

func (s *Int64Set) Add(i int64) *Int64Set {
	s.Lock()
	defer s.Unlock()
	s.m[i] = struct{}{}
	return s
}

func (s *Int64Set) Remove(i int64) *Int64Set {
	s.Lock()
	defer s.Unlock()
	delete(s.m, i)
	return s
}

func (s *Int64Set) Has(i int64) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.m[i]
	return ok
}

func (s *Int64Set) Int64List() []int64 {
	s.RLock()
	defer s.RUnlock()
	list := make([]int64, 0, len(s.m))
	for e, _ := range s.m {
		list = append(list, e)
	}
	return list
}
