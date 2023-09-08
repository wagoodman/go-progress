package progress

import "sync/atomic"

type Stager interface {
	Stage() string
}

type Stage struct {
	currentStage atomic.Value
}

func (s *Stage) getCurrentStage() string {
	return s.currentStage.Load().(string)
}

func (s *Stage) setCurrentStage(newStage string) {
	s.currentStage.Store(newStage)
}

func (s *Stage) Stage() string {
	return s.getCurrentStage()
}

func (s *Stage) Set(newStage string) {
	s.setCurrentStage(newStage)
}

func NewStage(current string) *Stage {
	result := Stage{}
	result.currentStage.Store(current)
	return &result
}

type StagerSettable interface {
	Stager
	Set(newStage string)
}

type StagedMonitorable interface {
	Stager
	Monitorable
}

type StagedProgressable interface {
	Stager
	Progressable
}

type StagedProgressor interface {
	Stager
	Progressor
}
