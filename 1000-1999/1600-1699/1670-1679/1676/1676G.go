package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	reader = bufio.NewReader(os.Stdin)
	writer = bufio.NewWriter(os.Stdout)
)

func main() {
	defer writer.Flush()
	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		solve()
	}
}

func solve() {
	var n int
	fmt.Fscan(reader, &n)
	parent := make([]int, n+1)
	for i := 2; i <= n; i++ {
		fmt.Fscan(reader, &parent[i])
	}
	var s string
	fmt.Fscan(reader, &s)

	tree := make([][]int, n+1)
	for i := 2; i <= n; i++ {
		p := parent[i]
		tree[p] = append(tree[p], i)
	}

	colors := []byte(s)
	ans := 0
	var dfs func(int) (int, int)
	dfs = func(v int) (int, int) {
		white, black := 0, 0
		if colors[v-1] == 'W' {
			white = 1
		} else {
			black = 1
		}
		for _, u := range tree[v] {
			w, b := dfs(u)
			white += w
			black += b
		}
		if white == black {
			ans++
		}
		return white, black
	}

	dfs(1)
	fmt.Fprintln(writer, ans)
}
