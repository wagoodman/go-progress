package progress

import (
	"errors"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type readResult struct {
	n   int
	err error
}

type readerState struct {
	bytes int64
	err   error
}

type readStep struct {
	result   readResult
	expected readerState
}

type testReader struct {
	iteration int
	steps     []readStep
	size      int64
}

func (t *testReader) Read(p []byte) (n int, err error) {
	s := t.steps[t.iteration]
	t.iteration++
	return s.result.n, s.result.err
}

func TestReader(t *testing.T) {
	cases := []struct {
		name  string
		size  int64
		steps []readStep
	}{
		{
			name: "go case",
			size: -1,
			steps: []readStep{
				{
					result:   readResult{n: 0, err: io.EOF},
					expected: readerState{bytes: 0, err: io.EOF},
				},
			},
		},
		{
			name: "error passthrough",
			size: -1,
			steps: []readStep{
				{
					result:   readResult{n: 0, err: ErrCompleted},
					expected: readerState{bytes: 0, err: ErrCompleted},
				},
			},
		},
		{
			name: "multi step",
			size: -1,
			steps: []readStep{
				{
					result:   readResult{n: 5},
					expected: readerState{bytes: 5},
				},
				{
					result:   readResult{n: 5},
					expected: readerState{bytes: 10},
				},
				{
					result:   readResult{n: 3, err: io.EOF},
					expected: readerState{bytes: 13, err: io.EOF},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			reader := &testReader{
				steps: c.steps,
				size:  c.size,
			}

			monitor := NewReader(reader)

			for i, s := range c.steps {
				actualN, actualErr := monitor.Read(nil)

				if monitor.Size() != c.size {
					t.Fatalf("step %d: unexpected size: %d", i, monitor.Size())
				}

				if s.result.n != actualN {
					t.Fatalf("step %d: mismatched N: '%+v'!='%+v'", i, monitor.Current(), actualN)
				}

				if !errors.Is(actualErr, s.result.err) {
					t.Fatalf("step %d: mismatched err: '%+v'!='%+v'", i, actualErr, s.result.err)
				}

				if !errors.Is(monitor.Error(), s.expected.err) {
					t.Fatalf("step %d: mismatched cumulative err: '%+v'!='%+v'", i, monitor.Error(), s.expected.err)
				}

				if monitor.Current() != s.expected.bytes {
					t.Fatalf("step %d: mismatched byte count: '%+v'!='%+v'", i, monitor.Current(), s.expected.bytes)
				}

			}

		})
	}
}

func TestReaderPasthrough(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{
			name:  "gocase",
			input: "this is a test",
		},
		{
			// from https://github.com/minimaxir/big-list-of-naughty-strings/blob/master/blns.txt
			name: "horrible string",
			input: `
Ω≈ç√∫˜µ≤≥÷
åß∂ƒ©˙∆˚¬…æ
œ∑´®†¥¨ˆøπ“‘
¡™£¢∞§¶•ªº–≠
¸˛Ç◊ı˜Â¯˘¿
ÅÍÎÏ˝ÓÔÒÚÆ☃
Œ„´‰ˇÁ¨ˆØ∏”’
⁄€‹›ﬁﬂ‡°·‚—±
⅛⅜⅝⅞
ЁЂЃЄЅІЇЈЉЊЋЌЍЎЏАБВГДЕЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯабвгдежзийклмнопрстуфхцчшщъыьэюя
٠١٢٣٤٥٦٧٨٩

#	Unicode Subscript/Superscript/Accents
#
#	Strings which contain unicode subscripts/superscripts; can cause rendering issues

⁰⁴⁵
₀₁₂
⁰⁴⁵₀₁₂
ด้้้้้็็็็็้้้้้็็็็็้้้้้้้้็็็็็้้้้้็็็็็้้้้้้้้็็็็็้้้้้็็็็็้้้้้้้้็็็็็้้้้้็็็็ ด้้้้้็็็็็้้้้้็็็็็้้้้้้้้็็็็็้้้้้็็็็็้้้้้้้้็็็็็้้้้้็็็็็้้้้้้้้็็็็็้้้้้็็็็ ด้้้้้็็็็็้้้้้็็็็็้้้้้้้้็็็็็้้้้้็็็็็้้้้้้้้็็็็็้้้้้็็็็็้้้้้้้้็็็็็้้้้้็็็็`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			stringReader := strings.NewReader(c.input)
			monitor := NewReader(stringReader)

			actualBytes, err := ioutil.ReadAll(monitor)
			if err != nil {
				t.Fatalf("error when reading reader: %+v", err)
			}

			if string(actualBytes) != c.input {
				dmp := diffmatchpatch.New()
				diffs := dmp.DiffMain(string(actualBytes), c.input, true)
				t.Errorf("mismatched output:\n%s", dmp.DiffPrettyText(diffs))
			}

		})
	}
}
