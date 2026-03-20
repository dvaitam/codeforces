package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func modPow(a, e, mod int64) int64 {
	r := int64(1)
	for e > 0 {
		if e&1 == 1 {
			r = r * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return r
}

func main() {
	in := bufio.NewReaderSize(os.Stdin, 1<<20)
	out := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer out.Flush()

	var n, k int
	var p int64
	fmt.Fscan(in, &n, &k, &p)

	total := modPow(int64(k), int64(n), p)
	if n&1 == 1 {
		fmt.Fprintln(out, total)
		return
	}

	fact := make([]int64, n+1)
	invFact := make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % p
	}
	invFact[n] = modPow(fact[n], p-2, p)
	for i := n; i >= 1; i-- {
		invFact[i-1] = invFact[i] * int64(i) % p
	}

	m := n / 2
	vals := make([]int, 0, 18)
	for x := m; x > 0; x &= x - 1 {
		vals = append(vals, x&-x)
	}

	s := len(vals)
	K := k
	if K > s {
		K = s
	}

	size := 1 << s
	sum := make([]int, size)
	pc := make([]uint8, size)
	for mask := 1; mask < size; mask++ {
		lb := mask & -mask
		idx := bits.TrailingZeros(uint(lb))
		prev := mask ^ lb
		sum[mask] = sum[prev] + vals[idx]
		pc[mask] = pc[prev] + 1
	}

	w := make([]int64, size)
	for mask := 1; mask < size; mask++ {
		w[mask] = invFact[2*sum[mask]]
	}

	stride := K + 1
	dp := make([]int64, size*stride)
	dp[0] = 1

	for mask := 1; mask < size; mask++ {
		lb := mask & -mask
		rem := mask ^ lb
		base := mask * stride
		for sub := rem; ; sub = (sub - 1) & rem {
			block := lb | sub
			prev := rem ^ sub
			pbase := prev * stride
			lim := int(pc[prev])
			if lim > K-1 {
				lim = K - 1
			}
			mul := w[block]
			for t := 0; t <= lim; t++ {
				v := dp[pbase+t]
				if v != 0 {
					dp[base+t+1] = (dp[base+t+1] + v*mul) % p
				}
			}
			if sub == 0 {
				break
			}
		}
	}

	fullBase := (size - 1) * stride
	lose := int64(0)
	perm := int64(1)
	for t := 1; t <= K; t++ {
		perm = perm * int64(k-t+1) % p
		lose = (lose + dp[fullBase+t]*perm) % p
	}
	lose = lose * fact[n] % p

	ans := total - lose
	if ans < 0 {
		ans += p
	}
	fmt.Fprintln(out, ans)
}