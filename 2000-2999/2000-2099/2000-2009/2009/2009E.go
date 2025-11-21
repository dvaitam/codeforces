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

func (fs *FastScanner) NextInt64() int64 {
	sign := int64(1)
	var val int64
	c, _ := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, _ = fs.r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, _ = fs.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int64(c-'0')
		c, _ = fs.r.ReadByte()
	}
	return sign * val
}

func value(i, k int64) int64 {
	return i * (2*k + i - 1)
}

func main() {
	in := NewFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := in.NextInt64()
	for ; t > 0; t-- {
		n := in.NextInt64()
		k := in.NextInt64()

		total := n * (2*k + n - 1) / 2

		l, r := int64(0), n
		for l < r {
			mid := (l + r + 1) >> 1
			if value(mid, k) <= total {
				l = mid
			} else {
				r = mid - 1
			}
		}

		ans := total - value(l, k)
		if ans < 0 {
			ans = -ans
		}
		if l < n {
			next := total - value(l+1, k)
			if next < 0 {
				next = -next
			}
			if next < ans {
				ans = next
			}
		}

		fmt.Fprintln(out, ans)
	}
}
