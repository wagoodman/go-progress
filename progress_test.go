package progress

import (
	"fmt"
	"io"
	"testing"
)

func TestProgress(t *testing.T) {
	cases := []struct {
		current  int64
		size     int64
		err      error
		percent  float64
		ratio    float64
		complete bool
	}{
		{current: 1, size: 10, percent: 10.0, complete: false, ratio: 0.1},
		{current: 10, size: 10, percent: 100.0, complete: true, ratio: 1.0},
		{current: 0, size: 10, percent: 0.0, complete: false, ratio: 0.0},
		{current: 0, size: -1, percent: 0.0, complete: false, ratio: 0.0},
		{current: 2, size: -10, percent: 0.0, complete: false, ratio: 0.0},
		{current: 12, size: 10, percent: 100.0, complete: true, ratio: 1.0},
		{current: 2, size: 10, percent: 20.0, complete: true, ratio: 0.2, err: io.EOF},
		{current: 2, size: 10, percent: 20.0, complete: true, ratio: 0.2, err: ErrCompleted},
		{current: 2, size: 10, percent: 20.0, complete: false, ratio: 0.2, err: fmt.Errorf("blerg, err!")},
	}

	for _, c := range cases {
		name := fmt.Sprintf("%+v", c)
		t.Run(name, func(t *testing.T) {
			p := Progress{
				current: c.current,
				size:    c.size,
				err:     c.err,
			}

			if c.size != p.Size() {
				t.Errorf("size: expected '%v', got '%v'", c.size, p.Size())
			}

			if c.current != p.Current() {
				t.Errorf("current: expected '%v', got '%v'", c.current, p.Current())
			}

			if c.percent != p.Percent() {
				t.Errorf("percent: expected '%v', got '%v'", c.percent, p.Percent())
			}

			if c.ratio != p.Ratio() {
				t.Errorf("ratio: expected '%v', got '%v'", c.ratio, p.Ratio())
			}

			if c.err != p.Error() {
				t.Errorf("error: expected '%v', got '%v'", c.err, p.Error())
			}

			if c.complete != p.Complete() {
				t.Errorf("complete: expected '%v', got '%v'", c.complete, p.Complete())
			}
		})
	}
}
