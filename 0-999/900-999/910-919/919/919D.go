package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	var s string
	fmt.Fscan(reader, &s)

	edges := make([][]int, n)
	indeg := make([]int, n)
	for i := 0; i < m; i++ {
		var x, y int
		fmt.Fscan(reader, &x, &y)
		x--
		y--
		edges[x] = append(edges[x], y)
		indeg[y]++
	}

	dp := make([][26]int, n)
	for i := 0; i < n; i++ {
		idx := int(s[i] - 'a')
		dp[i][idx] = 1
	}

	q := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if indeg[i] == 0 {
			q = append(q, i)
		}
	}

	processed := 0
	for head := 0; head < len(q); head++ {
		v := q[head]
		processed++
		for _, to := range edges[v] {
			for c := 0; c < 26; c++ {
				val := dp[v][c]
				if c == int(s[to]-'a') {
					val++
				}
				if val > dp[to][c] {
					dp[to][c] = val
				}
			}
			indeg[to]--
			if indeg[to] == 0 {
				q = append(q, to)
			}
		}
	}

	if processed < n {
		fmt.Fprintln(writer, -1)
		return
	}

	ans := 0
	for i := 0; i < n; i++ {
		for c := 0; c < 26; c++ {
			if dp[i][c] > ans {
				ans = dp[i][c]
			}
		}
	}
	fmt.Fprintln(writer, ans)
}
