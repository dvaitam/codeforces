package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 998244353

var inv2 = (mod + 1) / 2

func modPow(x, n int64) int64 {
	res := int64(1)
	for n > 0 {
		if n&1 == 1 {
			res = res * x % mod
		}
		x = x * x % mod
		n >>= 1
	}
	return res
}

func pairSum(sum, sumSq int64) int64 {
	diff := (sum*sum - sumSq) % mod
	if diff < 0 {
		diff += mod
	}
	return diff * int64(inv2) % mod
}

type edge struct {
	u int
	v int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)

		fallProb := make([]int64, n)
		stayProb := make([]int64, n)
		invFall := make([]int64, n)
		val := make([]int64, n)

		for i := 0; i < n; i++ {
			var p, q int64
			fmt.Fscan(in, &p, &q)
			p %= mod
			q %= mod
			invQ := modPow(q, mod-2)
			f := p * invQ % mod
			fallProb[i] = f
			stay := (1 - f + mod) % mod
			stayProb[i] = stay
			invFall[i] = modPow(f, mod-2)
			val[i] = stay * invFall[i] % mod
		}

		adj := make([][]int, n)
		edges := make([]edge, 0, n-1)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			u--
			v--
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
			edges = append(edges, edge{u, v})
		}

		F := make([]int64, n)
		Ssum := make([]int64, n)
		for i := 0; i < n; i++ {
			if len(adj[i]) == 0 {
				F[i] = 1
				continue
			}
			prod := int64(1)
			sum := int64(0)
			for _, nb := range adj[i] {
				prod = prod * fallProb[nb] % mod
				sum += val[nb]
				if sum >= mod {
					sum -= mod
				}
			}
			F[i] = prod
			Ssum[i] = sum
		}

		m := make([]int64, n)
		A := make([]int64, n)
		sumM := int64(0)
		sumSq := int64(0)
		for i := 0; i < n; i++ {
			Ai := stayProb[i] * F[i] % mod
			A[i] = Ai
			mi := Ai * Ssum[i] % mod
			m[i] = mi
			sumM += mi
			if sumM >= mod {
				sumM -= mod
			}
			sumSq += mi * mi % mod
			if sumSq >= mod {
				sumSq -= mod
			}
		}

		ans := pairSum(sumM, sumSq)

		for _, e := range edges {
			u, v := e.u, e.v
			actual := stayProb[u] * stayProb[v] % mod
			actual = actual * F[u] % mod * F[v] % mod
			actual = actual * invFall[u] % mod * invFall[v] % mod
			product := m[u] * m[v] % mod
			delta := (actual - product) % mod
			ans += delta
			if ans >= mod {
				ans -= mod
			} else if ans < 0 {
				ans += mod
			}
		}

		for b := 0; b < n; b++ {
			deg := len(adj[b])
			if deg < 2 {
				continue
			}
			sumA := int64(0)
			sumA2 := int64(0)
			sumB := int64(0)
			sumB2 := int64(0)
			sumMN := int64(0)
			sumMN2 := int64(0)
			valb := val[b]
			for _, u := range adj[b] {
				Au := A[u]
				sumA += Au
				if sumA >= mod {
					sumA -= mod
				}
				sumA2 += Au * Au % mod
				if sumA2 >= mod {
					sumA2 -= mod
				}
				term := Ssum[u] - valb
				if term < 0 {
					term += mod
				}
				Bu := stayProb[u] * F[u] % mod * term % mod
				sumB += Bu
				if sumB >= mod {
					sumB -= mod
				}
				sumB2 += Bu * Bu % mod
				if sumB2 >= mod {
					sumB2 -= mod
				}
				mu := m[u]
				sumMN += mu
				if sumMN >= mod {
					sumMN -= mod
				}
				sumMN2 += mu * mu % mod
				if sumMN2 >= mod {
					sumMN2 -= mod
				}
			}

			pairA := pairSum(sumA, sumA2)
			pairB := pairSum(sumB, sumB2)
			pairM := pairSum(sumMN, sumMN2)

			const1 := stayProb[b] * invFall[b] % mod * invFall[b] % mod
			contr1 := const1 * pairA % mod
			contr2 := invFall[b] * pairB % mod

			delta := (contr1 + contr2 - pairM) % mod
			ans += delta
			if ans >= mod {
				ans -= mod
			} else if ans < 0 {
				ans += mod
			}
		}

		if ans < 0 {
			ans += mod
		}
		fmt.Fprintln(out, ans%mod)
	}
}
