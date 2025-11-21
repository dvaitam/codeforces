package main

import (
	"bufio"
	"fmt"
	"os"
)

type fastScanner struct {
	r *bufio.Reader
}

func newFastScanner() *fastScanner {
	return &fastScanner{r: bufio.NewReaderSize(os.Stdin, 1<<20)}
}

func (fs *fastScanner) nextInt() int {
	sign := 1
	val := 0
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
		val = val*10 + int(c-'0')
		c, err = fs.r.ReadByte()
		if err != nil {
			break
		}
	}
	return sign * val
}

// countAtMost returns the number of subarrays whose length lies in [minLen, maxLen]
// and that contain at most distinctLimit distinct values.
func countAtMost(a []int, distinctLimit, minLen, maxLen int) int64 {
	if distinctLimit < 0 {
		return 0
	}
	freq := make(map[int]int)
	left := 0
	distinct := 0
	var ans int64
	for right, val := range a {
		freq[val]++
		if freq[val] == 1 {
			distinct++
		}
		for distinct > distinctLimit {
			cur := a[left]
			freq[cur]--
			if freq[cur] == 0 {
				delete(freq, cur)
				distinct--
			}
			left++
		}
		low := right - maxLen + 1
		if low < left {
			low = left
		}
		high := right - minLen + 1
		if high >= low {
			ans += int64(high - low + 1)
		}
	}
	return ans
}

func main() {
	fs := newFastScanner()
	out := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer out.Flush()

	t := fs.nextInt()
	for ; t > 0; t-- {
		n := fs.nextInt()
		k := fs.nextInt()
		l := fs.nextInt()
		r := fs.nextInt()
		a := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = fs.nextInt()
		}
		res := countAtMost(a, k, l, r) - countAtMost(a, k-1, l, r)
		fmt.Fprintln(out, res)
	}
}
