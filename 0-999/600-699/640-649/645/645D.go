package main

import (
	"bufio"
	"fmt"
	"os"
)

// check returns true if the first k edges uniquely define a total order on n nodes.
func check(n int, edges [][2]int, k int) bool {
	g := make([][]int, n+1)
	indeg := make([]int, n+1)
	for i := 0; i < k; i++ {
		u := edges[i][0]
		v := edges[i][1]
		g[u] = append(g[u], v)
		indeg[v]++
	}
	q := make([]int, 0, n)
	for i := 1; i <= n; i++ {
		if indeg[i] == 0 {
			q = append(q, i)
		}
	}
	if len(q) != 1 {
		return false
	}
	head := 0
	for head < len(q) {
		if head != len(q)-1 {
			return false
		}
		u := q[head]
		head++
		for _, v := range g[u] {
			indeg[v]--
			if indeg[v] == 0 {
				q = append(q, v)
			}
		}
	}
	return len(q) == n
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	edges := make([][2]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &edges[i][0], &edges[i][1])
	}

	if !check(n, edges, m) {
		fmt.Println(-1)
		return
	}

	lo, hi := 1, m
	ans := m
	for lo <= hi {
		mid := (lo + hi) / 2
		if check(n, edges, mid) {
			ans = mid
			hi = mid - 1
		} else {
			lo = mid + 1
		}
	}
	fmt.Println(ans)
}
