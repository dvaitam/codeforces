package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD = 998244353

func mul(a, b int) int {
	return int(int64(a) * int64(b) % int64(MOD))
}

func add(a, b int) int {
	r := a + b
	if r >= MOD {
		r -= MOD
	}
	return r
}

func sub(a, b int) int {
	r := a - b
	if r < 0 {
		r += MOD
	}
	return r
}

func modPow(a, e int) int {
	res := 1
	a %= MOD
	if a < 0 {
		a += MOD
	}
	for e > 0 {
		if e&1 == 1 {
			res = mul(res, a)
		}
		a = mul(a, a)
		e >>= 1
	}
	return res
}

type Fenwick struct {
	n    int
	tree []int
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{n, make([]int, n+1)}
}
func (f *Fenwick) Add(i, v int) {
	for ; i <= f.n; i += i & -i {
		f.tree[i] += v
	}
}
func (f *Fenwick) Sum(i int) int {
	s := 0
	for ; i > 0; i -= i & -i {
		s += f.tree[i]
	}
	return s
}
func (f *Fenwick) RangeSum(l, r int) int {
	if r < l {
		return 0
	}
	return f.Sum(r) - f.Sum(l-1)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	a := make([][]int, n)
	for i := 0; i < n; i++ {
		a[i] = make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &a[i][j])
		}
	}

	maxN := 2*n + 2
	fact := make([]int, maxN+1)
	fact[0] = 1
	for i := 1; i <= maxN; i++ {
		fact[i] = mul(fact[i-1], i)
	}
	inv_fact := make([]int, maxN+1)
	inv_fact[maxN] = modPow(fact[maxN], MOD-2)
	for i := maxN - 1; i >= 0; i-- {
		inv_fact[i] = mul(inv_fact[i+1], i+1)
	}

	comb := func(nn, k int) int {
		if k < 0 || k > nn || nn < 0 {
			return 0
		}
		return mul(fact[nn], mul(inv_fact[k], inv_fact[nn-k]))
	}

	// D(m, b) = permutations of m elements with b specific forbidden constraints
	// = sum_{i=0}^{b} (-1)^i * C(b,i) * (m-i)!
	dRestrict := func(m, b int) int {
		if m < 0 {
			return 0
		}
		if b > m {
			b = m
		}
		res := 0
		for i := 0; i <= b; i++ {
			if m-i < 0 {
				break
			}
			term := mul(comb(b, i), fact[m-i])
			if i%2 == 0 {
				res = add(res, term)
			} else {
				res = sub(res, term)
			}
		}
		return res
	}

	g := make([]int, n+1)
	g[0] = 1
	if n >= 1 {
		g[1] = 0
	}
	for i := 2; i <= n; i++ {
		g[i] = mul(i-1, add(g[i-1], g[i-2]))
	}
	G := g[n]
	gpows := make([]int, n+1)
	gpows[0] = 1
	for i := 1; i <= n; i++ {
		gpows[i] = mul(gpows[i-1], G)
	}

	ans := 0

	// First row: standard permutation rank
	{
		bit := NewFenwick(n)
		for i := 1; i <= n; i++ {
			bit.Add(i, 1)
		}
		rank := 0
		for j := 0; j < n; j++ {
			x := a[0][j]
			less := bit.RangeSum(1, x-1)
			rank = add(rank, mul(less%MOD, fact[n-j-1]))
			bit.Add(x, -1)
		}
		if n > 1 {
			ans = add(ans, mul(rank, gpows[n-1]))
		}
	}

	// Subsequent rows: derangement rank relative to previous row
	for i := 1; i < n; i++ {
		prev := a[i-1]
		curr := a[i]
		usedSet := make([]bool, n+1)

		// activeForbidden: for each future position k > j, prev[k] is forbidden at position k.
		// It's "active" if prev[k] is still unused.
		// We track this set explicitly using a Fenwick tree for counting < x queries.
		inActive := make([]bool, n+1) // whether value v is in the active forbidden set
		activeF := NewFenwick(n)      // Fenwick for active forbidden values
		unusedF := NewFenwick(n)      // Fenwick for unused values
		for v := 1; v <= n; v++ {
			unusedF.Add(v, 1)
		}
		// Initially, all positions 0..n-1 have active forbidden values
		for k := 0; k < n; k++ {
			if !inActive[prev[k]] {
				inActive[prev[k]] = true
				activeF.Add(prev[k], 1)
			}
		}

		rank := 0
		for j := 0; j < n; j++ {
			x := curr[j]
			forbidden := prev[j]

			// Position j's constraint is prev[j]. Remove it from active set.
			if inActive[forbidden] {
				inActive[forbidden] = false
				activeF.Add(forbidden, -1)
			}

			// Count valid values v < x: unused, v != forbidden
			lessUnused := unusedF.RangeSum(1, x-1)
			forbiddenSubtract := 0
			if !usedSet[forbidden] && forbidden < x {
				forbiddenSubtract = 1
			}
			validCount := lessUnused - forbiddenSubtract

			// Among these valid values, some are in activeF (placing them reduces active count)
			lessActive := activeF.RangeSum(1, x-1)
			// forbidden was already removed from activeF, so no need to subtract it
			inActiveCount := lessActive
			notInActiveCount := validCount - inActiveCount

			remaining := n - j - 1
			totalActive := activeF.RangeSum(1, n)

			comp0 := dRestrict(remaining, totalActive)   // v not in active: active stays same
			comp1 := dRestrict(remaining, totalActive-1)  // v in active: one fewer active

			contribution := add(mul(notInActiveCount%MOD, comp0), mul(inActiveCount%MOD, comp1))
			rank = add(rank, contribution)

			// Place x
			usedSet[x] = true
			unusedF.Add(x, -1)
			if inActive[x] {
				inActive[x] = false
				activeF.Add(x, -1)
			}
		}

		rem := n - 1 - i
		if rem >= 0 {
			ans = add(ans, mul(rank, gpows[rem]))
		}
	}
	fmt.Println(ans)
}
