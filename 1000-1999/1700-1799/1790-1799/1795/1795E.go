package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Fenwick struct {
	n   int
	bit []int64
}

func newFenwick(n int) *Fenwick {
	return &Fenwick{n: n, bit: make([]int64, n+2)}
}

func (f *Fenwick) add(i int, delta int64) {
	for i <= f.n {
		f.bit[i] += delta
		i += i & -i
	}
}

func (f *Fenwick) sum(i int) int64 {
	s := int64(0)
	for i > 0 {
		s += f.bit[i]
		i -= i & -i
	}
	return s
}

func (f *Fenwick) rangeSum(l, r int) int64 {
	if r < l {
		return 0
	}
	if l <= 1 {
		return f.sum(r)
	}
	return f.sum(r) - f.sum(l-1)
}

func compress(values []int64) ([]int64, map[int64]int) {
	uniq := make([]int64, len(values))
	copy(uniq, values)
	sort.Slice(uniq, func(i, j int) bool { return uniq[i] < uniq[j] })
	m := 1
	for i := 1; i < len(uniq); i++ {
		if uniq[i] != uniq[m-1] {
			uniq[m] = uniq[i]
			m++
		}
	}
	uniq = uniq[:m]
	mp := make(map[int64]int, len(uniq))
	for idx, v := range uniq {
		mp[v] = idx + 1 // fenwick is 1-indexed
	}
	return uniq, mp
}

func solve() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		h := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &h[i])
		}
		if n == 1 {
			fmt.Fprintln(out, h[0])
			continue
		}

		A := make([]int64, n)
		B := make([]int64, n)
		for i := 0; i < n; i++ {
			idx := int64(i + 1)
			A[i] = h[i] + idx
			B[i] = h[i] - idx
		}
		valsA, mapA := compress(A)
		valsB, mapB := compress(B)

		bitCntRight := newFenwick(len(valsA))
		bitSumRight := newFenwick(len(valsA))

		right := make([]int64, n)
		for i := n - 1; i >= 0; i-- {
			if i < n-1 {
				v := A[i+1]
				pos := mapA[v]
				bitCntRight.add(pos, 1)
				bitSumRight.add(pos, v)
			}
			t := A[i]
			// query values > t
			idx := sort.Search(len(valsA), func(j int) bool { return valsA[j] > t })
			if idx < len(valsA) {
				sum := bitSumRight.rangeSum(idx+1, len(valsA))
				cnt := bitCntRight.rangeSum(idx+1, len(valsA))
				right[i] = sum - int64(t)*cnt
			}
		}

		bitCntLeft := newFenwick(len(valsB))
		bitSumLeft := newFenwick(len(valsB))
		left := make([]int64, n)
		for i := 0; i < n; i++ {
			if i > 0 {
				v := B[i-1]
				pos := mapB[v]
				bitCntLeft.add(pos, 1)
				bitSumLeft.add(pos, v)
			}
			t := B[i]
			idx := sort.Search(len(valsB), func(j int) bool { return valsB[j] > t })
			if idx < len(valsB) {
				sum := bitSumLeft.rangeSum(idx+1, len(valsB))
				cnt := bitCntLeft.rangeSum(idx+1, len(valsB))
				left[i] = sum - int64(t)*cnt
			}
		}

		ans := int64(1 << 60)
		for i := 0; i < n; i++ {
			cost := h[i] + left[i] + right[i]
			if cost < ans {
				ans = cost
			}
		}
		fmt.Fprintln(out, ans)
	}
}

func main() {
	solve()
}
