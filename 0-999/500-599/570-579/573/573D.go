package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type item struct {
	val int64
	idx int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}

	w := make([]int64, n)
	for i := range w {
		fmt.Fscan(in, &w[i])
	}
	h := make([]int64, n)
	for i := range h {
		fmt.Fscan(in, &h[i])
	}

	type pair struct {
		val int64
		idx int
	}

	// wOrder stores indices of warriors sorted by strength descending
	wOrder := make([]int, n)
	for i := range wOrder {
		wOrder[i] = i
	}
	sort.Slice(wOrder, func(i, j int) bool {
		if w[wOrder[i]] == w[wOrder[j]] {
			return wOrder[i] < wOrder[j]
		}
		return w[wOrder[i]] > w[wOrder[j]]
	})

	for ; q > 0; q-- {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a--
		b--
		h[a], h[b] = h[b], h[a]

		horses := make([]pair, n)
		for i := 0; i < n; i++ {
			horses[i] = pair{val: h[i], idx: i}
		}
		sort.Slice(horses, func(i, j int) bool {
			if horses[i].val == horses[j].val {
				return horses[i].idx < horses[j].idx
			}
			return horses[i].val > horses[j].val
		})

		// fix conflicts by local swaps
		for i := 0; i < n; i++ {
			if wOrder[i] == horses[i].idx {
				if i+1 < n {
					horses[i], horses[i+1] = horses[i+1], horses[i]
				} else if i > 0 {
					horses[i], horses[i-1] = horses[i-1], horses[i]
				}
			}
		}

		var total int64
		for i := 0; i < n; i++ {
			total += w[wOrder[i]] * horses[i].val
		}
		fmt.Fprintln(out, total)
	}
}
