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
	ret := []byte("                                                  ")
	summary := []byte("|")
	offset := 0
	for i := range len(rawLine) {
		if rawLine[i] >= 32 && rawLine[i] <= 127 {
			summary = append(summary, rawLine[i])
		} else {
			summary = append(summary, '.')
		}
		ret[offset] = hexEncode(rawLine[i] / 16)
		ret[offset+1] = hexEncode(rawLine[i] % 16)
		offset += 3
		if i%8 == 7 {
			offset += 1
		}
	}
	summary = append(summary, '|')
	ret = append(ret, summary...)
	ret = append(ret, '\n')
	return ret
}
