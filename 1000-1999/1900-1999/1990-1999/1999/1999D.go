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

func main() {
	in := NewFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Sscan(in.Next(), &T)
	for ; T > 0; T-- {
		s := []byte(in.Next())
		t := in.Next()
		j := len(t) - 1
		for i := len(s) - 1; i >= 0 && j >= 0; i-- {
			if s[i] == t[j] {
				j--
			} else if s[i] == '?' {
				s[i] = t[j]
				j--
			}
		}
		if j >= 0 {
			fmt.Fprintln(out, "NO")
		} else {
			for i := 0; i < len(s); i++ {
				if s[i] == '?' {
					s[i] = 'a'
				}
			}
			fmt.Fprintln(out, "YES")
			fmt.Fprintln(out, string(s))
		}
	}
}
