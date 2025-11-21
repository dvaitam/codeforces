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
	return &fastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *fastScanner) NextInt64() int64 {
	sign := int64(1)
	var val int64
	c, err := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, err = fs.r.ReadByte()
		if err != nil {
			return 0
		}
	}
	if c == '-' {
		sign = -1
		c, err = fs.r.ReadByte()
		if err != nil {
			return 0
		}
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int64(c-'0')
		c, err = fs.r.ReadByte()
		if err != nil {
			break
		}
	}
	return val * sign
}

func (fs *fastScanner) NextInt() int {
	return int(fs.NextInt64())
}

var (
	n       int
	pref    []int64
	prefix2 []int64
	cumLen  []int64
	totPref []int64
)

func sumBlock(l int, t int64) int64 {
	if t <= 0 {
		return 0
	}
	r := l + int(t) - 1
	return (prefix2[r] - prefix2[l-1]) - t*pref[l-1]
}

func findBlock(idx int64) int {
	lo, hi := 1, n
	for lo < hi {
		mid := (lo + hi) >> 1
		if cumLen[mid] >= idx {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return lo
}

func prefixSum(idx int64) int64 {
	if idx <= 0 {
		return 0
	}
	totalElements := cumLen[n]
	if idx >= totalElements {
		return totPref[n]
	}
	block := findBlock(idx)
	prevCnt := cumLen[block-1]
	t := idx - prevCnt
	return totPref[block-1] + sumBlock(block, t)
}

func main() {
	fs := newFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	n = fs.NextInt()
	pref = make([]int64, n+1)
	prefix2 = make([]int64, n+1)

	for i := 1; i <= n; i++ {
		val := fs.NextInt64()
		pref[i] = pref[i-1] + val
		prefix2[i] = prefix2[i-1] + pref[i]
	}

	cumLen = make([]int64, n+1)
	totPref = make([]int64, n+1)

	for l := 1; l <= n; l++ {
		lenSeg := int64(n - l + 1)
		cumLen[l] = cumLen[l-1] + lenSeg
		blockSum := (prefix2[n] - prefix2[l-1]) - lenSeg*pref[l-1]
		totPref[l] = totPref[l-1] + blockSum
	}

	q := fs.NextInt()
	for ; q > 0; q-- {
		l := fs.NextInt64()
		r := fs.NextInt64()
		ans := prefixSum(r) - prefixSum(l-1)
		fmt.Fprintln(out, ans)
	}
}
