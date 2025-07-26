package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

const inf int64 = 1<<63 - 1

func cost(L, p int64) int64 {
	if p <= 0 {
		return inf
	}
	q := L / p
	r := L % p
	qq := uint64(q * q)
	hi1, lo1 := bits.Mul64(qq, uint64(p-r))
	qq1 := uint64(q+1) * uint64(q+1)
	hi2, lo2 := bits.Mul64(qq1, uint64(r))
	lo, carry := bits.Add64(lo1, lo2, 0)
	hi, _ := bits.Add64(hi1, hi2, carry)
	if hi > 0 || lo > uint64(inf) {
		return inf
	}
	return int64(lo)
}

func delta(L, p int64) int64 {
	c1 := cost(L, p)
	c2 := cost(L, p+1)
	if c1 == inf || c2 == inf {
		return inf
	}
	return c1 - c2
}

func minParts(L, lam int64) int64 {
	low, high := int64(1), L
	for low < high {
		mid := (low + high) >> 1
		if delta(L, mid) <= lam {
			high = mid
		} else {
			low = mid + 1
		}
	}
	return low
}

func feasible(lengths []int64, lam, m int64) bool {
	var total int64
	for _, L := range lengths {
		p := minParts(L, lam)
		total += cost(L, p)
		if total > m {
			return false
		}
	}
	return total <= m
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	a := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}
	var m int64
	fmt.Fscan(in, &m)

	lengths := make([]int64, n)
	prev := int64(0)
	for i := 0; i < n; i++ {
		lengths[i] = a[i+1] - prev
		prev = a[i+1]
	}

	lo, hi := int64(0), int64(1e18)
	best := int64(0)
	for lo <= hi {
		mid := (lo + hi) >> 1
		if feasible(lengths, mid, m) {
			best = mid
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}

	parts := make([]int64, len(lengths))
	var costNow int64
	for i, L := range lengths {
		p := minParts(L, best)
		parts[i] = p
		costNow += cost(L, p)
	}
	K := int64(0)
	for _, p := range parts {
		K += p - 1
	}
	leftover := m - costNow

	type item struct {
		d   int64
		idx int
	}
	pq := make([]item, 0)
	for i, L := range lengths {
		if parts[i] > 1 {
			d := delta(L, parts[i]-1)
			pq = append(pq, item{d: d, idx: i})
		}
	}
	sort.Slice(pq, func(i, j int) bool { return pq[i].d < pq[j].d })
	for len(pq) > 0 {
		it := pq[0]
		if it.d > leftover {
			break
		}
		leftover -= it.d
		K--
		idx := it.idx
		parts[idx]--
		if parts[idx] > 1 {
			d := delta(lengths[idx], parts[idx]-1)
			pq[0] = item{d: d, idx: idx}
		} else {
			pq = pq[1:]
		}
		sort.Slice(pq, func(i, j int) bool { return pq[i].d < pq[j].d })
	}
	fmt.Println(K)
}
