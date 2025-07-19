package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

type pair struct{ x, y float64 }
type interval struct{ l, r float64 }

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	pts := make([]pair, n)
	for i := 0; i < n; i++ {
		var xi, yi float64
		fmt.Fscan(in, &xi, &yi)
		pts[i] = pair{xi, yi}
		if xi == 0 && yi == 0 {
			fmt.Fprintln(out, "0")
			return
		}
	}
	sort.Slice(pts, func(i, j int) bool {
		if pts[i].x != pts[j].x {
			return pts[i].x < pts[j].x
		}
		return pts[i].y < pts[j].y
	})
	// remove duplicates
	uniq := pts[:1]
	for i := 1; i < len(pts); i++ {
		if pts[i].x != pts[i-1].x || pts[i].y != pts[i-1].y {
			uniq = append(uniq, pts[i])
		}
	}
	pts = uniq
	n = len(pts)
	dist := make([]float64, n)
	th := make([]float64, n)
	rb := math.Inf(1)
	for i, p := range pts {
		d := math.Hypot(p.x, p.y)
		dist[i] = d
		if d < rb {
			rb = d
		}
		th[i] = math.Atan2(p.y, p.x)
	}
	lb := 0.0
	b := make([]interval, n)
	var c []interval
	twoPI := 2 * math.Pi
	// check if radius r is feasible
	chk := func(r float64) bool {
		// build intervals
		for i := 0; i < n; i++ {
			ac := math.Acos(r / dist[i])
			l := th[i] - ac
			rr := th[i] + ac
			if l < 0 {
				l += twoPI
				rr += twoPI
			}
			b[i].l = l
			b[i].r = rr
		}
		sort.Slice(b, func(i, j int) bool { return b[i].l < b[j].l })
		// filter intervals
		c = c[:0]
		for i := 0; i < n; i++ {
			for len(c) > 0 && c[len(c)-1].r >= b[i].r {
				c = c[:len(c)-1]
			}
			if len(c) == 0 || c[0].r > b[i].r-twoPI {
				c = append(c, b[i])
			}
		}
		aa := len(c)
		if aa == 0 {
			return true
		}
		// duplicate intervals shifted by 2π
		orig := make([]interval, aa)
		copy(orig, c)
		for i := 0; i < aa; i++ {
			c = append(c, interval{orig[i].l + twoPI, orig[i].r + twoPI})
		}
		// build sparse table
		st := make([][17]int, aa)
		j := 0
		for i := 0; i < aa; i++ {
			for j < i+aa && c[j].l <= c[i].r {
				j++
			}
			st[i][0] = j - i
		}
		for k := 1; k < 17; k++ {
			for i := 0; i < aa; i++ {
				nxt := (i + st[i][k-1]) % aa
				st[i][k] = st[i][k-1] + st[nxt][k-1]
			}
		}
		// try covering with at most m intervals
		for s := 0; s < aa; s++ {
			used, covered, idx := 0, 0, s
			for k := 16; k >= 0; k-- {
				if used+(1<<k) <= m {
					covered += st[idx][k]
					idx = (idx + st[idx][k]) % aa
					used += 1 << k
				}
			}
			if covered >= aa {
				return true
			}
		}
		return false
	}
	// binary search for maximum r
	for rb-lb > 1e-7 {
		mb := (lb + rb) / 2
		if chk(mb) {
			lb = mb
		} else {
			rb = mb
		}
	}
	ans := (lb + rb) / 2
	fmt.Fprintf(out, "%.9f", ans)
}
