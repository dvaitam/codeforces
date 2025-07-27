package main

import (
	"bufio"
	"fmt"
	"os"
)

var n, m, k int
var grid [][]int
var memo [][]int

func drop(r, c int) int {
	if r >= n {
		return c
	}
	if grid[r][c] == 2 && memo[r][c] != -1 {
		return memo[r][c]
	}
	var res int
	switch grid[r][c] {
	case 1:
		grid[r][c] = 2
		res = drop(r, c+1)
	case 2:
		res = drop(r+1, c)
	case 3:
		grid[r][c] = 2
		res = drop(r, c-1)
	}
	memo[r][c] = res
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Fscan(in, &n, &m, &k)
	grid = make([][]int, n)
	memo = make([][]int, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]int, m)
		memo[i] = make([]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(in, &grid[i][j])
			memo[i][j] = -1
		}
	}

	out := bufio.NewWriter(os.Stdout)
	for i := 0; i < k; i++ {
		var c int
		fmt.Fscan(in, &c)
		ans := drop(0, c-1)
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, ans+1)
	}
	out.WriteByte('\n')
	out.Flush()
}
