package main

import (
	"bufio"
	"fmt"
	"io"
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
		if err == io.EOF {
			return 0
		}
		c, err = fs.r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, err = fs.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, err = fs.r.ReadByte()
		if err != nil {
			break
		}
	}
	return sign * val
}

func maxKMin(prefMin []int, l int) int {
	lo, hi := 0, len(prefMin)
	for lo < hi {
		mid := (lo + hi) >> 1
		if prefMin[mid] >= l {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return lo - 1
}

func maxKMax(prefMax []int, r int) int {
	lo, hi := 0, len(prefMax)
	for lo < hi {
		mid := (lo + hi) >> 1
		if prefMax[mid] <= r {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return lo - 1
}

func main() {
	in := NewFastScanner()
	out := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer out.Flush()

	t := in.NextInt()
	for ; t > 0; t-- {
		n := in.NextInt()
		q := in.NextInt()
		pos := make([]int, n)
		for i := 1; i <= n; i++ {
			val := in.NextInt()
			pos[val] = i
		}
		prefMin := make([]int, n+1)
		prefMax := make([]int, n+1)
		prefMin[0] = n + 1
		for i := 1; i <= n; i++ {
			prefMin[i] = prefMin[i-1]
			if pos[i-1] < prefMin[i] {
				prefMin[i] = pos[i-1]
			}
			if pos[i-1] > prefMax[i-1] {
				prefMax[i] = pos[i-1]
			} else {
				prefMax[i] = prefMax[i-1]
			}
		}

		best := 0
		for i := 0; i < q; i++ {
			l := in.NextInt()
			r := in.NextInt()
			a := maxKMin(prefMin, l)
			b := maxKMax(prefMax, r)
			mex := a
			if b < mex {
				mex = b
			}
			if mex > best {
				best = mex
			}
		}
		fmt.Fprintln(out, best)
	}
}
