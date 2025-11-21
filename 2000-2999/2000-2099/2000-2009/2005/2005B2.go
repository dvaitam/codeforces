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

func lowerBound(arr []int64, target int64) int {
	l, r := 0, len(arr)
	for l < r {
		mid := (l + r) >> 1
		if arr[mid] < target {
			l = mid + 1
		} else {
			r = mid
		}
	}
	return l
}

func main() {
	in := NewFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := in.NextInt64()
	for ; t > 0; t-- {
		n := in.NextInt64()
		m := int(in.NextInt64())
		q := int(in.NextInt64())

		teachers := make([]int64, m)
		for i := 0; i < m; i++ {
			teachers[i] = in.NextInt64()
		}
		sort.Slice(teachers, func(i, j int) bool {
			return teachers[i] < teachers[j]
		})

		for i := 0; i < q; i++ {
			pos := in.NextInt64()
			idx := lowerBound(teachers, pos)
			if idx == 0 {
				ans := teachers[0] - 1
				fmt.Fprintln(out, ans)
			} else if idx == len(teachers) {
				ans := n - teachers[len(teachers)-1]
				fmt.Fprintln(out, ans)
			} else {
				gap := teachers[idx] - teachers[idx-1] - 1
				ans := (gap + 1) / 2
				fmt.Fprintln(out, ans)
			}
		}
	}
}
