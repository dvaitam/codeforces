package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const inf64 = int64(1 << 60)

var pattern = []byte("docker")

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var s string
		fmt.Fscan(in, &s)
		m := len(s)

		// Group mismatch costs by starting position residue modulo 6.
		costs := make([][]int64, 6)
		if m >= 6 {
			for i := 0; i+6 <= m; i++ {
				var c int64
				for k := 0; k < 6; k++ {
					if s[i+k] != pattern[k] {
						c++
					}
				}
				costs[i%6] = append(costs[i%6], c)
			}
		}

		maxC := m / 6
		bestCost := make([]int64, maxC+1)
		for i := range bestCost {
			bestCost[i] = inf64
		}
		bestCost[0] = 0

		for r := 0; r < 6; r++ {
			if len(costs[r]) == 0 {
				continue
			}
			sort.Slice(costs[r], func(i, j int) bool { return costs[r][i] < costs[r][j] })
			pre := make([]int64, len(costs[r])+1)
			for i, v := range costs[r] {
				pre[i+1] = pre[i] + v
			}
			limit := len(costs[r])
			if limit > maxC {
				limit = maxC
			}
			for c := 1; c <= limit; c++ {
				if pre[c] < bestCost[c] {
					bestCost[c] = pre[c]
				}
			}
		}

		var n int
		fmt.Fscan(in, &n)
		diff := make([]int, maxC+2)
		for i := 0; i < n; i++ {
			var l, r int64
			fmt.Fscan(in, &l, &r)
			if l > int64(maxC) || r < 0 {
				continue
			}
			L := int(l)
			if L < 0 {
				L = 0
			}
			R := int(r)
			if R > maxC {
				R = maxC
			}
			if L <= R {
				diff[L]++
				diff[R+1]--
			}
		}

		bestCover := -1
		minChanges := inf64
		cur := 0
		for c := 0; c <= maxC; c++ {
			cur += diff[c]
			if bestCost[c] == inf64 { // unreachable count
				continue
			}
			if cur > bestCover {
				bestCover = cur
				minChanges = bestCost[c]
			} else if cur == bestCover && bestCost[c] < minChanges {
				minChanges = bestCost[c]
			}
		}

		fmt.Fprintln(out, minChanges)
	}
}
