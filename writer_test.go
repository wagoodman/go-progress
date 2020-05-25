package progress

import (
	"errors"
	"strings"
	"testing"
)

type writeResult struct {
	n int
}

type writerState struct {
	bytes int64
}

type writeStep struct {
	result   writeResult
	expected writerState
}

func TestWriter(t *testing.T) {
	cases := []struct {
		name  string
		size  int64
		steps []writeStep
	}{
		{
			name: "go case",
			size: -1,
			steps: []writeStep{
				{
					result:   writeResult{n: 0},
					expected: writerState{bytes: 0},
				},
			},
		},
		{
			name: "multi step",
			size: -1,
			steps: []writeStep{
				{
					result:   writeResult{n: 5},
					expected: writerState{bytes: 5},
				},
				{
					result:   writeResult{n: 5},
					expected: writerState{bytes: 10},
				},
				{
					result:   writeResult{n: 3},
					expected: writerState{bytes: 13},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			writer := NewWriter()

			for i, s := range c.steps {

				p := []byte(strings.Repeat(".", s.result.n))

				actualN, actualErr := writer.Write(p)

				if writer.Size() != c.size {
					t.Fatalf("step %d: unexpected size: %d", i, writer.Size())
				}

				if s.result.n != actualN {
					t.Fatalf("step %d: mismatched N: '%+v'!='%+v'", i, writer.Current(), actualN)
				}

				if !errors.Is(actualErr, nil) {
					t.Fatalf("step %d: mismatched err: '%+v'!='%+v'", i, actualErr, nil)
				}

				if !errors.Is(writer.Error(), nil) {
					t.Fatalf("step %d: mismatched cumulative err: '%+v'!='%+v'", i, writer.Error(), nil)
				}

				if writer.Current() != s.expected.bytes {
					t.Fatalf("step %d: mismatched byte count: '%+v'!='%+v'", i, writer.Current(), s.expected.bytes)
				}

			}

		})
	}
}
