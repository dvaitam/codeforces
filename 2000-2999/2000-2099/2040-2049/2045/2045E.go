package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const mod int64 = 998244353

type fenwick struct {
	n   int
	bit []int
}

func newFenwick(n int) *fenwick {
	return &fenwick{n: n, bit: make([]int, n+2)}
}

func (f *fenwick) add(idx, val int) {
	for idx <= f.n {
		f.bit[idx] += val
		idx += idx & -idx
	}
}

func (f *fenwick) sum(idx int) int {
	res := 0
	for idx > 0 {
		res += f.bit[idx]
		idx -= idx & -idx
	}
	return res
}

// find smallest index such that prefix sum >= k (k >= 1)
func (f *fenwick) kth(k int) int {
	idx := 0
	bitMask := 1
	for bitMask<<1 <= f.n {
		bitMask <<= 1
	}
	for bitMask > 0 {
		next := idx + bitMask
		if next <= f.n && f.bit[next] < k {
			k -= f.bit[next]
			idx = next
		}
		bitMask >>= 1
	}
	return idx + 1
}

func modAdd(a, b int64) int64 {
	res := a + b
	if res >= mod {
		res -= mod
	}
	return res
}

func modSub(a, b int64) int64 {
	res := a - b
	if res < 0 {
		res += mod
	}
	return res
}

func precomputePows(n int) ([]int64, []int64, []int64) {
	pow2 := make([]int64, n+2)
	pref := make([]int64, n+2)
	prefPref := make([]int64, n+2)
	inv2 := (mod + 1) / 2
	pow2[0] = 1
	for i := 1; i <= n+1; i++ {
		pow2[i] = pow2[i-1] * inv2 % mod
	}
	pref[0] = pow2[0]
	for i := 1; i <= n+1; i++ {
		pref[i] = modAdd(pref[i-1], pow2[i])
	}
	prefPref[0] = pref[0]
	for i := 1; i <= n+1; i++ {
		prefPref[i] = modAdd(prefPref[i-1], pref[i])
	}
	return pow2, pref, prefPref
}

// sum_{k=0}^{a-1} sum_{t=0}^{b-1} pow2[1+k+t]
func rectSum(a, b int, prefPref []int64) int64 {
	if a <= 0 || b <= 0 {
		return 0
	}
	idx1 := a - 1 + b
	idx2 := b - 1
	res := prefPref[idx1]
	if idx2 >= 0 {
		res = modSub(res, prefPref[idx2])
	}
	res = modSub(res, prefPref[a-1])
	return res
}

func powRange(l, r int, pref []int64) int64 {
	if l > r {
		return 0
	}
	res := pref[r]
	if l > 0 {
		res = modSub(res, pref[l-1])
	}
	return res
}

// prev greater (strict) and next greater or equal indices (sentinels 0, n+1)
func maxBounds(arr []int) ([]int, []int) {
	n := len(arr) - 1
	prev := make([]int, n+2)
	next := make([]int, n+2)

	stack := make([]int, 0, n)
	for i := 1; i <= n; i++ {
		for len(stack) > 0 && arr[stack[len(stack)-1]] <= arr[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			prev[i] = 0
		} else {
			prev[i] = stack[len(stack)-1]
		}
		stack = append(stack, i)
	}

	stack = stack[:0]
	for i := n; i >= 1; i-- {
		for len(stack) > 0 && arr[stack[len(stack)-1]] < arr[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			next[i] = n + 1
		} else {
			next[i] = stack[len(stack)-1]
		}
		stack = append(stack, i)
	}
	return prev, next
}

// For each position in primary, find nearest positions in other with value > primary value
func nearestGreaterOther(primary []int, other []int) ([]int, []int) {
	n := len(primary) - 1
	type pair struct {
		v int
		i int
	}
	prim := make([]pair, n)
	oth := make([]pair, n)
	for i := 1; i <= n; i++ {
		prim[i-1] = pair{primary[i], i}
		oth[i-1] = pair{other[i], i}
	}
	sort.Slice(prim, func(i, j int) bool { return prim[i].v > prim[j].v })
	sort.Slice(oth, func(i, j int) bool { return oth[i].v > oth[j].v })

	bit := newFenwick(n)
	resL := make([]int, n+2)
	resR := make([]int, n+2)
	ptr := 0
	for _, p := range prim {
		v := p.v
		for ptr < n && oth[ptr].v > v {
			bit.add(oth[ptr].i, 1)
			ptr++
		}
		idx := p.i
		cntLeft := bit.sum(idx - 1)
		if cntLeft == 0 {
			resL[idx] = 0
		} else {
			resL[idx] = bit.kth(cntLeft)
		}
		total := bit.sum(n)
		cntHere := bit.sum(idx)
		if total == cntHere {
			resR[idx] = n + 1
		} else {
			resR[idx] = bit.kth(cntHere + 1)
		}
	}
	for i := 1; i <= n; i++ {
		if other[i] > primary[i] {
			resL[i] = i
			resR[i] = i
		}
	}
	return resL, resR
}

func weightSum(i, L, R, n int, pref []int64, prefPref []int64, pow2 []int64) int64 {
	if L >= i || R <= i {
		return 0
	}
	left := i - L
	right := R - i
	base := rectSum(left, right, prefPref)
	addLeft := int64(0)
	addRight := int64(0)
	addBoth := int64(0)
	if L == 0 {
		addLeft = powRange(i, i+right-1, pref)
	}
	if R == n+1 {
		addRight = powRange(n-i+1, n-i+left, pref)
	}
	if L == 0 && R == n+1 {
		addBoth = pow2[n]
	}
	inv4 := (mod + 1) / 2
	inv4 = inv4 * inv4 % mod
	total := base
	total = modAdd(total, addLeft)
	total = modAdd(total, addRight)
	total = modAdd(total, addBoth)
	total = total * inv4 % mod
	return total
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	a := make([]int, n+1)
	b := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &b[i])
	}

	pow2, pref, prefPref := precomputePows(n)

	prevA, nextA := maxBounds(a)
	prevB, nextB := maxBounds(b)

	gaL, gaR := nearestGreaterOther(a, b)
	gbL, gbR := nearestGreaterOther(b, a)

	res := int64(0)
	for i := 1; i <= n; i++ {
		total := weightSum(i, prevA[i], nextA[i], n, pref, prefPref, pow2)
		l2 := prevA[i]
		if gaL[i] > l2 {
			l2 = gaL[i]
		}
		r2 := nextA[i]
		if gaR[i] < r2 {
			r2 = gaR[i]
		}
		restricted := weightSum(i, l2, r2, n, pref, prefPref, pow2)
		contrib := modSub(total, restricted)
		res = (res + int64(a[i])*contrib) % mod
	}

	for i := 1; i <= n; i++ {
		total := weightSum(i, prevB[i], nextB[i], n, pref, prefPref, pow2)
		l2 := prevB[i]
		if gbL[i] > l2 {
			l2 = gbL[i]
		}
		r2 := nextB[i]
		if gbR[i] < r2 {
			r2 = gbR[i]
		}
		restricted := weightSum(i, l2, r2, n, pref, prefPref, pow2)
		contrib := modSub(total, restricted)
		res = (res + int64(b[i])*contrib) % mod
	}

	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, res)
	out.Flush()
}
