package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var m, n int
	if _, err := fmt.Fscan(reader, &m, &n); err != nil {
		return
	}

	profits := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &profits[i])
	}

	pos := make([][]int, m)
	for i := 0; i < m; i++ {
		pos[i] = make([]int, n)
		for j := 0; j < n; j++ {
			var r int
			fmt.Fscan(reader, &r)
			pos[i][j] = r
		}
	}

	order := make([]int, n)
	for i := 0; i < n; i++ {
		order[i] = i
	}
	sort.Slice(order, func(a, b int) bool {
		ra := pos[0][order[a]]
		rb := pos[0][order[b]]
		if ra == rb {
			return order[a] < order[b]
		}
		return ra < rb
	})

	dp := make([]int64, n)
	for _, idx := range order {
		dp[idx] = profits[idx]
	}

	for i := 0; i < n; i++ {
		a := order[i]
		for j := i + 1; j < n; j++ {
			b := order[j]
			ok := true
			for city := 1; city < m; city++ {
				if pos[city][a] >= pos[city][b] {
					ok = false
					break
				}
			}
			if ok {
				if dp[a]+profits[b] > dp[b] {
					dp[b] = dp[a] + profits[b]
				}
			}
		}
	}

	var ans int64
	for _, v := range dp {
		if v > ans {
			ans = v
		}
	}

	fmt.Fprintln(writer, ans)
}
