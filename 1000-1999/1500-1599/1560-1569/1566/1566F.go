package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type seg struct{ l, r int }

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		pts := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &pts[i])
		}
		sort.Ints(pts)
		segs := make([]seg, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &segs[i].l, &segs[i].r)
		}
		sort.Slice(segs, func(i, j int) bool { return segs[i].l < segs[j].l })
		// simple greedy: assign each segment to nearest point
		cost := 0
		for _, s := range segs {
			// check if any point already inside
			idx := sort.SearchInts(pts, s.l)
			if idx < len(pts) && pts[idx] <= s.r {
				continue
			}
			// consider point to the left
			bestIdx := -1
			bestDist := int(1 << 60)
			if idx-1 >= 0 {
				d := 0
				p := pts[idx-1]
				if p < s.l {
					d = s.l - p
				} else if p > s.r {
					d = p - s.r
				} else {
					d = 0
				}
				if d < bestDist {
					bestDist = d
					bestIdx = idx - 1
				}
			}
			if idx < len(pts) {
				d := 0
				p := pts[idx]
				if p < s.l {
					d = s.l - p
				} else if p > s.r {
					d = p - s.r
				} else {
					d = 0
				}
				if d < bestDist {
					bestDist = d
					bestIdx = idx
				}
			}
			if bestIdx == -1 {
				continue
			}
			cost += bestDist
			// move point to nearest boundary inside segment
			if pts[bestIdx] < s.l {
				pts[bestIdx] = s.l
			} else if pts[bestIdx] > s.r {
				pts[bestIdx] = s.r
			}
		}
		fmt.Fprintln(out, cost)
	}
}
