package main

import (
	"bufio"
	"fmt"
	"os"
)

type pair struct {
	win int
	sum int64
}

var (
	n, m int
	diff []int64
	g    [][]int
)

func dfs(u, p int) ([]pair, int) {
	dp := make([]pair, 2)
	dp[1] = pair{0, diff[u]}
	size := 1
	for _, v := range g[u] {
		if v == p {
			continue
		}
		child, sz := dfs(v, u)
		newSize := size + sz
		if newSize > m {
			newSize = m
		}
		ndp := make([]pair, newSize+1)
		for i := 1; i <= size && i <= m; i++ {
			for j := 1; j <= sz && i+j-1 <= m; j++ {
				// merge without cutting edge
				w1 := dp[i].win + child[j].win
				s1 := dp[i].sum + child[j].sum
				if w1 > ndp[i+j-1].win || (w1 == ndp[i+j-1].win && s1 > ndp[i+j-1].sum) {
					ndp[i+j-1] = pair{w1, s1}
				}
				// cut edge between u and v
				w2 := dp[i].win + child[j].win
				if child[j].sum > 0 {
					w2++
				}
				s2 := dp[i].sum
				if w2 > ndp[i+j].win || (w2 == ndp[i+j].win && s2 > ndp[i+j].sum) {
					ndp[i+j] = pair{w2, s2}
				}
			}
		}
		size = newSize
		dp = ndp
	}
	return dp, size
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		fmt.Fscan(reader, &n, &m)
		bees := make([]int64, n+1)
		wasps := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &bees[i])
		}
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &wasps[i])
		}
		g = make([][]int, n+1)
		for i := 0; i < n-1; i++ {
			var x, y int
			fmt.Fscan(reader, &x, &y)
			g[x] = append(g[x], y)
			g[y] = append(g[y], x)
		}
		diff = make([]int64, n+1)
		for i := 1; i <= n; i++ {
			diff[i] = wasps[i] - bees[i]
		}
		dp, _ := dfs(1, 0)
		ans := dp[m].win
		if dp[m].sum > 0 {
			ans++
		}
		fmt.Fprintln(writer, ans)
	}
}
