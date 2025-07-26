package main

import (
	"bufio"
	"fmt"
	"os"
)

// Solution for problem E from the 1697 directory.
// The task is to count the number of valid colorings of points on a plane
// so that points sharing a color form special clusters.
// A cluster of size >1 must have all pairwise distances equal and
// strictly smaller than distances from its points to any outside point.

const MOD int64 = 998244353

func abs64(v int64) int64 {
	if v < 0 {
		return -v
	}
	return v
}

func modPow(a, b int64) int64 {
	res := int64(1)
	a %= MOD
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
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

	x := make([]int64, n)
	y := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &x[i], &y[i])
	}

	dist := make([][]int64, n)
	dmin := make([]int64, n)
	for i := 0; i < n; i++ {
		dist[i] = make([]int64, n)
		dmin[i] = 1 << 62
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			d := abs64(x[i]-x[j]) + abs64(y[i]-y[j])
			dist[i][j] = d
			if d < dmin[i] {
				dmin[i] = d
			}
		}
	}

	neighbors := make([][]bool, n)
	neighCnt := make([]int, n)
	for i := 0; i < n; i++ {
		neighbors[i] = make([]bool, n)
		for j := 0; j < n; j++ {
			if i != j && dist[i][j] == dmin[i] {
				neighbors[i][j] = true
				neighCnt[i]++
			}
		}
	}

	visited := make([]bool, n)
	groupSizes := []int{}
	singles := 0

	for i := 0; i < n; i++ {
		if visited[i] {
			continue
		}
		// candidate group formed by i and its closest points
		cand := []int{i}
		for j := 0; j < n; j++ {
			if neighbors[i][j] {
				cand = append(cand, j)
			}
		}
		boolSet := make([]bool, n)
		for _, v := range cand {
			boolSet[v] = true
		}
		base := dmin[i]
		valid := true
		for _, idx := range cand {
			if dmin[idx] != base || neighCnt[idx]+1 != len(cand) {
				valid = false
				break
			}
			for p := 0; p < n; p++ {
				want := boolSet[p]
				have := neighbors[idx][p]
				if p == idx {
					have = true
				}
				if want != have {
					valid = false
					break
				}
			}
			if !valid {
				break
			}
		}
		if valid && len(cand) > 1 {
			for _, idx := range cand {
				visited[idx] = true
			}
			groupSizes = append(groupSizes, len(cand))
		} else {
			visited[i] = true
			singles++
		}
	}

	dp := make([]int64, n+1)
	dp[singles] = 1
	for _, s := range groupSizes {
		next := make([]int64, n+1)
		for k := 0; k <= n; k++ {
			if dp[k] == 0 {
				continue
			}
			if k+1 <= n {
				next[k+1] = (next[k+1] + dp[k]) % MOD
			}
			if k+s <= n {
				next[k+s] = (next[k+s] + dp[k]) % MOD
			}
		}
		dp = next
	}

	fact := make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invFact := make([]int64, n+1)
	invFact[n] = modPow(fact[n], MOD-2)
	for i := n; i >= 1; i-- {
		invFact[i-1] = invFact[i] * int64(i) % MOD
	}
	perm := func(n, k int) int64 {
		if k > n {
			return 0
		}
		return fact[n] * invFact[n-k] % MOD
	}

	ans := int64(0)
	for k := 0; k <= n; k++ {
		if dp[k] == 0 {
			continue
		}
		ans = (ans + dp[k]*perm(n, k)) % MOD
	}
	fmt.Fprintln(out, ans)
}
