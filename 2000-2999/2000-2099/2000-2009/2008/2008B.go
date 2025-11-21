package main

import (
	"bufio"
	"fmt"
	"math"
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

		m := int(math.Round(math.Sqrt(float64(n))))
		if m*m != n {
			fmt.Fprintln(out, "No")
			continue
		}

		ok := true
		for i := 0; i < m && ok; i++ {
			for j := 0; j < m; j++ {
				ch := s[i*m+j]
				if i == 0 || i == m-1 || j == 0 || j == m-1 {
					if ch != '1' {
						ok = false
						break
					}
				} else {
					if ch != '0' {
						ok = false
						break
					}
				}
			}
		}

		if ok {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
	}
}
