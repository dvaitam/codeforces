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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	a := make([]int, n+1)
	pos := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
		pos[a[i]] = i
	}

	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		adj[u] = append(adj[u], v)
	}

	minVal := make([]int, n+1)
	var dfsMin func(int)
	dfsMin = func(u int) {
		minVal[u] = a[u]
		for _, v := range adj[u] {
			dfsMin(v)
			if minVal[v] < minVal[u] {
				minVal[u] = minVal[v]
			}
		}
	}
	dfsMin(1)

	for u := 1; u <= n; u++ {
		sort.Slice(adj[u], func(i, j int) bool {
			return minVal[adj[u][i]] < minVal[adj[u][j]]
		})
	}

	aOrig := make([]int, n+1)
	origPos := make([]int, n+1)
	timer := 1
	var buildOrig func(int)
	buildOrig = func(u int) {
		aOrig[u] = timer
		origPos[timer] = u
		timer++
		for _, v := range adj[u] {
			buildOrig(v)
		}
	}
	buildOrig(1)

	bit := make([]int, n+2)
	add := func(idx, val int) {
		for ; idx <= n; idx += idx & -idx {
			bit[idx] += val
		}
	}
	query := func(idx int) int {
		sum := 0
		for ; idx > 0; idx -= idx & -idx {
			sum += bit[idx]
		}
		return sum
	}

	var D int64 = 0
	var countInversions func(int)
	countInversions = func(u int) {
		greater := query(n) - query(a[u])
		D += int64(greater)
		add(a[u], 1)
		for _, v := range adj[u] {
			countInversions(v)
		}
		add(a[u], -1)
	}
	countInversions(1)

	maxX := 0
	for u := 1; u <= n; u++ {
		if a[u] < aOrig[u] {
			if a[u] > maxX {
				maxX = a[u]
			}
		}
	}

	for y := 1; y < maxX; y++ {
		u := pos[y]
		for _, v := range adj[u] {
			if a[v] > y {
				fmt.Fprintln(writer, "NO")
				return
			}
		}
	}

	for u := 1; u <= n; u++ {
		actualMin := a[u]
		var checkMin func(int)
		checkMin = func(curr int) {
			if a[curr] < actualMin {
				actualMin = a[curr]
			}
			for _, child := range adj[curr] {
				checkMin(child)
			}
		}
		checkMin(u)
		if actualMin != minVal[u] {
			fmt.Fprintln(writer, "NO")
			return
		}
	}

	fmt.Fprintln(writer, "YES")
	fmt.Fprintln(writer, D)
	for i := 1; i <= n; i++ {
		fmt.Fprint(writer, aOrig[i])
		if i < n {
			fmt.Fprint(writer, " ")
		}
	}
	fmt.Fprintln(writer)
}
