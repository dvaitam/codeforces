package main

import (
	"bufio"
	"fmt"
	"os"
)

type fastScanner struct {
	r *bufio.Reader
}

func newScanner() *fastScanner {
	return &fastScanner{r: bufio.NewReaderSize(os.Stdin, 1<<20)}
}

func (fs *fastScanner) nextInt() int {
	sign, val := 1, 0
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

func (fs *fastScanner) nextString() string {
	var buf []byte
	c, _ := fs.r.ReadByte()
	for c == ' ' || c == '\n' || c == '\r' || c == '\t' {
		c, _ = fs.r.ReadByte()
	}
	for c != ' ' && c != '\n' && c != '\r' && c != '\t' {
		buf = append(buf, c)
		c, _ = fs.r.ReadByte()
	}
	return string(buf)
}

func main() {
	in := newScanner()
	out := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer out.Flush()

	T := in.nextInt()

	trans := [4][]int{
		0: {0, 1, 2, 3}, // both empty
		1: {2},          // top filled/blocked, bottom empty -> horizontal bottom
		2: {1},          // bottom filled/blocked, top empty -> horizontal top
		3: {0},          // both filled/blocked -> nothing to place
	}

	capAdd := func(a, b int) int {
		res := a + b
		if res > 2 {
			return 2
		}
		return res
	}

	for ; T > 0; T-- {
		n := in.nextInt()
		s1 := in.nextString()
		s2 := in.nextString()

		dp := [4]int{1, 0, 0, 0}

		for i := 0; i < n; i++ {
			block := 0
			if s1[i] == '#' {
				block |= 1
			}
			if s2[i] == '#' {
				block |= 2
			}
			var ndp [4]int
			for mask := 0; mask < 4; mask++ {
				if dp[mask] == 0 {
					continue
				}
				if mask&block != 0 {
					continue // domino would land on blocked cell
				}
				cur := mask | block
				for _, nm := range trans[cur] {
					ndp[nm] = capAdd(ndp[nm], dp[mask])
				}
			}
			dp = ndp
		}

		ans := dp[0]
		switch ans {
		case 0:
			fmt.Fprintln(out, "None")
		case 1:
			fmt.Fprintln(out, "Unique")
		default:
			fmt.Fprintln(out, "Multiple")
		}
	}
}
