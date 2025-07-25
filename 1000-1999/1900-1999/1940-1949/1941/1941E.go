package main

import (
	"bufio"
	"fmt"
	"os"
)

func rowCost(row []int, d int) int64 {
	m := len(row)
	const INF int64 = 1 << 60
	cost := make([]int64, m)
	for i, v := range row {
		cost[i] = int64(v) + 1
	}
	dp := make([]int64, m)
	for i := range dp {
		dp[i] = INF
	}
	dp[0] = cost[0]
	type pair struct {
		idx int
		val int64
	}
	deque := make([]pair, 0)
	for j := 1; j < m; j++ {
		val := dp[j-1]
		for len(deque) > 0 && deque[len(deque)-1].val >= val {
			deque = deque[:len(deque)-1]
		}
		deque = append(deque, pair{j - 1, val})
		limit := j - (d + 1)
		for len(deque) > 0 && deque[0].idx < limit {
			deque = deque[1:]
		}
		if len(deque) > 0 {
			dp[j] = cost[j] + deque[0].val
		}
	}
	return dp[m-1]
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m, k, d int
		fmt.Fscan(reader, &n, &m, &k, &d)
		grid := make([][]int, n)
		for i := 0; i < n; i++ {
			grid[i] = make([]int, m)
			for j := 0; j < m; j++ {
				fmt.Fscan(reader, &grid[i][j])
			}
		}
		costs := make([]int64, n)
		for i := 0; i < n; i++ {
			costs[i] = rowCost(grid[i], d)
		}
		var sum int64
		for i := 0; i < k; i++ {
			sum += costs[i]
		}
		ans := sum
		for i := k; i < n; i++ {
			sum += costs[i] - costs[i-k]
			if sum < ans {
				ans = sum
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
