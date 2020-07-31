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
	ColorCompleted = color.HEX("#fcba03")
	ColorTodo      = color.HEX("#777777")
)

type Simple struct {
	width     int
	charSet   []string
	doneColor color.RGBColor
	todoColor color.RGBColor
}

func NewSimple(width int) Simple {
	return Simple{
		width:     width,
		charSet:   lookup[HeavySquashTheme],
		doneColor: ColorCompleted,
		todoColor: ColorTodo,
	}
}

func NewSimpleWithTheme(width int, theme SimpleTheme, doneHexColor, todoHexColor color.RGBColor) Simple {
	return Simple{
		width:     width,
		charSet:   lookup[theme],
		doneColor: doneHexColor,
		todoColor: todoHexColor,
	}
}

func (s Simple) Format(p progress.Progress) (string, error) {
	completedRatio := p.Ratio()
	if completedRatio < 0 {
		completedRatio = 0
	}
	completedCount := int(completedRatio * float64(s.width))
	todoCount := s.width - completedCount
	if todoCount < 0 {
		todoCount = 0
	}

	completedSection := s.doneColor.Sprint(strings.Repeat(s.charSet[fullPosition], completedCount))
	todoSection := s.todoColor.Sprint(strings.Repeat(s.charSet[fullPosition], todoCount))

	return s.charSet[leftCapPosition] + completedSection + todoSection + s.charSet[rightCapPosition], nil
}
