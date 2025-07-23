package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, T int
	if _, err := fmt.Fscan(in, &n, &T); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	valsMap := make(map[int]struct{})
	for _, v := range a {
		valsMap[v] = struct{}{}
	}
	vals := make([]int, 0, len(valsMap))
	for v := range valsMap {
		vals = append(vals, v)
	}
	sort.Ints(vals)
	m := len(vals)
	idx := make(map[int]int, m)
	for i, v := range vals {
		idx[v] = i
	}

	neg := -1 << 60
	g := make([][]int, m)
	for i := range g {
		g[i] = make([]int, m)
		for j := range g[i] {
			g[i][j] = neg
		}
	}

	for s := 0; s < m; s++ {
		dp := make([]int, m)
		for j := range dp {
			dp[j] = neg
		}
		dp[s] = 0
		for _, val := range a {
			id := idx[val]
			best := dp[0]
			for k := 1; k <= id; k++ {
				if dp[k] > best {
					best = dp[k]
				}
			}
			if best+1 > dp[id] {
				dp[id] = best + 1
			}
		}
		for j := 0; j < m; j++ {
			g[s][j] = dp[j]
		}
	}

	matMul := func(A, B [][]int) [][]int {
		C := make([][]int, m)
		for i := range C {
			C[i] = make([]int, m)
			for j := range C[i] {
				C[i][j] = neg
			}
		}
		for i := 0; i < m; i++ {
			for k := 0; k < m; k++ {
				if A[i][k] == neg {
					continue
				}
				for j := 0; j < m; j++ {
					val := A[i][k] + B[k][j]
					if val > C[i][j] {
						C[i][j] = val
					}
				}
			}
		}
		return C
	}

	vecMul := func(vec []int, M [][]int) []int {
		res := make([]int, m)
		for i := range res {
			res[i] = neg
		}
		for i := 0; i < m; i++ {
			if vec[i] == neg {
				continue
			}
			for j := 0; j < m; j++ {
				val := vec[i] + M[i][j]
				if val > res[j] {
					res[j] = val
				}
			}
		}
		return res
	}

	pow := g
	vec := make([]int, m)
	for T > 0 {
		if T&1 == 1 {
			vec = vecMul(vec, pow)
		}
		if T > 1 {
			pow = matMul(pow, pow)
		}
		T >>= 1
	}

	ans := vec[0]
	for i := 1; i < m; i++ {
		if vec[i] > ans {
			ans = vec[i]
		}
	}
	fmt.Fprintln(out, ans)
}
