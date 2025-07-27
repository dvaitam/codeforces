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
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		var row string
		fmt.Fscan(reader, &row)
		grid[i] = []byte(row)
	}
	// read the last line of a[i] but not used in easy version
	for i := 0; i < m; i++ {
		var x int
		fmt.Fscan(reader, &x)
		_ = x
	}

	id := make([][]int, n)
	for i := range id {
		id[i] = make([]int, m)
		for j := range id[i] {
			id[i][j] = -1
		}
	}
	rows := make([]int, 0)
	cols := make([]int, 0)
	idx := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '#' {
				id[i][j] = idx
				rows = append(rows, i)
				cols = append(cols, j)
				idx++
			}
		}
	}
	tot := idx
	if tot == 0 {
		fmt.Fprintln(writer, 0)
		return
	}

	firstBelow := make([][]int, m)
	for j := 0; j < m; j++ {
		firstBelow[j] = make([]int, n)
		nearest := -1
		for i := n - 1; i >= 0; i-- {
			if id[i][j] != -1 {
				nearest = i
			}
			firstBelow[j][i] = nearest
		}
	}

	belowID := make([]int, tot)
	for j := 0; j < m; j++ {
		lastID := -1
		for i := n - 1; i >= 0; i-- {
			if id[i][j] != -1 {
				cur := id[i][j]
				belowID[cur] = lastID
				lastID = cur
			}
		}
	}

	adj := make([][]int, tot)
	rev := make([][]int, tot)
	addEdge := func(u, v int) {
		adj[u] = append(adj[u], v)
		rev[v] = append(rev[v], u)
	}

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			u := id[i][j]
			if u == -1 {
				continue
			}
			if i > 0 && id[i-1][j] != -1 {
				addEdge(u, id[i-1][j])
			}
			if b := belowID[u]; b != -1 {
				addEdge(u, b)
			}
			if j > 0 {
				if r := firstBelow[j-1][i]; r != -1 {
					addEdge(u, id[r][j-1])
				}
			}
			if j+1 < m {
				if r := firstBelow[j+1][i]; r != -1 {
					addEdge(u, id[r][j+1])
				}
			}
		}
	}

	order := make([]int, 0, tot)
	visited := make([]bool, tot)
	for i := 0; i < tot; i++ {
		if visited[i] {
			continue
		}
		stack := []int{i}
		idxStack := []int{0}
		visited[i] = true
		for len(stack) > 0 {
			v := stack[len(stack)-1]
			idxCur := idxStack[len(idxStack)-1]
			if idxCur < len(adj[v]) {
				to := adj[v][idxCur]
				idxStack[len(idxStack)-1]++
				if !visited[to] {
					visited[to] = true
					stack = append(stack, to)
					idxStack = append(idxStack, 0)
				}
			} else {
				stack = stack[:len(stack)-1]
				idxStack = idxStack[:len(idxStack)-1]
				order = append(order, v)
			}
		}
	}

	comp := make([]int, tot)
	for i := range comp {
		comp[i] = -1
	}
	cnum := 0
	for k := len(order) - 1; k >= 0; k-- {
		v := order[k]
		if comp[v] != -1 {
			continue
		}
		stack := []int{v}
		comp[v] = cnum
		for len(stack) > 0 {
			x := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			for _, to := range rev[x] {
				if comp[to] == -1 {
					comp[to] = cnum
					stack = append(stack, to)
				}
			}
		}
		cnum++
	}

	indeg := make([]int, cnum)
	for u := 0; u < tot; u++ {
		cu := comp[u]
		for _, v := range adj[u] {
			cv := comp[v]
			if cu != cv {
				indeg[cv]++
			}
		}
	}

	ans := 0
	for i := 0; i < cnum; i++ {
		if indeg[i] == 0 {
			ans++
		}
	}
	fmt.Fprintln(writer, ans)
}
