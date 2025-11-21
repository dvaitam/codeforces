package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Fenwick struct {
	n   int
	bit []int
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{n: n, bit: make([]int, n+2)}
}

func (f *Fenwick) Add(idx, delta int) {
	for idx <= f.n {
		f.bit[idx] += delta
		idx += idx & -idx
	}
}

func (f *Fenwick) Sum(idx int) int {
	if idx > f.n {
		idx = f.n
	}
	res := 0
	for idx > 0 {
		res += f.bit[idx]
		idx -= idx & -idx
	}
	return res
}

func upperBound(arr []int64, target int64) int {
	l, r := 0, len(arr)
	for l < r {
		mid := (l + r) >> 1
		if arr[mid] > target {
			r = mid
		} else {
			l = mid + 1
		}
	}
	return l
}

func lowerBound(arr []int64, target int64) int {
	l, r := 0, len(arr)
	for l < r {
		mid := (l + r) >> 1
		if arr[mid] >= target {
			r = mid
		} else {
			l = mid + 1
		}
	}
	return l
}

func makerCanWin(a []int64, comp []int, vals []int64, D int64) bool {
	n := len(a)
	m := len(vals)
	deg := make([]int, n)

	fw := NewFenwick(m)
	for i := 0; i < n; i++ {
		threshold := a[i] - D
		idx := upperBound(vals, threshold)
		if idx > 0 {
			deg[i] += fw.Sum(idx)
		}
		fw.Add(comp[i], 1)
	}

	fw = NewFenwick(m)
	total := 0
	for i := n - 1; i >= 0; i-- {
		threshold := a[i] + D
		idx := lowerBound(vals, threshold)
		less := fw.Sum(idx)
		countGE := total - less
		deg[i] += countGE
		fw.Add(comp[i], 1)
		total++
	}

	ge3 := make([]int, 0)
	deg2Exists := false
	for i, d := range deg {
		if d >= 3 {
			ge3 = append(ge3, i)
		}
		if d == 2 {
			deg2Exists = true
		}
	}

	adjacent := func(i, j int) bool {
		if i < j {
			return a[j]-a[i] >= D
		}
		return a[i]-a[j] >= D
	}

	if len(ge3) >= 2 {
		return true
	}

	if len(ge3) == 1 {
		v := ge3[0]
		for i := 0; i < n; i++ {
			if i == v {
				continue
			}
			if deg[i] >= 2 && !adjacent(i, v) {
				return true
			}
		}
		return false
	}

	if !deg2Exists {
		return false
	}

	const INF int64 = 1 << 60
	const NEG_INF int64 = -INF

	pref := make([]int64, n)
	prefVal := NEG_INF
	for i := 0; i < n; i++ {
		pref[i] = prefVal
		if deg[i] == 2 {
			candidate := a[i] + D
			if candidate > prefVal {
				prefVal = candidate
			}
		}
	}

	suff := make([]int64, n)
	suffVal := INF
	for i := n - 1; i >= 0; i-- {
		suff[i] = suffVal
		if deg[i] == 2 {
			candidate := a[i] - D
			if candidate < suffVal {
				suffVal = candidate
			}
		}
	}

	for i := 0; i < n; i++ {
		left := pref[i]
		right := suff[i]
		if a[i] >= left && a[i] <= right {
			return false
		}
	}
	return true
}

func solveCase(a []int64) int64 {
	n := len(a)
	vals := make([]int64, n)
	copy(vals, a)
	sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })
	uniq := vals[:0]
	for _, v := range vals {
		if len(uniq) == 0 || uniq[len(uniq)-1] != v {
			uniq = append(uniq, v)
		}
	}
	vals = uniq
	comp := make([]int, n)
	for i, v := range a {
		comp[i] = lowerBound(vals, v) + 1
	}

	lo := int64(-1_000_000_000)
	hi := int64(1_000_000_000)
	for lo < hi {
		mid := (lo + hi + 1) / 2
		if makerCanWin(a, comp, vals, mid) {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	return lo
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			var x int64
			fmt.Fscan(in, &x)
			a[i] = x
		}
		ans := solveCase(a)
		fmt.Fprintln(out, ans)
	}
}
