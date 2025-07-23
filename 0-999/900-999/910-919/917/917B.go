package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	to int
	ch int
}

var (
	n, m int
	g    [][]Edge
	dp   [101][101][26]int8
)

func win(a, b, c int) int8 {
	res := &dp[a][b][c]
	if *res != 0 {
		return *res
	}
	for _, e := range g[a] {
		if e.ch >= c {
			if win(b, e.to, e.ch) == -1 {
				*res = 1
				return 1
			}
		}
	}
	*res = -1
	return -1
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	g = make([][]Edge, n+1)
	for i := 0; i < m; i++ {
		var v, u int
		var s string
		fmt.Fscan(reader, &v, &u, &s)
		g[v] = append(g[v], Edge{u, int(s[0] - 'a')})
	}

	for i := 1; i <= n; i++ {
		for j := 1; j <= n; j++ {
			if win(i, j, 0) == 1 {
				fmt.Fprint(writer, "A")
			} else {
				fmt.Fprint(writer, "B")
			}
		}
		fmt.Fprintln(writer)
	}
}
