package main

import (
	"bufio"
	"fmt"
	"os"
)

// find-next disjoint set for numbers
var next []int

func find(x int) int {
	if next[x] != x {
		next[x] = find(next[x])
	}
	return next[x]
}

func removeNode(x int) {
	next[x] = find(x + 1)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m, k int
	if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
		return
	}

	forb1 := make([]bool, n+1)
	banned := make([]map[int]struct{}, n+1)
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		if a == 1 || b == 1 {
			other := b
			if a == 1 {
				other = b
			} else {
				other = a
			}
			forb1[other] = true
		} else {
			if banned[a] == nil {
				banned[a] = make(map[int]struct{})
			}
			if banned[b] == nil {
				banned[b] = make(map[int]struct{})
			}
			banned[a][b] = struct{}{}
			banned[b][a] = struct{}{}
		}
	}

	forb1Count := 0
	for i := 2; i <= n; i++ {
		if forb1[i] {
			forb1Count++
		}
	}
	maxDeg1 := n - 1 - forb1Count
	if k > maxDeg1 {
		fmt.Println("impossible")
		return
	}

	next = make([]int, n+2)
	for i := 2; i <= n+1; i++ {
		next[i] = i
	}
	visited := make([]bool, n+1)
	blocked := make([]bool, n+2)

	compCount := 0

	for {
		start := find(2)
		if start > n {
			break
		}
		queue := []int{start}
		removeNode(start)
		visited[start] = true
		accessible := !forb1[start]

		for len(queue) > 0 {
			v := queue[0]
			queue = queue[1:]
			for to := range banned[v] {
				blocked[to] = true
			}
			x := find(2)
			for x <= n {
				if !blocked[x] {
					queue = append(queue, x)
					removeNode(x)
					visited[x] = true
					if !forb1[x] {
						accessible = true
					}
					x = find(x)
				} else {
					x = find(x + 1)
				}
			}
			for to := range banned[v] {
				blocked[to] = false
			}
		}

		if !accessible {
			fmt.Println("impossible")
			return
		}
		compCount++
	}

	if compCount <= k {
		fmt.Println("possible")
	} else {
		fmt.Println("impossible")
	}
}
