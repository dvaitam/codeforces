package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

func powmod(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
	}
	return res
}

func make2(n, m int) [][]int64 {
	r := make([][]int64, n)
	for i := 0; i < n; i++ {
		r[i] = make([]int64, m)
	}
	return r
}

func singleDist(L int) []int64 {
	if L == 0 {
		return []int64{1}
	}
	dpA := make([]int64, L+1)
	dpO := make([]int64, L+1)
	dpA[1] = 1
	dpO[0] = 25
	for pos := 1; pos < L; pos++ {
		newA := make([]int64, L+1)
		newO := make([]int64, L+1)
		for k := 0; k <= pos; k++ {
			valA := dpA[k]
			if valA != 0 {
				newO[k] = (newO[k] + valA*25) % MOD
			}
			valO := dpO[k]
			if valO != 0 {
				newA[k+1] = (newA[k+1] + valO) % MOD
				newO[k] = (newO[k] + valO*24) % MOD
			}
		}
		dpA, dpO = newA, newO
	}
	res := make([]int64, L+1)
	for k := 0; k <= L; k++ {
		res[k] = (dpA[k] + dpO[k]) % MOD
	}
	return res
}

func pairDist(L int) [][]int64 {
	dpA := make2(L+1, L+1)
	dpB := make2(L+1, L+1)
	dpO := make2(L+1, L+1)
	if L == 0 {
		dpO[0][0] = 1
		return dpO
	}
	dpA[1][0] = 1
	dpB[0][1] = 1
	dpO[0][0] = 24
	for pos := 1; pos < L; pos++ {
		newA := make2(L+1, L+1)
		newB := make2(L+1, L+1)
		newO := make2(L+1, L+1)
		for a := 0; a <= pos; a++ {
			for b := 0; b <= pos-a; b++ {
				if v := dpA[a][b]; v != 0 {
					if b+1 <= L {
						newB[a][b+1] = (newB[a][b+1] + v) % MOD
					}
					newO[a][b] = (newO[a][b] + v*24) % MOD
				}
				if v := dpB[a][b]; v != 0 {
					if a+1 <= L {
						newA[a+1][b] = (newA[a+1][b] + v) % MOD
					}
					newO[a][b] = (newO[a][b] + v*24) % MOD
				}
				if v := dpO[a][b]; v != 0 {
					if a+1 <= L {
						newA[a+1][b] = (newA[a+1][b] + v) % MOD
					}
					if b+1 <= L {
						newB[a][b+1] = (newB[a][b+1] + v) % MOD
					}
					newO[a][b] = (newO[a][b] + v*23) % MOD
				}
			}
		}
		dpA, dpB, dpO = newA, newB, newO
	}
	res := make2(L+1, L+1)
	for a := 0; a <= L; a++ {
		for b := 0; b <= L; b++ {
			res[a][b] = (dpA[a][b] + dpB[a][b] + dpO[a][b]) % MOD
		}
	}
	return res
}

func solve(n int, c []int) int64 {
	L0 := (n + 1) / 2
	L1 := n / 2
	d0 := singleDist(L0)
	d1 := singleDist(L1)
	single := make([]int64, n+1)
	for i := 0; i <= L0; i++ {
		for j := 0; j <= L1; j++ {
			single[i+j] = (single[i+j] + d0[i]*d1[j]) % MOD
		}
	}
	p0 := pairDist(L0)
	p1 := pairDist(L1)
	pair := make2(n+1, n+1)
	for a0 := 0; a0 <= L0; a0++ {
		for b0 := 0; b0 <= L0-a0; b0++ {
			v0 := p0[a0][b0]
			if v0 == 0 {
				continue
			}
			for a1 := 0; a1 <= L1; a1++ {
				for b1 := 0; b1 <= L1-a1; b1++ {
					pair[a0+a1][b0+b1] = (pair[a0+a1][b0+b1] + v0*p1[a1][b1]) % MOD
				}
			}
		}
	}
	singleSuf := make([]int64, n+2)
	for i := n; i >= 0; i-- {
		singleSuf[i] = (single[i] + singleSuf[i+1]) % MOD
	}
	pairSuf := make2(n+2, n+2)
	for i := n; i >= 0; i-- {
		for j := n; j >= 0; j-- {
			val := pair[i][j]
			pairSuf[i][j] = (val + pairSuf[i+1][j] + pairSuf[i][j+1] - pairSuf[i+1][j+1]) % MOD
			if pairSuf[i][j] < 0 {
				pairSuf[i][j] += MOD
			}
		}
	}
	total := int64(26*26%MOD) * powmod(25, int64(n-2)) % MOD
	for _, x := range c {
		total = (total - singleSuf[x+1]) % MOD
	}
	if total < 0 {
		total += MOD
	}
	for i := 0; i < 26; i++ {
		for j := i + 1; j < 26; j++ {
			total = (total + pairSuf[c[i]+1][c[j]+1]) % MOD
		}
	}
	if total < 0 {
		total += MOD
	}
	return total % MOD
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	c := make([]int, 26)
	for i := 0; i < 26; i++ {
		fmt.Fscan(in, &c[i])
	}
	fmt.Println(solve(n, c))
}
