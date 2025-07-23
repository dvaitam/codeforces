package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type pair struct {
	val int64
	id  int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}

	centersX := make([]pair, n)
	centersY := make([]pair, n)

	for i := 0; i < n; i++ {
		var x1, y1, x2, y2 int64
		fmt.Fscan(in, &x1, &y1, &x2, &y2)
		cx := x1 + x2
		cy := y1 + y2
		centersX[i] = pair{val: cx, id: i}
		centersY[i] = pair{val: cy, id: i}
	}

	sort.Slice(centersX, func(i, j int) bool { return centersX[i].val < centersX[j].val })
	sort.Slice(centersY, func(i, j int) bool { return centersY[i].val < centersY[j].val })

	// brute-force over counts removed from each side
	best := int64(1<<63 - 1)

	for a := 0; a <= k; a++ {
		for b := 0; b <= k; b++ {
			for c := 0; c <= k; c++ {
				for d := 0; d <= k; d++ {
					removed := make(map[int]struct{})
					for i := 0; i < a && i < n; i++ {
						removed[centersX[i].id] = struct{}{}
					}
					for i := 0; i < b && i < n; i++ {
						removed[centersX[n-1-i].id] = struct{}{}
					}
					for i := 0; i < c && i < n; i++ {
						removed[centersY[i].id] = struct{}{}
					}
					for i := 0; i < d && i < n; i++ {
						removed[centersY[n-1-i].id] = struct{}{}
					}
					if len(removed) > k {
						continue
					}
					li := 0
					for li < n {
						if _, ok := removed[centersX[li].id]; !ok {
							break
						}
						li++
					}
					ri := n - 1
					for ri >= 0 {
						if _, ok := removed[centersX[ri].id]; !ok {
							break
						}
						ri--
					}
					lj := 0
					for lj < n {
						if _, ok := removed[centersY[lj].id]; !ok {
							break
						}
						lj++
					}
					rj := n - 1
					for rj >= 0 {
						if _, ok := removed[centersY[rj].id]; !ok {
							break
						}
						rj--
					}
					if li > ri || lj > rj {
						continue
					}
					dx := centersX[ri].val - centersX[li].val
					dy := centersY[rj].val - centersY[lj].val
					w := dx / 2
					if dx%2 != 0 {
						w++
					}
					h := dy / 2
					if dy%2 != 0 {
						h++
					}
					if w < 1 {
						w = 1
					}
					if h < 1 {
						h = 1
					}
					area := w * h
					if area < best {
						best = area
					}
				}
			}
		}
	}

	if best == int64(1<<63-1) {
		best = 1
	}
	fmt.Fprintln(out, best)
}
