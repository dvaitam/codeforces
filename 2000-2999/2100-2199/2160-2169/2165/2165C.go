package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type FastScanner struct {
	r *bufio.Reader
}

func NewFastScanner() *FastScanner {
	return &FastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *FastScanner) NextInt() int64 {
	sign := int64(1)
	var val int64
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
	for c >= '0' && c <= '9' {
		val = val*10 + int64(c-'0')
		c, err = fs.r.ReadByte()
		if err != nil {
			break
		}
	}
	return sign * val
}

func main() {
	fs := NewFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := int(fs.NextInt())
	for ; t > 0; t-- {
		n := int(fs.NextInt())
		q := int(fs.NextInt())
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			a[i] = fs.NextInt()
		}
		sort.Slice(a, func(i, j int) bool { return a[i] > a[j] })
		k := n
		if k > 30 {
			k = 30
		}
		top := make([]int64, k)
		if k > 0 {
			copy(top, a[:k])
		}
		for ; q > 0; q-- {
			c := fs.NextInt()
			if c == 0 {
				fmt.Fprintln(out, 0)
				continue
			}
			if k == 0 {
				fmt.Fprintln(out, c)
				continue
			}
			slacks := make([]int64, k)
			copy(slacks, top)
			var cost int64
			for bit := 29; bit >= 0; bit-- {
				if (c>>bit)&1 == 0 {
					continue
				}
				v := int64(1) << bit
				idx := 0
				for i := 1; i < k; i++ {
					if slacks[i] > slacks[idx] {
						idx = i
					}
				}
				free := slacks[idx]
				var used int64
				if free > 0 {
					if free >= v {
						used = v
					} else {
						used = free
					}
				}
				cost += v - used
				slacks[idx] = free - v
			}
			fmt.Fprintln(out, cost)
		}
	}
}
