package format

import (
	"strings"

	"github.com/gookit/color"
	"github.com/wagoodman/go-progress"
)

const (
	LiteTheme SimpleTheme = iota
	LiteSquashTheme
	HeavyTheme
	HeavySquashTheme
	ReallyHeavySquashTheme
	HeavyNoBarTheme
)

const (
	emptyPosition SimplePosition = iota
	fullPosition
	leftCapPosition
	rightCapPosition
)

type SimpleTheme int
type SimplePosition int

var lookup = map[SimpleTheme][]string{
	LiteTheme:              {" ", "─", "├", "┤"},
	LiteSquashTheme:        {" ", "─", "▕", "▏"},
	HeavyTheme:             {"━", "━", "┝", "┥"},
	HeavySquashTheme:       {"━", "━", "▕", "▏"},
	ReallyHeavySquashTheme: {"━", "━", "▐", "▌"},
	HeavyNoBarTheme:        {"━", "━", " ", " "},
}

var (
	doneColor = color.HEX("#ff8700")
	todoColor = color.HEX("#c6c6c6")
)

type Simple struct {
	width   int
	theme   SimpleTheme
	charSet []string
}

func NewSimple(width int, themes ...SimpleTheme) Simple {
	var theme SimpleTheme
	switch len(themes) {
	case 1:
		theme = themes[0]
	default:
		theme = HeavySquashTheme
	}

	return Simple{
		width:   width,
		theme:   theme,
		charSet: lookup[theme],
	}
}

func (s Simple) Format(p progress.Progress) (string, error) {

	completedRatio := p.Ratio()
	completedCount := int(completedRatio * float64(s.width))
	todoCount := s.width - completedCount

	completedSection := doneColor.Sprint(strings.Repeat(string(s.charSet[fullPosition]), completedCount))
	todoSection := todoColor.Sprint(strings.Repeat(string(s.charSet[fullPosition]), todoCount))

	return s.charSet[leftCapPosition] + completedSection + todoSection + s.charSet[rightCapPosition], nil
}
