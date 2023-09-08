package progress

import "sync/atomic"

type AtomicStage struct {
	current atomic.Value
}

func (s *AtomicStage) getCurrentStage() string {
	return s.current.Load().(string)
}

func (s *AtomicStage) setCurrent(new string) {
	s.current.Store(new)
}

func (s *AtomicStage) Stage() string {
	return s.getCurrentStage()
}

func (s *AtomicStage) Set(new string) {
	s.setCurrent(new)
}

func NewAtomicStage(current string) *AtomicStage {
	result := AtomicStage{}
	result.current.Store(current)
	return &result
}
