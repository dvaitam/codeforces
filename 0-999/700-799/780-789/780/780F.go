package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

// Bitset is a simple fixed-size bitset implementation.
type bitset []uint64

func newBitset(n int) bitset {
	return make([]uint64, (n+63)>>6)
}

func (b bitset) set(i int) {
	b[i>>6] |= 1 << uint(i&63)
}

func (b bitset) or(src bitset) {
	for i := range b {
		b[i] |= src[i]
	}
}

func (b bitset) isZero() bool {
	for _, w := range b {
		if w != 0 {
			return false
		}
	}
	return true
}

func (b bitset) clone() bitset {
	dst := make(bitset, len(b))
	copy(dst, b)
	return dst
}

func forEach(b bitset, f func(i int)) {
	for idx, w := range b {
		for w != 0 {
			t := w & -w
			bit := bits.TrailingZeros64(w)
			pos := idx*64 + bit
			f(pos)
			w &^= t
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	edges := make([][]bitset, 2)
	edges[0] = make([]bitset, n)
	edges[1] = make([]bitset, n)
	for i := 0; i < n; i++ {
		edges[0][i] = newBitset(n)
		edges[1][i] = newBitset(n)
	}

	for i := 0; i < m; i++ {
		var v, u, t int
		fmt.Fscan(in, &v, &u, &t)
		v--
		u--
		edges[t][v].set(u)
	}

	const maxK = 60
	dp := make([][][]bitset, maxK+1)
	dp[0] = make([][]bitset, 2)
	dp[0][0] = make([]bitset, n)
	dp[0][1] = make([]bitset, n)
	for i := 0; i < n; i++ {
		dp[0][0][i] = edges[0][i].clone()
		dp[0][1][i] = edges[1][i].clone()
	}

	for k := 1; k <= maxK; k++ {
		dp[k] = make([][]bitset, 2)
		dp[k][0] = make([]bitset, n)
		dp[k][1] = make([]bitset, n)
		for v := 0; v < n; v++ {
			bs0 := newBitset(n)
			forEach(dp[k-1][0][v], func(x int) {
				if x < n {
					bs0.or(dp[k-1][1][x])
				}
			})
			dp[k][0][v] = bs0

			bs1 := newBitset(n)
			forEach(dp[k-1][1][v], func(x int) {
				if x < n {
					bs1.or(dp[k-1][0][x])
				}
			})
			dp[k][1][v] = bs1
		}
	}

	canReach := func(L int64) bool {
		cur := newBitset(n)
		cur.set(0)
		prefix := int64(0)
		for k := maxK; k >= 0; k-- {
			if (L>>uint(k))&1 == 1 {
				orient := bits.OnesCount64(uint64(prefix)) & 1
				nxt := newBitset(n)
				forEach(cur, func(v int) {
					nxt.or(dp[k][orient][v])
				})
				if nxt.isZero() {
					return false
				}
				cur = nxt
				prefix += 1 << uint(k)
			}
		}
		return !cur.isZero()
	}

	const limit = int64(1000000000000000000)
	if canReach(limit + 1) {
		fmt.Println(-1)
		return
	}

	low, high := int64(0), limit
	ans := int64(0)
	for low <= high {
		mid := (low + high) / 2
		if canReach(mid) {
			ans = mid
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	fmt.Println(ans)
}
