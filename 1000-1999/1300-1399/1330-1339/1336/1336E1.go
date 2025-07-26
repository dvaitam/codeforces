package main

import (
	"bufio"
	"fmt"
	"math/bits"
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

// buildBasis returns basis vectors in RREF order and pivot indices
func buildBasis(arr []uint64, m int) ([]uint64, []int) {
	basis := make([]uint64, m)
	for _, x := range arr {
		v := x
		for i := m - 1; i >= 0; i-- {
			if (v>>i)&1 == 0 {
				continue
			}
			if basis[i] != 0 {
				v ^= basis[i]
			} else {
				basis[i] = v
				// remove bit i from higher pivot vectors
				for j := i + 1; j < m; j++ {
					if basis[j] != 0 && ((basis[j]>>i)&1) == 1 {
						basis[j] ^= v
					}
				}
				break
			}
		}
	}
	// eliminate upwards
	for i := m - 1; i >= 0; i-- {
		if basis[i] == 0 {
			continue
		}
		for j := i - 1; j >= 0; j-- {
			if basis[j] != 0 && ((basis[j]>>i)&1) == 1 {
				basis[j] ^= basis[i]
			}
		}
	}
	pivots := make([]int, 0)
	for i := 0; i < m; i++ {
		if basis[i] != 0 {
			pivots = append(pivots, i)
		}
	}
	vecs := make([]uint64, len(pivots))
	for idx, p := range pivots {
		vecs[idx] = basis[p]
	}
	return vecs, pivots
}

func dualBasis(basis []uint64, pivots []int, m int) []uint64 {
	pivotSet := make(map[int]bool)
	for _, p := range pivots {
		pivotSet[p] = true
	}
	dual := make([]uint64, 0, m-len(pivots))
	for j := 0; j < m; j++ {
		if pivotSet[j] {
			continue
		}
		w := uint64(1) << j
		for idx, p := range pivots {
			if ((basis[idx] >> j) & 1) == 1 {
				w |= uint64(1) << p
			}
		}
		dual = append(dual, w)
	}
	return dual
}

func enumerate(vecs []uint64, idx int, cur uint64, res []int64) {
	if idx == len(vecs) {
		w := bits.OnesCount64(cur)
		res[w]++
		return
	}
	enumerate(vecs, idx+1, cur, res)
	enumerate(vecs, idx+1, cur^vecs[idx], res)
}

func macWilliams(dual []int64, m int, r int) []int64 {
	d := m - r
	comb := make([][]int64, m+1)
	for i := 0; i <= m; i++ {
		comb[i] = make([]int64, m+1)
		comb[i][0] = 1
		for j := 1; j <= i; j++ {
			comb[i][j] = (comb[i-1][j-1] + comb[i-1][j]) % MOD
		}
	}
	res := make([]int64, m+1)
	powInv := powmod(2, int64(MOD-1)-int64(d)) // inverse of 2^d
	for k := 0; k <= m; k++ {
		var sum int64
		for i := 0; i <= m; i++ {
			if dual[i] == 0 {
				continue
			}
			var kk int64
			for j := 0; j <= k && j <= i; j++ {
				sign := int64(1)
				if j%2 == 1 {
					sign = MOD - 1
				}
				kk = (kk + sign*comb[i][j]%MOD*comb[m-i][k-j]) % MOD
			}
			sum = (sum + dual[i]*kk) % MOD
		}
		res[k] = sum * powInv % MOD
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	var m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	arr := make([]uint64, n)
	for i := 0; i < n; i++ {
		var x uint64
		fmt.Fscan(reader, &x)
		arr[i] = x
	}

	vecs, pivots := buildBasis(arr, m)
	r := len(vecs)
	pow2 := powmod(2, int64(n-r))

	ans := make([]int64, m+1)
	if r <= m-r {
		enumerate(vecs, 0, 0, ans)
		for i := range ans {
			ans[i] = ans[i] * pow2 % MOD
		}
	} else {
		dualVecs := dualBasis(vecs, pivots, m)
		dualCnt := make([]int64, m+1)
		enumerate(dualVecs, 0, 0, dualCnt)
		ans = macWilliams(dualCnt, m, r)
		for i := range ans {
			ans[i] = ans[i] * pow2 % MOD
		}
	}

	for i := 0; i <= m; i++ {
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, ans[i])
	}
	fmt.Fprintln(writer)
}
