package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

type interval struct {
	l, r int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		segs := make([]interval, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &segs[i].l, &segs[i].r)
		}
		sort.Slice(segs, func(i, j int) bool {
			if segs[i].l == segs[j].l {
				return segs[i].r < segs[j].r
			}
			return segs[i].l < segs[j].l
		})

		constInf := int64(math.MaxInt64 / 4)
		dp := make([]int64, n+1)
		ndp := make([]int64, n+1)
		for i := range dp {
			dp[i] = constInf
			ndp[i] = constInf
		}
		dp[0] = 0
		maxS := 0
		pos := 1

		for _, seg := range segs {
			if seg.l > pos {
				lenO := seg.l - pos
				if lenO > 0 {
					newMax := maxS + lenO
					for i := 0; i <= n; i++ {
						ndp[i] = constInf
					}
					for s := 0; s <= maxS; s++ {
						if dp[s] == constInf {
							continue
						}
						for y := 0; y <= lenO; y++ {
							newS := s + lenO - y
							cost := dp[s] + int64(y)*int64(s)
							if cost < ndp[newS] {
								ndp[newS] = cost
							}
						}
					}
					dp, ndp = ndp, dp
					maxS = newMax
				}
			}
			length := seg.r - seg.l + 1
			newMax := maxS + length - 1
			for i := 0; i <= n; i++ {
				ndp[i] = constInf
			}
			for s := 0; s <= maxS; s++ {
				if dp[s] == constInf {
					continue
				}
				for x := 1; x <= length-1; x++ {
					newS := s + x
					penalty := int64(length-x) * int64(s+x)
					cost := dp[s] + penalty
					if cost < ndp[newS] {
						ndp[newS] = cost
					}
				}
			}
			dp, ndp = ndp, dp
			maxS = newMax
			pos = seg.r + 1
		}

		if pos <= n {
			lenO := n - pos + 1
			if lenO > 0 {
				newMax := maxS + lenO
				for i := 0; i <= n; i++ {
					ndp[i] = constInf
				}
				for s := 0; s <= maxS; s++ {
					if dp[s] == constInf {
						continue
					}
					for y := 0; y <= lenO; y++ {
						newS := s + lenO - y
						cost := dp[s] + int64(y)*int64(s)
						if cost < ndp[newS] {
							ndp[newS] = cost
						}
					}
				}
				dp, ndp = ndp, dp
				maxS = newMax
			}
		}

		best := constInf
		for s := 0; s <= maxS; s++ {
			if dp[s] < best {
				best = dp[s]
			}
		}
		totalPairs := int64(n) * int64(n-1) / 2
		fmt.Fprintln(out, totalPairs-best)
	}
}
