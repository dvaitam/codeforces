package main

import (
	"bufio"
	"fmt"
	"os"
)

func sumMin(gaps []int64, l, d int64) int64 {
	r := l + d
	var s int64
	for _, g := range gaps {
		s += (g+r-1)/r - 1
	}
	return s
}

func sumMax(gaps []int64, l int64) int64 {
	var s int64
	for _, g := range gaps {
		s += g/l - 1
	}
	return s
}

func possible(diff int64, gaps []int64, k int64) bool {
	maxL := gaps[0]
	for _, g := range gaps {
		if g < maxL {
			maxL = g
		}
	}
	// minimal l with sumMin <= k
	lo, hi := int64(1), maxL+1
	for lo < hi {
		mid := (lo + hi) / 2
		if sumMin(gaps, mid, diff) <= k {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	lLow := lo
	if lLow == maxL+1 {
		return false
	}
	// maximal l with sumMax >= k
	lo, hi = 0, maxL
	for lo < hi {
		mid := (lo + hi + 1) / 2
		if sumMax(gaps, mid) >= k {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	lHigh := lo
	if lHigh == 0 {
		return false
	}
	return lLow <= lHigh
}

func solve(n int64, pos []int64, k int64) int64 {
	m := len(pos)
	gaps := make([]int64, m+1)
	prev := int64(0)
	for i := 0; i < m; i++ {
		gaps[i] = pos[i] - prev
		prev = pos[i]
	}
	gaps[m] = n - prev

	var lo, hi int64
	for _, g := range gaps {
		if g > hi {
			hi = g
		}
	}
	for lo < hi {
		mid := (lo + hi) / 2
		if possible(mid, gaps, k) {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return lo
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n, m, k int64
		fmt.Fscan(reader, &n, &m, &k)
		pos := make([]int64, m)
		for i := int64(0); i < m; i++ {
			fmt.Fscan(reader, &pos[i])
		}
		ans := solve(n, pos, k)
		fmt.Fprintln(writer, ans)
	}
}
