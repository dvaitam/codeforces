package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const mod int64 = 998244353

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

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	// Precompute factorials for combinations up to n.
	fact := make([]int64, n+1)
	invFact := make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	invFact[n] = modPow(fact[n], mod-2)
	for i := n; i >= 1; i-- {
		invFact[i-1] = invFact[i] * int64(i) % mod
	}
	comb := func(N, K int) int64 {
		if K < 0 || K > N {
			return 0
		}
		return fact[N] * invFact[K] % mod * invFact[N-K] % mod
	}

	b := make([]int64, n)
	freq := make(map[int64]int)
	for i := 0; i < n; i++ {
		b[i] = a[i] + int64(i+1) // i is 1-based in formula
		freq[b[i]]++
	}

	sort.Slice(b, func(i, j int) bool { return b[i] > b[j] })

	prefixSum := int64(0)
	maxVal := int64(-1 << 60)
	waysAns := int64(0)

	// k = 0: train every day.
	maxVal = 0
	waysAns = 1

	prefCnt := make(map[int64]int)
	for k := 1; k <= n; k++ {
		prefixSum += b[k-1]
		val := prefixSum - int64(k*(k+1)/2)
		prefCnt[b[k-1]]++
		need := prefCnt[b[k-1]]
		total := freq[b[k-1]]
		ways := comb(total, need)

		if val > maxVal {
			maxVal = val
			waysAns = ways
		} else if val == maxVal {
			waysAns = (waysAns + ways) % mod
		}
	}

	fmt.Fprintf(out, "%d %d\n", maxVal, waysAns%mod)
}
