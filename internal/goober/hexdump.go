package goober

import (
	"fmt"
	"io"
	"strings"
)

func Hexdump(in io.Reader, out io.Writer) error {
	const chunkSize = 4096
	buf := make([]byte, chunkSize)
	var offset int64

	for {
		n, err := in.Read(buf)
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
			_, err := fmt.Fprintf(out, "%08x\n", offset)
			if err != nil {
				return err
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

func generateOneLine(rawLine []byte) []byte {
	var line strings.Builder
	// TODO: If we're at the end figure out the offset thing
	var summary strings.Builder
	for i := range len(rawLine) {
		if rawLine[i] >= 32 && rawLine[i] <= 127 {
			summary.WriteByte(rawLine[i])
		} else {
			summary.WriteRune('.')
		}
		line.WriteByte(hexEncode(rawLine[i] / 16))
		line.WriteByte(hexEncode(rawLine[i] % 16))
		line.WriteByte(' ')
		if i%8 == 7 {
			line.WriteRune(' ')
		}
	}
	line.WriteString(strings.Repeat(" ", (16-len(rawLine))*3))
	if len(rawLine) < 16 {
		line.WriteRune(' ')
	}
	if len(rawLine) < 8 {
		line.WriteRune(' ')
	}
	line.WriteRune('|')
	line.WriteString(summary.String())
	line.WriteRune('|')
	line.WriteRune('\n')
	// TODO: Don't make string then bytes
	return []byte(line.String())
}
