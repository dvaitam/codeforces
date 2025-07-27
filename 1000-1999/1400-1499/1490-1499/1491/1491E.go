package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	to int
	id int
}

var (
	fib                 []int
	fibIndex            map[int]int
	g                   [][]Edge
	removed             []bool
	sizeArr             []int
	cutU, cutV, cutEdge int
)

func buildFib(limit int) {
	fib = []int{1, 1}
	fibIndex = map[int]int{1: 1}
	for fib[len(fib)-1] < limit {
		n := len(fib)
		val := fib[n-1] + fib[n-2]
		fib = append(fib, val)
		fibIndex[val] = n
	}
}

func find(u, p, edgeID, k int) bool {
	sizeArr[u] = 1
	for _, e := range g[u] {
		if e.to == p || removed[e.id] {
			continue
		}
		if find(e.to, u, e.id, k) {
			return true
		}
		sizeArr[u] += sizeArr[e.to]
	}
	if sizeArr[u] == fib[k-1] || sizeArr[u] == fib[k-2] {
		cutU = u
		cutV = p
		cutEdge = edgeID
		return true
	}
	return false
}

func solve(root, n int) bool {
	k, ok := fibIndex[n]
	if !ok {
		return false
	}
	if n <= 3 {
		return true
	}
	cutU = -1
	if !find(root, -1, -1, k) {
		return false
	}
	removed[cutEdge] = true
	s1 := sizeArr[cutU]
	s2 := n - s1
	return solve(cutU, s1) && solve(cutV, s2)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(reader, &n)
	g = make([][]Edge, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		u--
		v--
		g[u] = append(g[u], Edge{to: v, id: i})
		g[v] = append(g[v], Edge{to: u, id: i})
	}
	buildFib(n)
	removed = make([]bool, n-1)
	sizeArr = make([]int, n)
	if solve(0, n) {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
