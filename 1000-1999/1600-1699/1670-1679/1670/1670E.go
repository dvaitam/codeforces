package main

import (
	"bufio"
	"fmt"
	"os"
)

type edge struct {
	to  int
	idx int
}

func solve(reader *bufio.Reader, writer *bufio.Writer) {
	var x int
	fmt.Fscan(reader, &x)
	n := 1 << x
	g := make([][]edge, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		u--
		v--
		g[u] = append(g[u], edge{v, i})
		g[v] = append(g[v], edge{u, i})
	}
	p1 := make([]int, n)
	p2 := make([]int, n-1)
	p1[0] = n
	var dfs func(u, parent int, left bool)
	dfs = func(u, parent int, left bool) {
		for _, e := range g[u] {
			v := e.to
			if v == parent {
				continue
			}
			if left {
				p1[v] = v ^ n
				p2[e.idx] = v
			} else {
				p1[v] = v
				p2[e.idx] = v ^ n
			}
			dfs(v, u, !left)
		}
	}
	dfs(0, -1, false)
	fmt.Fprintln(writer, 1)
	for i := 0; i < n; i++ {
		fmt.Fprint(writer, p1[i], " ")
	}
	fmt.Fprintln(writer)
	for i := 0; i < n-1; i++ {
		fmt.Fprint(writer, p2[i], " ")
	}
	fmt.Fprintln(writer)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var t int
	fmt.Fscan(reader, &t)
	for t > 0 {
		solve(reader, writer)
		t--
	}
}
