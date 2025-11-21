package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type pair struct {
	sum int64
	i   int
	j   int
}

type query struct {
	sum int64
	l   int
	r   int
}

type fenwick struct {
	n   int
	bit []int
}

func newFenwick(n int) *fenwick {
	return &fenwick{n: n, bit: make([]int, n+1)}
}

func (f *fenwick) add(idx, delta int) {
	for idx <= f.n {
		f.bit[idx] += delta
		idx += idx & -idx
	}
}

func (f *fenwick) pref(idx int) int {
	res := 0
	for idx > 0 {
		res += f.bit[idx]
		idx -= idx & -idx
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var T int
	fmt.Fscan(in, &T)
	out := bufio.NewWriter(os.Stdout)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := range a {
			fmt.Fscan(in, &a[i])
		}

		pref := make([]int64, n+1)
		for i, v := range a {
			pref[i+1] = pref[i] + int64(v)
		}

		pos := make(map[int64]int, n+1)
		for i, v := range pref {
			pos[v] = i
		}

		// Collect pairs of internal boundaries.
		pairs := make([]pair, 0, n*(n-1)/2)
		if n >= 2 {
			for i := 1; i <= n-2; i++ {
				for j := i + 1; j <= n-1; j++ {
					pairs = append(pairs, pair{sum: pref[i] + pref[j], i: i, j: j})
				}
			}
		}

		// Collect queries for all subarrays and count middle contributions.
		queries := make([]query, 0, n*(n+1)/2)
		var midTotal int64
		for l := 1; l <= n; l++ {
			for r := l; r <= n; r++ {
				curSum := pref[l-1] + pref[r]
				queries = append(queries, query{sum: curSum, l: l, r: r})
				if curSum%2 == 0 {
					want := curSum / 2
					if p, ok := pos[want]; ok && p >= l && p <= r-1 {
						midTotal++
					}
				}
			}
		}

		// Sum of (len-1) over all subarrays.
		var base int64
		for length := 1; length <= n; length++ {
			base += int64(length-1) * int64(n-length+1)
		}

		sort.Slice(pairs, func(i, j int) bool { return pairs[i].sum < pairs[j].sum })
		sort.Slice(queries, func(i, j int) bool { return queries[i].sum < queries[j].sum })

		var totalPairs int64
		pi, qi := 0, 0
		for pi < len(pairs) || qi < len(queries) {
			var curSum int64
			if qi == len(queries) {
				curSum = pairs[pi].sum
			} else if pi == len(pairs) {
				curSum = queries[qi].sum
			} else if pairs[pi].sum < queries[qi].sum {
				curSum = pairs[pi].sum
			} else {
				curSum = queries[qi].sum
			}

			pStart := pi
			for pi < len(pairs) && pairs[pi].sum == curSum {
				pi++
			}
			qStart := qi
			for qi < len(queries) && queries[qi].sum == curSum {
				qi++
			}

			groupPairs := pairs[pStart:pi]
			groupQueries := queries[qStart:qi]
			if len(groupPairs) == 0 || len(groupQueries) == 0 {
				continue
			}

			sort.Slice(groupPairs, func(i, j int) bool { return groupPairs[i].i > groupPairs[j].i })
			sort.Slice(groupQueries, func(i, j int) bool { return groupQueries[i].l > groupQueries[j].l })

			ft := newFenwick(n)
			pIdx, qIdx := 0, 0
			for pIdx < len(groupPairs) || qIdx < len(groupQueries) {
				curL := 0
				if pIdx < len(groupPairs) {
					curL = groupPairs[pIdx].i
				}
				if qIdx < len(groupQueries) && groupQueries[qIdx].l > curL {
					curL = groupQueries[qIdx].l
				}

				for pIdx < len(groupPairs) && groupPairs[pIdx].i == curL {
					ft.add(groupPairs[pIdx].j, 1)
					pIdx++
				}
				for qIdx < len(groupQueries) && groupQueries[qIdx].l == curL {
					totalPairs += int64(ft.pref(groupQueries[qIdx].r - 1))
					qIdx++
				}
			}
		}

		ans := base - 2*totalPairs - midTotal
		fmt.Fprintln(out, ans)
	}
	out.Flush()
}
