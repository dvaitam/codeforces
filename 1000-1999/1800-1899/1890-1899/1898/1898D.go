package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Solution for problemD.txt in folder 1898.
// For each test case we are given arrays a and b and may swap at most one
// pair of positions in b. If we denote x_i = a_i - b_i and s_i = a_i + b_i,
// the change in absolute beauty after swapping i and j equals
//   max(|x_i + x_j|, |s_i - s_j|) - |x_i| - |x_j|.
// Since |x_i + x_j| \le |x_i| + |x_j|, only the second term can give a
// positive improvement. Thus we maximize |s_i - s_j| - |x_i| - |x_j| over all
// pairs. After sorting indices by s_i we can maintain prefix maxima of q_j =
// s_j + |x_j| and suffix maxima of p_j = s_j - |x_j| to find this value in
// O(n log n). The final answer is the initial sum of |x_i| plus the best gain
// (or zero if all swaps are bad).

type item struct {
	s int64
	p int64
	q int64
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int64, n)
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		items := make([]item, n)
		var sum int64
		for i := 0; i < n; i++ {
			diff := a[i] - b[i]
			if diff < 0 {
				diff = -diff
			}
			sum += diff
			s := a[i] + b[i]
			items[i] = item{s: s, p: s - diff, q: s + diff}
		}
		sort.Slice(items, func(i, j int) bool { return items[i].s < items[j].s })
		pre := make([]int64, n)
		cur := int64(-1 << 60)
		for i := 0; i < n; i++ {
			if items[i].q > cur {
				cur = items[i].q
			}
			pre[i] = cur
		}
		suf := make([]int64, n)
		cur = int64(-1 << 60)
		for i := n - 1; i >= 0; i-- {
			if items[i].p > cur {
				cur = items[i].p
			}
			suf[i] = cur
		}
		var best int64
		for i := 0; i < n; i++ {
			if i > 0 {
				cand := items[i].p - pre[i-1]
				if cand > best {
					best = cand
				}
			}
			if i+1 < n {
				cand := suf[i+1] - items[i].q
				if cand > best {
					best = cand
				}
			}
		}
		if best < 0 {
			best = 0
		}
		fmt.Fprintln(out, sum+best)
	}
}
