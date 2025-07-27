package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// segCost returns minimal flips to go from (sa,ca) to (sb,cb).
func segCost(sa, ca, sb, cb int64) int64 {
	dr := sb - sa
	dc := cb - ca
	if dr < 0 || dc < 0 || dc > dr {
		return 0
	}
	if (sa+ca)&1 == 0 {
		if dc == dr {
			return dr
		}
		return (dr - dc) / 2
	}
	return (dr - dc + 1) / 2
}

func solve() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		r := make([]int64, n)
		c := make([]int64, n)
		for i := range r {
			fmt.Fscan(in, &r[i])
		}
		for i := range c {
			fmt.Fscan(in, &c[i])
		}
		pts := make([][2]int64, n)
		for i := 0; i < n; i++ {
			pts[i] = [2]int64{r[i], c[i]}
		}
		sort.Slice(pts, func(i, j int) bool {
			if pts[i][0] == pts[j][0] {
				return pts[i][1] < pts[j][1]
			}
			return pts[i][0] < pts[j][0]
		})

		var ans int64
		curR, curC := int64(1), int64(1)
		for _, p := range pts {
			ans += segCost(curR, curC, p[0], p[1])
			curR, curC = p[0], p[1]
		}
		fmt.Fprintln(out, ans)
	}
}

func main() { solve() }
