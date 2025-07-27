package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const mod = 998244353

func powMod(a, e int) int {
	res := 1
	base := a
	for e > 0 {
		if e&1 == 1 {
			res = int(int64(res) * int64(base) % mod)
		}
		base = int(int64(base) * int64(base) % mod)
		e >>= 1
	}
	return res
}

type state struct {
	idx int
	val uint64
}

func enumerate(vects []uint64, m int) []int {
	res := make([]int, m+1)
	stack := []state{{0, 0}}
	k := len(vects)
	for len(stack) > 0 {
		s := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if s.idx == k {
			res[bits.OnesCount64(s.val)]++
		} else {
			stack = append(stack, state{s.idx + 1, s.val})
			stack = append(stack, state{s.idx + 1, s.val ^ vects[s.idx]})
		}
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	arr := make([]uint64, n)
	for i := 0; i < n; i++ {
		var x uint64
		fmt.Fscan(reader, &x)
		arr[i] = x
	}

	basis := make([]uint64, m)
	for _, v := range arr {
		x := v
		for i := m - 1; i >= 0; i-- {
			if ((x >> uint(i)) & 1) == 0 {
				continue
			}
			if basis[i] != 0 {
				x ^= basis[i]
			} else {
				basis[i] = x
				break
			}
		}
	}

	// reduce to row echelon form
	for i := 0; i < m; i++ {
		if basis[i] != 0 {
			for j := i + 1; j < m; j++ {
				if (basis[j]>>uint(i))&1 == 1 {
					basis[j] ^= basis[i]
				}
			}
		}
	}

	var vectors []uint64
	var pivots []int
	for i := 0; i < m; i++ {
		if basis[i] != 0 {
			vectors = append(vectors, basis[i])
			pivots = append(pivots, i)
		}
	}
	r := len(vectors)
	d := m - r

	pow2 := make([]int, n+1)
	pow2[0] = 1
	for i := 1; i <= n; i++ {
		pow2[i] = int(int64(pow2[i-1]) * 2 % mod)
	}

	if r <= d {
		cnt := enumerate(vectors, m)
		for i := 0; i <= m; i++ {
			ans := int(int64(cnt[i]) * int64(pow2[n-r]) % mod)
			if i > 0 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, ans)
		}
		fmt.Fprintln(writer)
		return
	}

	// build nullspace basis
	used := make([]bool, m)
	for _, p := range pivots {
		used[p] = true
	}
	var dual []uint64
	for f := 0; f < m; f++ {
		if used[f] {
			continue
		}
		x := uint64(1) << uint(f)
		for idx, p := range pivots {
			if (vectors[idx]>>uint(f))&1 == 1 {
				x |= uint64(1) << uint(p)
			}
		}
		dual = append(dual, x)
	}

	cntDual := enumerate(dual, m)

	// precompute binomial coefficients
	comb := make([][]int, m+1)
	for i := 0; i <= m; i++ {
		comb[i] = make([]int, i+1)
		comb[i][0] = 1
		comb[i][i] = 1
		for j := 1; j < i; j++ {
			comb[i][j] = (comb[i-1][j-1] + comb[i-1][j]) % mod
		}
	}

	inv := powMod(pow2[d], mod-2)
	ans := make([]int, m+1)
	for i := 0; i <= m; i++ {
		sum := 0
		for j := 0; j <= m; j++ {
			kcoef := 0
			maxT := i
			if j < maxT {
				maxT = j
			}
			for t := 0; t <= maxT; t++ {
				if i-t > m-j {
					continue
				}
				term := int(int64(comb[j][t]) * int64(comb[m-j][i-t]) % mod)
				if t%2 == 1 {
					term = (mod - term) % mod
				}
				kcoef += term
				if kcoef >= mod {
					kcoef -= mod
				}
			}
			sum = (sum + int(int64(cntDual[j])*int64(kcoef)%mod)) % mod
		}
		ans[i] = int(int64(sum) * int64(inv) % mod)
		ans[i] = int(int64(ans[i]) * int64(pow2[n-r]) % mod)
	}

	for i := 0; i <= m; i++ {
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, ans[i])
	}
	fmt.Fprintln(writer)
}
