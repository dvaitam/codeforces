package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 1000000007

var (
	n, k    int
	w       [][]int
	dp      [][]int
	subtree []int
	tmp     []int
	answer  []int
	C       [][]int
)

func add(a, b int) int {
	a += b
	if a >= mod {
		a -= mod
	}
	return a
}

func sub(a, b int) int {
	a -= b
	if a < 0 {
		a += mod
	}
	return a
}

func dfs(a, parent int) {
	dp[a][0] = 2
	subtree[a] = 1
	for _, b := range w[a] {
		if b == parent {
			continue
		}
		dfs(b, a)
		// subtract dp[b] from answer
		limb := subtree[b]
		if limb > k {
			limb = k
		}
		for i := 0; i <= limb; i++ {
			answer[i] = sub(answer[i], dp[b][i])
		}
		// merge dp[a] and dp[b] into tmp
		limA := subtree[a]
		if limA > k+1 {
			limA = k + 1
		}
		for i := 0; i < limA; i++ {
			maxj := k - i
			if subtree[b] < maxj {
				maxj = subtree[b]
			}
			for j := 0; j <= maxj; j++ {
				tmp[i+j] = (tmp[i+j] + dp[a][i]*dp[b][j]) % mod
			}
		}
		subtree[a] += subtree[b]
		// copy back
		limA = subtree[a]
		if limA > k {
			limA = k
		}
		for i := 0; i <= limA; i++ {
			dp[a][i] = tmp[i]
			tmp[i] = 0
		}
	}
	lim := subtree[a]
	if lim > k {
		lim = k
	}
	for i := 0; i <= lim; i++ {
		answer[i] = add(answer[i], dp[a][i])
	}
	// shift dp[a] for non-root
	if a != 1 {
		for i := k; i >= 1; i-- {
			val := dp[a][i] + dp[a][i-1]
			if i == 1 {
				val--
			}
			if val < 0 {
				val += mod
			}
			if val >= mod {
				val -= mod
			}
			dp[a][i] = val
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	fmt.Fscan(in, &n, &k)
	w = make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		w[a] = append(w[a], b)
		w[b] = append(w[b], a)
	}
	dp = make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, k+1)
	}
	subtree = make([]int, n+1)
	tmp = make([]int, k+1)
	answer = make([]int, k+1)
	// compute binomial C
	C = make([][]int, k+1)
	for i := 0; i <= k; i++ {
		C[i] = make([]int, k+1)
		C[i][0], C[i][i] = 1, 1
		for j := 1; j < i; j++ {
			C[i][j] = add(C[i-1][j-1], C[i-1][j])
		}
	}
	// run dfs from root = 1
	dfs(1, 0)
	// compute ple for expansion
	ple := make([]int, k+1)
	ple[0] = 1
	total := 0
	for rozne := 1; rozne <= k; rozne++ {
		ple2 := make([]int, k+1)
		for already := 0; already <= k; already++ {
			pa := ple[already]
			if pa == 0 {
				continue
			}
			for nowe := 1; already+nowe <= k; nowe++ {
				ple2[already+nowe] = (ple2[already+nowe] + pa*C[already+nowe][nowe]) % mod
			}
		}
		ple = ple2
		total = (total + ple[k]*answer[rozne]) % mod
	}
	fmt.Fprint(out, total)
}
