package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func powMod(a, e int64) int64 {
	a %= mod
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

	var N, M int
	if _, err := fmt.Fscan(in, &N, &M); err != nil {
		return
	}
	D := make([]int64, N)
	total := int64(0)
	for i := 0; i < N; i++ {
		fmt.Fscan(in, &D[i])
		total += D[i]
	}

	if total%2 == 1 {
		fmt.Fprintln(out, powMod(int64(M), int64(N))%mod)
		return
	}
	half := total / 2
	pos := make([]int64, N)
	cur := int64(0)
	for i := 0; i < N; i++ {
		pos[i] = cur
		cur += D[i]
	}

	index := make(map[int64]int, N)
	for i := 0; i < N; i++ {
		index[pos[i]] = i
	}
	P := 0
	for i := 0; i < N; i++ {
		target := (pos[i] + half) % total
		if j, ok := index[target]; ok && i < j {
			P++
		}
	}
	U := N - 2*P

	size := M
	if P > size {
		size = P
	}
	fac := make([]int64, size+1)
	fac[0] = 1
	for i := 1; i <= size; i++ {
		fac[i] = fac[i-1] * int64(i) % mod
	}
	invFac := make([]int64, size+1)
	invFac[size] = powMod(fac[size], mod-2)
	for i := size; i > 0; i-- {
		invFac[i-1] = invFac[i] * int64(i) % mod
	}

	comb := func(n, k int) int64 {
		if k < 0 || k > n {
			return 0
		}
		return fac[n] * invFac[k] % mod * invFac[n-k] % mod
	}
	perm := func(n, k int) int64 {
		if k < 0 || k > n {
			return 0
		}
		return fac[n] * invFac[n-k] % mod
	}

	maxT := P
	if M < maxT {
		maxT = M
	}
	ans := int64(0)
	for t := 0; t <= maxT; t++ {
		if P-t > 0 && M-t < 2 {
			continue
		}
		ways := comb(P, t)
		ways = ways * perm(M, t) % mod
		ways = ways * powMod(int64(M-t), int64(U+P-t)) % mod
		ways = ways * powMod(int64(M-t-1), int64(P-t)) % mod
		ans = (ans + ways) % mod
	}

	fmt.Fprintln(out, ans)
}
