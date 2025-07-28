package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	g := make([][]int, n)
	for i := 0; i < m; i++ {
		var u, v, w int
		fmt.Fscan(reader, &u, &v, &w)
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}

	size := int(math.Sqrt(float64(n)))
	if size <= 0 {
		size = 1
	}

	visited := make([]bool, n)
	type state struct{ v, idx int }
	dfsStack := []state{{0, 0}}
	visited[0] = true
	path := []int{0}
	groups := [][]int{}

	for len(dfsStack) > 0 {
		st := &dfsStack[len(dfsStack)-1]
		v := st.v
		if st.idx < len(g[v]) {
			to := g[v][st.idx]
			st.idx++
			if !visited[to] {
				visited[to] = true
				dfsStack = append(dfsStack, state{to, 0})
				path = append(path, to)
				if len(path) >= size {
					grp := append([]int(nil), path[len(path)-size:]...)
					groups = append(groups, grp)
					path = path[:len(path)-size]
				}
			}
		} else {
			dfsStack = dfsStack[:len(dfsStack)-1]
		}
	}

	if len(path) > 0 {
		groups = append(groups, append([]int(nil), path...))
	}

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	fmt.Fprintln(writer, len(groups))
	for _, grp := range groups {
		fmt.Fprint(writer, len(grp))
		for _, v := range grp {
			fmt.Fprint(writer, " ", v)
		}
		fmt.Fprintln(writer)
	}
}
