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
	return &FastScanner{r: bufio.NewReaderSize(os.Stdin, 1<<20)}
}

func (fs *FastScanner) NextInt() int {
	sign, val := 1, 0
	c, err := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		if err != nil {
			return 0
		}
		c, err = fs.r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, err = fs.r.ReadByte()
	}
	for '0' <= c && c <= '9' {
		val = val*10 + int(c-'0')
		c, err = fs.r.ReadByte()
		if err != nil {
			break
		}
	}
	return sign * val
}

func main() {
	in := NewFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := in.NextInt()
	const limit = 60

	for ; t > 0; t-- {
		n := in.NextInt()
		a := make([]int, n)
		evens := make([]int, 0, 2)

		for i := 0; i < n; i++ {
			a[i] = in.NextInt()
			if a[i]%2 == 0 && len(evens) < 2 {
				evens = append(evens, a[i])
			}
		}

		if len(evens) == 2 {
			fmt.Fprintf(out, "%d %d\n", evens[0], evens[1])
			continue
		}

		check := n
		if check > limit {
			check = limit
		}

		found := false
		for i := 0; i < check && !found; i++ {
			for j := i + 1; j < n; j++ {
				if a[j]%a[i]%2 == 0 {
					fmt.Fprintf(out, "%d %d\n", a[i], a[j])
					found = true
					break
				}
			}
		}

		if !found {
			fmt.Fprintln(out, -1)
		}
	}
}
