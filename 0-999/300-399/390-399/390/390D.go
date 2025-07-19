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

	var n, m, k int
	if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
		return
	}

	type item struct {
		x, y int
		s    string
	}
	// Preallocate slice for BFS
	h := make([]item, 0, n*m)
	// Start at (1,1)
	h = append(h, item{1, 1, fmt.Sprintf("(1,1)")})
	done := 1
	// BFS-like generation: only right from first row and down
outer:
	for st := 0; st < len(h); st++ {
		if done >= k {
			break
		}
		x, y := h[st].x, h[st].y
		// move right if on first row
		if x == 1 && y < m {
			ns := fmt.Sprintf("%s (%d,%d)", h[st].s, 1, y+1)
			h = append(h, item{1, y + 1, ns})
			done++
			if done >= k {
				break outer
			}
		}
		// move down if possible
		if x < n {
			ns := fmt.Sprintf("%s (%d,%d)", h[st].s, x+1, y)
			h = append(h, item{x + 1, y, ns})
			done++
			if done >= k {
				break outer
			}
		}
	}
	// Keep only first k items
	if len(h) > k {
		h = h[:k]
	}
	// Sort by decreasing x+y
	sort.Slice(h, func(i, j int) bool {
		return h[i].x+h[i].y > h[j].x+h[j].y
	})
	// Compute answer: sum of path lengths = x+y-1
	ans := 0
	for _, it := range h {
		ans += it.x + it.y - 1
	}
	fmt.Fprintln(out, ans)
	for _, it := range h {
		fmt.Fprintln(out, it.s)
	}
}
