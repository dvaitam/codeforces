package main

import (
	"bufio"
	"fmt"
	"os"
)

func checkCandidate(x int, edges [][2]int, colors []int) bool {
	for _, e := range edges {
		u, v := e[0], e[1]
		if colors[u] != colors[v] && x != u && x != v {
			return false
		}
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	edges := make([][2]int, n-1)
	for i := 0; i < n-1; i++ {
		fmt.Fscan(reader, &edges[i][0], &edges[i][1])
	}
	colors := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &colors[i])
	}

	diffFound := false
	var cand1, cand2 int
	for _, e := range edges {
		if colors[e[0]] != colors[e[1]] {
			cand1 = e[0]
			cand2 = e[1]
			diffFound = true
			break
		}
	}

	if !diffFound {
		fmt.Println("YES")
		fmt.Println(1)
		return
	}

	if checkCandidate(cand1, edges, colors) {
		fmt.Println("YES")
		fmt.Println(cand1)
		return
	}
	if checkCandidate(cand2, edges, colors) {
		fmt.Println("YES")
		fmt.Println(cand2)
		return
	}

	fmt.Println("NO")
}
