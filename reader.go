package progress

import (
	"io"

	"github.com/hashicorp/go-multierror"
)

// Reader should wrap another reader (acts as a bytes pass through)
type Reader struct {
	reader io.Reader
	bytes  int
	err    error
	size   int64
}

func NewSizedReader(reader io.Reader, size int64) *Reader {
	return &Reader{
		reader: reader,
		size:   size,
	}
}

func NewReader(reader io.Reader) *Reader {
	return &Reader{
		reader: reader,
		size:   -1,
	}
}

func (r *Reader) SetReader(reader io.Reader) {
	r.reader = reader
}

func (r *Reader) SetCompleted() {
	r.err = multierror.Append(r.err, ErrCompleted)
}

func (r *Reader) Read(p []byte) (n int, err error) {
	bytes, err := r.reader.Read(p)
	r.bytes += bytes
	if err != nil {
		r.err = multierror.Append(r.err, err)
	}
	return bytes, err
}

func (r *Reader) Current() int64 {
	return int64(r.bytes)
}

func (r *Reader) Size() int64 {
	return int64(r.size)
}

func (r *Reader) Error() error {
	return error(r.err)
}
