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
	for t := 0; t < T; t++ {
		solve()
	}
}

func solve() {
	var n int
	fmt.Fscan(reader, &n)
	p := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &p[i])
		p[i]--
	}
	tree := make([][]int, n)
	root := 0
	for i := 0; i < n; i++ {
		if p[i] == i {
			root = i
		} else {
			tree[p[i]] = append(tree[p[i]], i)
		}
	}
	ans := make([][]int, 0, n)
	ans = append(ans, []int{})
	var dfs func(v, pth int)
	dfs = func(v, pth int) {
		ans[pth] = append(ans[pth], v)
		children := tree[v]
		if len(children) == 0 {
			return
		}
		// continue path with first child
		dfs(children[0], pth)
		// spawn new paths for other children
		for i := 1; i < len(children); i++ {
			newpth := len(ans)
			ans = append(ans, []int{})
			dfs(children[i], newpth)
		}
	}
	dfs(root, 0)
	m := len(ans)
	fmt.Fprintln(writer, m)
	for i := 0; i < m; i++ {
		k := len(ans[i])
		fmt.Fprintln(writer, k)
		for j := 0; j < k; j++ {
			fmt.Fprint(writer, ans[i][j]+1)
			if j+1 < k {
				fmt.Fprint(writer, " ")
			}
		}
		fmt.Fprintln(writer)
	}
	// blank line after test case
	fmt.Fprintln(writer)
}
