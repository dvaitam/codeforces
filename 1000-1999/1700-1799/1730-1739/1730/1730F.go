package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func minInversions(p []int, k int) int {
	n := len(p)
	arr := make([]int, n+1)
	for i, v := range p {
		arr[v] = i + 1
	}
	w := k + 1
	future := make([][]int, n+1)
	for i := range future {
		future[i] = make([]int, w)
	}
	bit := make([]int, n+2)
	update := func(i, val int) {
		for i <= n {
			bit[i] += val
			i += i & -i
		}
	}
	query := func(i int) int {
		s := 0
		for i > 0 {
			s += bit[i]
			i -= i & -i
		}
		return s
	}
	for i := n; i >= 1; i-- {
		limit := i - k
		if limit < 1 {
			limit = 1
		}
		for j := i; j >= limit; j-- {
			future[i][i-j] = query(arr[j] - 1)
		}
		update(arr[i], 1)
	}

	const INF = int(1e18)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, 1<<uint(w))
		for j := range dp[i] {
			dp[i][j] = -1
		}
	}

	var dfs func(int, int) int
	dfs = func(i, mask int) int {
		if i == n && mask == 0 {
			return 0
		}
		if dp[i][mask] != -1 {
			return dp[i][mask]
		}
		res := INF
		if mask != 0 {
			for j := 0; j < w; j++ {
				if mask>>uint(j)&1 == 1 {
					idx := i - j
					if idx <= 0 {
						continue
					}
					val := arr[idx]
					newMask := mask & ^(1 << uint(j))
					small := future[i][i-idx]
					for q := 0; q < w; q++ {
						if newMask>>uint(q)&1 == 1 {
							idx2 := i - q
							if idx2 > 0 && arr[idx2] < val {
								small++
							}
						}
					}
					cand := dfs(i, newMask) + small
					if cand < res {
						res = cand
					}
				}
			}
		}
		if i < n && bits.OnesCount(uint(mask)) < w && (mask>>uint(w-1)&1) == 0 {
			newMask := (mask << 1) | 1
			cand := dfs(i+1, newMask)
			if cand < res {
				res = cand
			}
		}
		dp[i][mask] = res
		return res
	}

	return dfs(0, 0)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	fmt.Fscan(reader, &n, &k)
	p := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &p[i])
	}
	ans := minInversions(p, k)
	fmt.Fprintln(writer, ans)
}
