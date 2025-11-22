package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

type fenwick struct {
	n int
	f []int
}

func newFenwick(n int) *fenwick {
	return &fenwick{n: n, f: make([]int, n+1)}
}

func (ft *fenwick) add(idx, delta int) {
	for idx <= ft.n {
		ft.f[idx] += delta
		idx += idx & -idx
	}
}

func (ft *fenwick) sum(idx int) int {
	res := 0
	for idx > 0 {
		res += ft.f[idx]
		idx -= idx & -idx
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	s := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &s[i])
	}

	// Coordinate compression for Fenwick queries on values.
	all := make([]int64, n)
	copy(all, s)
	sort.Slice(all, func(i, j int) bool { return all[i] < all[j] })
	all = unique(all)

	// bad[i] counts descents between indices < i (size n)
	bad := make([]int, n)
	lastBad := make([]int, n)
	last := -1
	bad[0] = 0
	for i := 1; i < n; i++ {
		if s[i-1] > s[i] {
			last = i - 1
			bad[i] = bad[i-1] + 1
		} else {
			bad[i] = bad[i-1]
		}
		lastBad[i] = last
	}
	lastBad[0] = -1

	prefMax := make([]int64, n)
	prefMax[0] = s[0]
	for i := 1; i < n; i++ {
		if s[i] > prefMax[i-1] {
			prefMax[i] = s[i]
		} else {
			prefMax[i] = prefMax[i-1]
		}
	}

	// prefOk[i]: after sorting prefix ending at i-1, boundary to index i is fine (max prefix <= s[i])
	prefOk := make([]bool, n)
	prefOk[0] = true
	for i := 1; i < n; i++ {
		prefOk[i] = prefMax[i-1] <= s[i]
	}

	// nextOk[i]: smallest idx >= i with prefOk[idx] true
	nextOk := make([]int, n+1)
	nextOk[n] = n
	for i := n - 1; i >= 0; i-- {
		if prefOk[i] {
			nextOk[i] = i
		} else {
			nextOk[i] = nextOk[i+1]
		}
	}

	rightMin := make([]int64, n+1)
	rightMin[n] = math.MaxInt64
	for i := n - 1; i >= 0; i-- {
		if s[i] < rightMin[i+1] {
			rightMin[i] = s[i]
		} else {
			rightMin[i] = rightMin[i+1]
		}
	}

	best := int64(math.MaxInt64)

	// Helper to compute prefix query index for value <= x
	getLEIdx := func(x int64) int {
		// returns index in Fenwick (1-based) for last value <= x
		pos := sort.Search(len(all), func(i int) bool { return all[i] > x })
		return pos // pos equals count of <= x, OK for Fenwick sum
	}
	getLTIdx := func(x int64) int {
		pos := sort.Search(len(all), func(i int) bool { return all[i] >= x })
		return pos
	}

	// Case 1: non-overlap/middle untouched (order irrelevant)
	for k := 0; k <= n; k++ { // k is start index of suffix sort (length n-k)
		sufLen := n - k
		sufCost := int64(sufLen * sufLen)
		sufMin := rightMin[k]

		// Case: empty middle (i == k)
		if k == 0 {
			// prefix length 0
			cost := sufCost
			if cost < best {
				best = cost
			}
		} else {
			// need maxPrefix <= sufMin
			if prefMax[k-1] <= sufMin {
				cost := int64(k*k) + sufCost
				if cost < best {
					best = cost
				}
			}
		}

		if k == 0 {
			continue // no non-empty middle possible
		}

		// Case: non-empty middle [i..k-1] untouched
		if sufMin < s[k-1] {
			continue
		}
		i0 := lastBad[k-1] + 1
		if i0 >= k {
			continue
		}
		iCand := nextOk[i0]
		if iCand >= k {
			continue
		}
		cost := int64(iCand*iCand) + sufCost
		if cost < best {
			best = cost
		}
	}

	// Case 2: order = suffix first, prefix last (overlap allowed)
	// Sweep split from n down to 0 while maintaining suffix elements in BIT.
	ftSuf := newFenwick(len(all))
	for split := n; split >= 0; split-- { // split = number of elements before suffix start (n-b)
		if split < n {
			idx := getLEIdx(s[split]) // position for value <= val equals upper bound
			ftSuf.add(idx, 1)
		}
		sufLen := n - split
		maxPref := int64(math.MinInt64)
		if split > 0 {
			maxPref = prefMax[split-1]
		}
		cntSmaller := 0
		if split > 0 {
			idxLT := getLTIdx(maxPref)
			cntSmaller = ftSuf.sum(idxLT)
		}
		a := split + cntSmaller
		if a > n {
			continue
		}
		cost := int64(a*a + sufLen*sufLen)
		if cost < best {
			best = cost
		}
	}

	// Case 3: order = prefix first, suffix last (overlap allowed)
	// Sweep split from 0 to n, maintaining prefix in BIT and pointer aPtr.
	ftPref := newFenwick(len(all))
	aPtr := 0                             // current prefix length built in BIT
	for split := 0; split <= n; split++ { // split = prefix-only length (n-b)
		thr := rightMin[split] // min of suffix starting at split
		// Ensure prefix length at least split
		for aPtr < split {
			idx := getLEIdx(s[aPtr])
			ftPref.add(idx, 1)
			aPtr++
		}
		// Grow prefix until enough elements <= thr
		for aPtr < n {
			cnt := ftPref.sum(getLEIdx(thr))
			if cnt >= split {
				break
			}
			idx := getLEIdx(s[aPtr])
			ftPref.add(idx, 1)
			aPtr++
		}
		cnt := ftPref.sum(getLEIdx(thr))
		if cnt >= split {
			prefLen := aPtr
			sufLen := n - split
			cost := int64(prefLen*prefLen + sufLen*sufLen)
			if cost < best {
				best = cost
			}
		}
	}

	fmt.Fprintln(out, best)
}

func unique(a []int64) []int64 {
	if len(a) == 0 {
		return a
	}
	j := 1
	for i := 1; i < len(a); i++ {
		if a[i] != a[i-1] {
			a[j] = a[i]
			j++
		}
	}
	return a[:j]
}
