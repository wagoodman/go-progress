package progress

import (
	"errors"
	"io"

	"github.com/hashicorp/go-multierror"
)

const (
	DefaultStrategy AggregationStrategy = iota
	NormalizeStrategy
)

type AggregationStrategy int

type AggregateGenerator struct {
	progs        []Progressable
	names        map[int]string
	nextIndex    int
	currentIndex int
	strategy     AggregationStrategy
}

func NewAggregateGenerator(p ...Progressable) *AggregateGenerator {
	strategy := NormalizeStrategy
	if p == nil {
		p = make([]Progressable, 0)
	}
	return &AggregateGenerator{
		progs:    p,
		names:    make(map[int]string),
		strategy: strategy,
	}
}

func (a *AggregateGenerator) SetStrategy(strategy AggregationStrategy) {
	a.strategy = strategy
}

func (a *AggregateGenerator) Add(p ...Progressable) {
	a.progs = append(a.progs, p...)
	a.nextIndex += len(p)
}

func (a *AggregateGenerator) AddNamed(name string, p Progressable) {
	a.names[a.nextIndex] = name
	a.nextIndex++
}

func (a *AggregateGenerator) CurrentName() string {
	if name, ok := a.names[a.currentIndex]; ok {
		return name
	}
	return ""
}

func (a *AggregateGenerator) Progress() Progress {
	result := Progress{}
	var completedProgs int
	var currentIndex = -1
	for idx, p := range a.progs {

		switch a.strategy {
		case NormalizeStrategy:
			if p.Size() < 0 {
				result.current = 0
			} else {
				result.current += int64(100 / (float64(p.Size()) / float64(p.Current())))
			}
			result.size += 100
		default:
			result.current += p.Current()
			result.size += p.Size()
		}

		// capture notable errors
		err := p.Error()
		if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, ErrCompleted) {
			result.err = multierror.Append(result.err, err)
		}
		if isCompleted(p) {
			completedProgs++
		} else {
			// the first non-completed task is the current task
			if currentIndex <= 0 {
				currentIndex = idx
			}
		}
	}
	if currentIndex <= 0 {
		a.currentIndex = 0
	} else {
		a.currentIndex = currentIndex
	}

	if completedProgs == len(a.progs) {
		result.err = multierror.Append(result.err, ErrCompleted)
	}
	return result
}
