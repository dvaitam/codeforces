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
		uniqA := make(map[int]struct{})
		for i := 0; i < n; i++ {
			val := in.NextInt()
			uniqA[val] = struct{}{}
		}
		uniqB := make(map[int]struct{})
		for i := 0; i < n; i++ {
			val := in.NextInt()
			uniqB[val] = struct{}{}
		}

		sums := make(map[int]struct{})
		ans := "NO"
		for x := range uniqA {
			for y := range uniqB {
				sums[x+y] = struct{}{}
				if len(sums) >= 3 {
					ans = "YES"
					break
				}
			}
			if ans == "YES" {
				break
			}
		}
		fmt.Fprintln(out, ans)
	}
}
