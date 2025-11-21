package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
)

type FastScanner struct {
	r *bufio.Reader
}

func NewFastScanner(reader io.Reader) *FastScanner {
	return &FastScanner{r: bufio.NewReaderSize(reader, 1<<20)}
}

func (fs *FastScanner) NextInt64() int64 {
	sign := int64(1)
	val := int64(0)
	for {
		c, err := fs.r.ReadByte()
		if err != nil {
			return val * sign
		}
		if c == '-' {
			sign = -1
			c, err = fs.r.ReadByte()
			if err != nil {
				return 0
			}
		}
		if c >= '0' && c <= '9' {
			for {
				val = val*10 + int64(c-'0')
				c, err = fs.r.ReadByte()
				if err != nil {
					return val * sign
				}
				if c < '0' || c > '9' {
					break
				}
			}
			return val * sign
		}
	}
}

func can(dist int64, friends []int64, x int64, need int64) bool {
	if dist == 0 {
		return x+1 >= need
	}
	spread := dist - 1
	prevEnd := int64(-1)
	available := int64(0)
	for _, a := range friends {
		left := a - spread
		right := a + spread
		if left < 0 {
			left = 0
		}
		if right > x {
			right = x
		}
		if left > right {
			continue
		}
		if left > prevEnd+1 {
			available += left - (prevEnd + 1)
			if available >= need {
				return true
			}
		}
		if right > prevEnd {
			prevEnd = right
		}
	}
	if prevEnd < x {
		available += x - prevEnd
		if available >= need {
			return true
		}
	}
	return available >= need
}

func collect(dist int64, friends []int64, x int64, need int) []int64 {
	res := make([]int64, 0, need)
	if need == 0 {
		return res
	}
	if dist == 0 {
		for i := 0; i < need; i++ {
			res = append(res, int64(i))
		}
		return res
	}
	spread := dist - 1
	prevEnd := int64(-1)
	for _, a := range friends {
		left := a - spread
		right := a + spread
		if left < 0 {
			left = 0
		}
		if right > x {
			right = x
		}
		if left > right {
			continue
		}
		if left > prevEnd+1 {
			start := prevEnd + 1
			end := left - 1
			if start < 0 {
				start = 0
			}
			for pos := start; pos <= end && len(res) < need; pos++ {
				res = append(res, pos)
			}
			if len(res) == need {
				return res
			}
		}
		if right > prevEnd {
			prevEnd = right
		}
	}
	if prevEnd < x {
		start := prevEnd + 1
		if start < 0 {
			start = 0
		}
		for pos := start; pos <= x && len(res) < need; pos++ {
			res = append(res, pos)
		}
	}
	return res
}

func main() {
	fs := NewFastScanner(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := int(fs.NextInt64())
	for ; t > 0; t-- {
		n := int(fs.NextInt64())
		k := int(fs.NextInt64())
		x := fs.NextInt64()
		friends := make([]int64, n)
		for i := 0; i < n; i++ {
			friends[i] = fs.NextInt64()
		}
		sort.Slice(friends, func(i, j int) bool {
			return friends[i] < friends[j]
		})
		lo, hi := int64(0), x+1
		need := int64(k)
		for lo < hi {
			mid := (lo + hi + 1) / 2
			if can(mid, friends, x, need) {
				lo = mid
			} else {
				hi = mid - 1
			}
		}
		ans := collect(lo, friends, x, k)
		for i, pos := range ans {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, pos)
		}
		fmt.Fprintln(out)
	}
}
