package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	g    [][]int
	lval []int64
	rval []int64
	res  int
)

func dfs(v int) int64 {
	sum := int64(0)
	for _, to := range g[v] {
		sum += dfs(to)
	}
	if sum < lval[v] {
		res++
		return rval[v]
	}
	if sum > rval[v] {
		return rval[v]
	}
	return sum
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		g = make([][]int, n)
		for i := 1; i < n; i++ {
			var p int
			fmt.Fscan(reader, &p)
			p--
			g[p] = append(g[p], i)
		}
		lval = make([]int64, n)
		rval = make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &lval[i], &rval[i])
		}
		res = 0
		dfs(0)
		fmt.Fprintln(writer, res)
	}
}
