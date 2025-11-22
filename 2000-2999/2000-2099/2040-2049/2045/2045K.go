package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 998244353

func modPow(a, e int) int {
	res := 1
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

	var N int
	if _, err := fmt.Fscan(in, &N); err != nil {
		return
	}
	freq := make([]int, N+1)
	for i := 1; i <= N; i++ {
		fmt.Fscan(in, &freq[i])
	}

	// mobius up to N
	mu := make([]int, N+1)
	primes := make([]int, 0)
	isComp := make([]bool, N+1)
	mu[1] = 1
	for i := 2; i <= N; i++ {
		if !isComp[i] {
			primes = append(primes, i)
			mu[i] = -1
		}
		for _, p := range primes {
			v := i * p
			if v > N {
				break
			}
			isComp[v] = true
			if i%p == 0 {
				mu[v] = 0
				break
			} else {
				mu[v] = -mu[i]
			}
		}
	}

	// cntMul[x] = count of cards divisible by x
	cntMul := make([]int, N+1)
	for x := 1; x <= N; x++ {
		for m := x; m <= N; m += x {
			cntMul[x] += freq[m]
		}
	}

	// precompute powers up to N
	pow2 := make([]int, N+1)
	pow3 := make([]int, N+1)
	powInv4 := make([]int, N+1)
	pow34 := make([]int, N+1) // (3/4)^k
	pow2[0], pow3[0], powInv4[0], pow34[0] = 1, 1, 1, 1
	inv4 := modPow(4, mod-2)
	for i := 1; i <= N; i++ {
		pow2[i] = pow2[i-1] * 2 % mod
		pow3[i] = pow3[i-1] * 3 % mod
		powInv4[i] = powInv4[i-1] * inv4 % mod
		pow34[i] = pow3[i] * powInv4[i] % mod
	}

	// reusable buffers
	c := make([]int, N+1)
	val2 := make([]int, N+1)
	S := make([]int, N+1)
	F := make([]int, N+1)
	g := make([]int, N+1)

	total := 0
	for i := 1; i <= N; i++ {
		if cntMul[i] < 2 {
			continue
		}
		len := N / i
		for j := 1; j <= len; j++ {
			c[j] = cntMul[i*j]
			val2[j] = pow2[c[j]]
			S[j], F[j], g[j] = 0, 0, 0
		}

		// S[l] = sum_{d|l} mu[d] * 2^{c[d]}
		for d := 1; d <= len; d++ {
			if mu[d] == 0 {
				continue
			}
			md := val2[d]
			if mu[d] == -1 {
				md = mod - md
			}
			for m := d; m <= len; m += d {
				S[m] += md
				if S[m] >= mod {
					S[m] -= mod
				}
			}
		}

		for l := 1; l <= len; l++ {
			F[l] = S[l] * S[l] % mod
			g[l] = 0
		}

		// g[l] = sum_{m|l} mu[l/m] * F[m]
		for m := 1; m <= len; m++ {
			fm := F[m]
			if fm == 0 {
				continue
			}
			for l := m; l <= len; l += m {
				muVal := mu[l/m]
				if muVal == 0 {
					continue
				}
				if muVal == 1 {
					g[l] += fm
				} else { // -1
					g[l] -= fm
				}
				if g[l] >= mod {
					g[l] -= mod
				} else if g[l] < 0 {
					g[l] += mod
				}
			}
		}

		ways := 0
		for l := 1; l <= len; l++ {
			if g[l] == 0 {
				continue
			}
			ways = (ways + g[l]*pow34[c[l]]) % mod
		}

		total = (total + i%mod*ways) % mod
	}

	fmt.Fprintln(out, total)
}
