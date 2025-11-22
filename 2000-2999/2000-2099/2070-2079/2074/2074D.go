package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

type interval struct {
	l, r int64
}

func isqrt(x int64) int64 {
	r := int64(math.Sqrt(float64(x)))
	for (r+1)*(r+1) <= x {
		r++
	}
	for r*r > x {
		r--
	}
	return r
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}

	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)

		x := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &x[i])
		}
		r := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &r[i])
		}

		byY := make(map[int][]interval)

		for i := 0; i < n; i++ {
			ri := r[i]
			rSq := int64(ri) * int64(ri)
			xi := x[i]
			for dy := -ri; dy <= ri; dy++ {
				ddy := int64(dy)
				val := rSq - ddy*ddy
				dx := isqrt(val)
				y := dy
				byY[y] = append(byY[y], interval{l: xi - dx, r: xi + dx})
			}
		}

		var ans int64
		for _, segs := range byY {
			sort.Slice(segs, func(i, j int) bool {
				if segs[i].l == segs[j].l {
					return segs[i].r < segs[j].r
				}
				return segs[i].l < segs[j].l
			})
			curL, curR := segs[0].l, segs[0].r
			for _, s := range segs[1:] {
				if s.l > curR+1 {
					ans += curR - curL + 1
					curL, curR = s.l, s.r
				} else if s.r > curR {
					curR = s.r
				}
			}
			ans += curR - curL + 1
		}

		fmt.Fprintln(out, ans)
	}
}
