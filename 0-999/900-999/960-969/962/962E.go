package main

import (
	"bufio"
	"fmt"
	"os"
)

func solve(xs []int64, cs []byte) int64 {
	n := len(xs)
	var pIdx []int
	for i, c := range cs {
		if c == 'P' {
			pIdx = append(pIdx, i)
		}
	}
	if len(pIdx) == 0 {
		var ans int64
		var lastB int64
		var haveB bool
		var lastR int64
		var haveR bool
		for i := 0; i < n; i++ {
			switch cs[i] {
			case 'B':
				if haveB {
					ans += xs[i] - lastB
				}
				lastB = xs[i]
				haveB = true
			case 'R':
				if haveR {
					ans += xs[i] - lastR
				}
				lastR = xs[i]
				haveR = true
			}
		}
		return ans
	}
	ans := int64(0)
	firstP := pIdx[0]
	lastP := pIdx[len(pIdx)-1]
	// left of firstP
	firstB, firstR := int64(0), int64(0)
	haveB := false
	haveR := false
	for i := 0; i < firstP; i++ {
		if cs[i] == 'B' && !haveB {
			firstB = xs[i]
			haveB = true
		}
		if cs[i] == 'R' && !haveR {
			firstR = xs[i]
			haveR = true
		}
	}
	if haveB {
		ans += xs[firstP] - firstB
	}
	if haveR {
		ans += xs[firstP] - firstR
	}
	// right of lastP
	lastBpos, lastRpos := int64(0), int64(0)
	haveB = false
	haveR = false
	for i := n - 1; i > lastP; i-- {
		if cs[i] == 'B' && !haveB {
			lastBpos = xs[i]
			haveB = true
		}
		if cs[i] == 'R' && !haveR {
			lastRpos = xs[i]
			haveR = true
		}
	}
	if haveB {
		ans += lastBpos - xs[lastP]
	}
	if haveR {
		ans += lastRpos - xs[lastP]
	}
	// segments
	for idx := 0; idx < len(pIdx)-1; idx++ {
		i := pIdx[idx]
		j := pIdx[idx+1]
		L := xs[i]
		R := xs[j]
		seg := R - L
		maxB := int64(0)
		prev := L
		var hasB bool
		for t := i + 1; t < j; t++ {
			if cs[t] == 'B' {
				hasB = true
				if xs[t]-prev > maxB {
					maxB = xs[t] - prev
				}
				prev = xs[t]
			}
		}
		if R-prev > maxB {
			maxB = R - prev
		}
		maxR := int64(0)
		prev = L
		var hasR bool
		for t := i + 1; t < j; t++ {
			if cs[t] == 'R' {
				hasR = true
				if xs[t]-prev > maxR {
					maxR = xs[t] - prev
				}
				prev = xs[t]
			}
		}
		if R-prev > maxR {
			maxR = R - prev
		}
		cost2 := 3*seg - maxB - maxR
		if hasB && hasR {
			cost1 := 2 * seg
			if cost1 < cost2 {
				ans += cost1
			} else {
				ans += cost2
			}
		} else {
			ans += cost2
		}
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	xs := make([]int64, n)
	cs := make([]byte, n)
	for i := 0; i < n; i++ {
		var c string
		fmt.Fscan(in, &xs[i], &c)
		cs[i] = c[0]
	}
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	fmt.Fprintln(out, solve(xs, cs))
}
