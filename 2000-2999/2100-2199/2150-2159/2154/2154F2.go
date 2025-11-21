package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 998244353

var fact []int64
var invFact []int64
var precomputed bool

func initFactorials(maxN int) {
	if precomputed {
		return
	}
	fact = make([]int64, maxN+1)
	invFact = make([]int64, maxN+1)
	fact[0] = 1
	for i := 1; i <= maxN; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	invFact[maxN] = modPow(fact[maxN], mod-2)
	for i := maxN; i >= 1; i-- {
		invFact[i-1] = invFact[i] * int64(i) % mod
	}
	precomputed = true
}

func modPow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func comb(n, k int) int64 {
	if k < 0 || k > n {
		return 0
	}
	return fact[n] * invFact[k] % mod * invFact[n-k] % mod
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)

	const maxN = 1_000_000
	initFactorials(maxN)

	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		pos := make([]int, n+1)
		for i := 1; i <= n; i++ {
			pos[i] = -1
		}
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &p[i])
			if p[i] != -1 {
				pos[p[i]] = i
			}
		}

		prefKnown := make([]int, n+1)
		prefixValid := make([]bool, n+1)
		last := -1
		ok := true
		for v := 1; v <= n; v++ {
			prefKnown[v] = prefKnown[v-1]
			if pos[v] != -1 {
				prefKnown[v]++
				if pos[v] > last {
					last = pos[v]
				} else {
					ok = false
				}
			}
			prefixValid[v] = ok
		}

	suffixValid := make([]bool, n+2)
	lastPos := n + 1
	ok = true
	for v := n; v >= 1; v-- {
		if pos[v] != -1 {
			if pos[v] < lastPos {
				lastPos = pos[v]
			} else {
				ok = false
			}
		}
		suffixValid[v] = ok
	}
	suffixValid[n+1] = true

	totalKnown := prefKnown[n]
	remainingPositions := n - totalKnown

	answer := int64(0)
	for k := 1; k <= n-1; k++ {
		if !prefixValid[k] || !suffixValid[k+1] {
			continue
		}
		forcedIn := prefKnown[k]
		need := k - forcedIn
		if need < 0 || need > remainingPositions {
			continue
		}
		answer = (answer + comb(remainingPositions, need)) % mod
	}
	fmt.Fprintln(out, answer)
	}
}
