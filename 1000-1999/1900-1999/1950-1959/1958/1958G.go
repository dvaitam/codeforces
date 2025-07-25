package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	x := make([]int, k)
	for i := 0; i < k; i++ {
		fmt.Fscan(in, &x[i])
	}
	h := make([]int, k)
	for i := 0; i < k; i++ {
		fmt.Fscan(in, &h[i])
	}

	type tower struct{ x, h int }
	towers := make([]tower, k)
	for i := 0; i < k; i++ {
		towers[i] = tower{x[i], h[i]}
	}
	sort.Slice(towers, func(i, j int) bool { return towers[i].x < towers[j].x })
	for i := 0; i < k; i++ {
		x[i] = towers[i].x
		h[i] = towers[i].h
	}

	const INF = int(1e18)
	prefix := make([]int, n+2)
	maxR := -INF
	t := 0
	for i := 1; i <= n; i++ {
		for t < k && x[t] <= i {
			if v := x[t] + h[t]; v > maxR {
				maxR = v
			}
			t++
		}
		prefix[i] = maxR
	}

	suffix := make([]int, n+2)
	minL := INF
	t = k - 1
	for i := n; i >= 1; i-- {
		for t >= 0 && x[t] >= i {
			if v := x[t] - h[t]; v < minL {
				minL = v
			}
			t--
		}
		suffix[i] = minL
	}

	costPoint := func(p int) int {
		best := INF
		if prefix[p] != -INF {
			v := p - prefix[p]
			if v < 0 {
				v = 0
			}
			if v < best {
				best = v
			}
		}
		if suffix[p] != INF {
			v := suffix[p] - p
			if v < 0 {
				v = 0
			}
			if v < best {
				best = v
			}
		}
		return best
	}

	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var l, r int
		fmt.Fscan(in, &l, &r)
		if l > r {
			l, r = r, l
		}
		cpL := costPoint(l)
		cpR := costPoint(r)
		separate := cpL + cpR

		mid := (l + r) / 2
		candidates := []int{INF, INF, INF, INF}
		if prefix[l] != -INF {
			v := r - prefix[l]
			if v < 0 {
				v = 0
			}
			candidates[0] = v
		}
		if suffix[r] != INF {
			v := suffix[r] - l
			if v < 0 {
				v = 0
			}
			candidates[1] = v
		}
		if prefix[mid] != -INF {
			v := r - prefix[mid]
			if v < 0 {
				v = 0
			}
			candidates[2] = v
		}
		if mid+1 <= n && suffix[mid+1] != INF {
			v := suffix[mid+1] - l
			if v < 0 {
				v = 0
			}
			candidates[3] = v
		}
		bestOne := INF
		for _, v := range candidates {
			if v < bestOne {
				bestOne = v
			}
		}
		if separate < bestOne {
			bestOne = separate
		}
		if bestOne < separate {
			fmt.Fprintln(out, bestOne)
		} else {
			fmt.Fprintln(out, separate)
		}
	}
}
