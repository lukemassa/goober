package goober

import "io"

// fixedReader reads from an underlying io.Reader but
// guarantees that Read() returns fixed-size chunks (except at final EOF).
type fixedReader struct {
	in  io.Reader
	buf []byte
	off int
}

func newFixedReader(in io.Reader, readSize int) io.Reader {
	return &fixedReader{
		in:  in,
		buf: make([]byte, readSize),
	}
}

func (f *fixedReader) Read(p []byte) (int, error) {
	// Never copy more than caller expects.
	if len(p) < len(f.buf) {
		return 0, io.ErrShortBuffer
	}

	// Fill f.buf unless we hit EOF.
	for f.off < len(f.buf) {
		n, err := f.in.Read(f.buf[f.off:])
		f.off += n

		if err == io.EOF {
			if f.off == 0 { // no data at all -> true EOF
				return 0, io.EOF
			}
			// return partial chunk + io.EOF
			copy(p, f.buf[:f.off])
			n = f.off
			f.off = 0
			return n, io.EOF
		}
		if err != nil {
			return f.off, err
		}
	}

	// Full chunk ready
	copy(p, f.buf)
	n := f.off
	f.off = 0
	return n, nil
}
