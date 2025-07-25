package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const MOD int = 998244353
const ROOT int = 3

func modAdd(a, b int) int {
	s := a + b
	if s >= MOD {
		s -= MOD
	}
	return s
}

func modSub(a, b int) int {
	s := a - b
	if s < 0 {
		s += MOD
	}
	return s
}

func modMul(a, b int) int {
	return int((int64(a) * int64(b)) % int64(MOD))
}

func modPow(a, e int) int {
	res := 1
	x := a % MOD
	for e > 0 {
		if e&1 != 0 {
			res = modMul(res, x)
		}
		x = modMul(x, x)
		e >>= 1
	}
	return res
}

func modInv(a int) int {
	return modPow(a, MOD-2)
}

func ntt(a []int, invert bool) {
	n := len(a)
	j := 0
	for i := 1; i < n; i++ {
		bit := n >> 1
		for ; j&bit != 0; bit >>= 1 {
			j ^= bit
		}
		j |= bit
		if i < j {
			a[i], a[j] = a[j], a[i]
		}
	}
	for length := 2; length <= n; length <<= 1 {
		wlen := modPow(ROOT, (MOD-1)/length)
		if invert {
			wlen = modInv(wlen)
		}
		for i := 0; i < n; i += length {
			w := 1
			half := length >> 1
			for j := 0; j < half; j++ {
				u := a[i+j]
				v := modMul(a[i+j+half], w)
				a[i+j] = modAdd(u, v)
				a[i+j+half] = modSub(u, v)
				w = modMul(w, wlen)
			}
		}
	}
	if invert {
		invN := modInv(n)
		for i := range a {
			a[i] = modMul(a[i], invN)
		}
	}
}

func polyMul(a, b []int) []int {
	n := len(a) + len(b) - 1
	sz := 1
	for sz < n {
		sz <<= 1
	}
	fa := make([]int, sz)
	fb := make([]int, sz)
	copy(fa, a)
	copy(fb, b)
	ntt(fa, false)
	ntt(fb, false)
	for i := 0; i < sz; i++ {
		fa[i] = modMul(fa[i], fb[i])
	}
	ntt(fa, true)
	return fa[:n]
}

func prepareFactorials(maxN int) ([]int, []int) {
	fact := make([]int, maxN+1)
	invFact := make([]int, maxN+1)
	fact[0] = 1
	for i := 1; i <= maxN; i++ {
		fact[i] = modMul(fact[i-1], i)
	}
	invFact[maxN] = modInv(fact[maxN])
	for i := maxN; i > 0; i-- {
		invFact[i-1] = modMul(invFact[i], i)
	}
	return fact, invFact
}

func comb(n, k int, fact, invFact []int) int {
	if k < 0 || k > n {
		return 0
	}
	return modMul(fact[n], modMul(invFact[k], invFact[n-k]))
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}

	const MAXL = 300000
	freq := make([]int, MAXL+2)
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		if x <= MAXL {
			freq[x]++
		}
	}

	reds := make([]int, k)
	for i := 0; i < k; i++ {
		fmt.Fscan(reader, &reds[i])
	}

	var q int
	fmt.Fscan(reader, &q)
	queries := make([]int, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(reader, &queries[i])
	}

	one := make([]int, MAXL+2)
	two := make([]int, MAXL+2)
	for i := 1; i <= MAXL; i++ {
		if freq[i] == 1 {
			one[i] = 1
		} else if freq[i] >= 2 {
			two[i] = 1
		}
	}
	prefixOne := make([]int, MAXL+2)
	prefixTwo := make([]int, MAXL+2)
	for i := 1; i <= MAXL; i++ {
		prefixOne[i] = prefixOne[i-1] + one[i]
		prefixTwo[i] = prefixTwo[i-1] + two[i]
	}

	sort.Ints(reds)
	maxB := 0
	for _, r := range reds {
		if prefixTwo[r-1] > maxB {
			maxB = prefixTwo[r-1]
		}
	}
	maxN := 2 * n
	if 2*maxB > maxN {
		maxN = 2 * maxB
	}

	fact, invFact := prepareFactorials(maxN)
	pow2 := make([]int, n+1)
	pow2[0] = 1
	for i := 1; i <= n; i++ {
		pow2[i] = modAdd(pow2[i-1], pow2[i-1])
	}

	type polyInfo struct {
		r   int
		val []int
	}
	polys := make([]polyInfo, k)
	for idx, r := range reds {
		A := prefixOne[r-1]
		B := prefixTwo[r-1]
		poly1 := make([]int, A+1)
		for i := 0; i <= A; i++ {
			poly1[i] = modMul(comb(A, i, fact, invFact), pow2[i])
		}
		poly2 := make([]int, 2*B+1)
		for j := 0; j <= 2*B; j++ {
			poly2[j] = comb(2*B, j, fact, invFact)
		}
		conv := polyMul(poly1, poly2)
		polys[idx] = polyInfo{r: r, val: conv}
	}

	for _, Q := range queries {
		ans := 0
		half := Q / 2
		for _, p := range polys {
			w := half - p.r - 1
			if w >= 0 && w < len(p.val) {
				ans = modAdd(ans, p.val[w])
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
