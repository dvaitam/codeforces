package main

import (
	"fmt"
)

const MAX_SIEVE = 5000000

var prefPhi []uint32
var primes []int
var memoPhi map[int]uint32

func initSieve() {
	prefPhi = make([]uint32, MAX_SIEVE+1)
	phi := make([]uint32, MAX_SIEVE+1)
	for i := 1; i <= MAX_SIEVE; i++ {
		phi[i] = uint32(i)
	}

	for i := 2; i <= MAX_SIEVE; i++ {
		if phi[i] == uint32(i) {
			for j := i; j <= MAX_SIEVE; j += i {
				phi[j] = phi[j] / uint32(i) * uint32(i-1)
			}
		}
	}
	for i := 1; i <= MAX_SIEVE; i++ {
		prefPhi[i] = prefPhi[i-1] + phi[i]
	}

	const MAX_PRIME = 13000000
	isP := make([]bool, MAX_PRIME+1)
	for i := 2; i <= MAX_PRIME; i++ {
		isP[i] = true
	}
	for i := 2; i*i <= MAX_PRIME; i++ {
		if isP[i] {
			for j := i * i; j <= MAX_PRIME; j += i {
				isP[j] = false
			}
		}
	}
	for i := 2; i <= MAX_PRIME; i++ {
		if isP[i] {
			primes = append(primes, i)
		}
	}
	memoPhi = make(map[int]uint32)
}

func PhiSum(x int) uint32 {
	if x <= MAX_SIEVE {
		return prefPhi[x]
	}
	if val, ok := memoPhi[x]; ok {
		return val
	}
	var res uint32 = uint32(x) * uint32(x+1) / 2
	for l := 2; l <= x; {
		q := x / l
		r := x / q
		res -= uint32(r-l+1) * PhiSum(q)
		l = r + 1
	}
	memoPhi[x] = res
	return res
}

type state struct{ X, D int }

var memoSD = make(map[state]uint32)

func S_dense(X int, D int) uint32 {
	if X == 0 {
		return 0
	}
	if D == 1 {
		return PhiSum(X)
	}
	st := state{X, D}
	if val, ok := memoSD[st]; ok {
		return val
	}

	var p int
	for _, pr := range []int{2, 3, 5, 7, 11} {
		if D%pr == 0 {
			p = pr
			break
		}
	}
	res := uint32(p-1)*S_dense(X, D/p) + S_dense(X/p, D)
	memoSD[st] = res
	return res
}

var n int
var sparse_sum [6]uint32
var D_c = [6]int{1, 1, 2, 6, 210, 2310}

func getMaxRatio(k int, phi int, p_idx int) float64 {
	rem := n / k
	ratio := float64(k) / float64(phi)
	for i := p_idx; i < len(primes); i++ {
		p := primes[i]
		if p > rem {
			break
		}
		ratio *= float64(p) / float64(p-1)
		rem /= p
	}
	return ratio
}

func dfs(p_idx int, k int, phi int) {
	for i := p_idx; i < len(primes); i++ {
		p := primes[i]
		if k*p > n {
			break
		}

		possible := false
		for c := 2; c <= 5; c++ {
			if (k*p)%D_c[c] != 0 {
				nr := float64(k*p) / float64(phi*(p-1))
				if nr >= float64(c) {
					possible = true
					break
				}
				mr := getMaxRatio(k*p, phi*(p-1), i+1)
				if mr >= float64(c) {
					possible = true
					break
				}
			}
		}
		if !possible {
			break
		}

		nk := k * p
		nphi := phi * (p - 1)
		for nk <= n {
			nr := float64(nk) / float64(nphi)
			for c := 2; c <= 5; c++ {
				if nk%D_c[c] != 0 && nr >= float64(c) {
					sparse_sum[c] += uint32(nphi)
				}
			}
			dfs(i+1, nk, nphi)
			nk *= p
			nphi *= p
		}
	}
}

func main() {
	initSieve()
	if _, err := fmt.Scan(&n); err != nil {
		return
	}

	dfs(0, 1, 1)

	var totalWPhi uint32
	totalWPhi += S_dense(n, 1)
	totalWPhi += S_dense(n/2, 2) + sparse_sum[2]
	totalWPhi += S_dense(n/6, 6) + sparse_sum[3]
	totalWPhi += S_dense(n/210, 210) + sparse_sum[4]
	totalWPhi += S_dense(n/2310, 2310) + sparse_sum[5]

	var sumK uint32 = uint32(n % (1 << 32))
	var sumK1 uint32 = uint32((n + 1) % (1 << 32))
	var totalK uint32
	if sumK%2 == 0 {
		totalK = (sumK / 2) * sumK1
	} else {
		totalK = sumK * (sumK1 / 2)
	}

	ans := totalK - totalWPhi
	fmt.Println(ans)
}
