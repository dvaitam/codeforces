package main

import (
	"bufio"
	"fmt"
	"os"
)

type FastScanner struct {
	r *bufio.Reader
}

func NewFastScanner() *FastScanner {
	return &FastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *FastScanner) Next() string {
	c, _ := fs.r.ReadByte()
	for c == ' ' || c == '\n' || c == '\r' || c == '\t' {
		c, _ = fs.r.ReadByte()
	}
	buf := []byte{c}
	for {
		c, err := fs.r.ReadByte()
		if err != nil || c == ' ' || c == '\n' || c == '\r' || c == '\t' {
			break
		}
		buf = append(buf, c)
	}
	return string(buf)
}

func (fs *FastScanner) NextInt() int {
	sign := 1
	val := 0
	c, _ := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, _ = fs.r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, _ = fs.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, _ = fs.r.ReadByte()
	}
	return sign * val
}

func main() {
	in := NewFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := in.NextInt()
	for ; t > 0; t-- {
		n := in.NextInt()
		s := in.Next()
		r := in.Next()

		zeros := 0
		for i := 0; i < n; i++ {
			if s[i] == '0' {
				zeros++
			}
		}
		ones := n - zeros

		totalZeroR := 0
		for i := 0; i < n-1; i++ {
			if r[i] == '0' {
				totalZeroR++
			}
		}
		totalOneR := (n - 1) - totalZeroR

		if totalZeroR > ones || totalOneR > zeros {
			fmt.Fprintln(out, "NO")
			continue
		}

		prefixLen := n - 2
		prefixZero := 0
		prefixOne := 0
		for i := 0; i < prefixLen; i++ {
			if r[i] == '0' {
				prefixZero++
			} else {
				prefixOne++
			}
		}

		if prefixZero >= ones || prefixOne >= zeros {
			fmt.Fprintln(out, "NO")
		} else {
			fmt.Fprintln(out, "YES")
		}
	}
}
