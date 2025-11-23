package goober

import (
	"fmt"
	"io"
)

func Hexdump(in io.Reader, out io.Writer) error {
	const chunkSize = 4096
	buf := make([]byte, chunkSize)
	var offset int64
	fullReader := newFixedReader(in, chunkSize)
	for {
		n, err := fullReader.Read(buf)
		if n > 0 {
			for i := 0; i < n; i += 16 {
				size := 16
				if i+size > n {
					size = n - i
				}
				_, err := fmt.Fprintf(out, "%08x  ", offset)
				if err != nil {
					return err
				}
				line := generateOneLine(buf[i : i+size])
				_, err = out.Write(line)
				if err != nil {
					return err
				}
				offset += int64(size)
			}
		}

		if err == io.EOF {
			if offset%16 != 0 {
				_, err := fmt.Fprintf(out, "%08x\n", offset)

				if err != nil {
					return err
				}
			}
			return nil
		}
		if err != nil {
			return err
		}

	}
}

func hexEncode(i byte) byte {
	if i < 10 {
		return i + '0'
	}
	return (i - 10) + 'a'
}

func generateOneLine(raw []byte) []byte {
	const max = 16

	// Hex section – fixed width, positionally determined:
	// XX XX XX XX XX XX XX XX  XX XX XX XX XX XX XX XX
	hexbuf := make([]byte, 0, max*3+2) // +2 for mid-space
	for i := 0; i < max; i++ {
		if i < len(raw) {
			b := raw[i]
			hexbuf = append(hexbuf, hexEncode(b>>4), hexEncode(b&0x0F))
		} else {
			hexbuf = append(hexbuf, ' ', ' ')
		}
		if i == 7 {
			hexbuf = append(hexbuf, ' ', ' ') // double-space after 8 bytes
		} else {
			hexbuf = append(hexbuf, ' ')
		}
	}

	// ASCII section — flexible length, safe to append:
	asciibuf := make([]byte, 0, len(raw)+2) // +2 for pipes
	asciibuf = append(asciibuf, '|')
	for _, b := range raw {
		if b >= 32 && b <= 126 {
			asciibuf = append(asciibuf, b)
		} else {
			asciibuf = append(asciibuf, '.')
		}
	}
	asciibuf = append(asciibuf, '|')

	// Join them — NO mixing “layout code” with append:
	line := make([]byte, 0, len(hexbuf)+len(asciibuf)+2)
	line = append(line, hexbuf...)
	line = append(line, ' ')
	line = append(line, asciibuf...)
	line = append(line, '\n')
	return line
}
