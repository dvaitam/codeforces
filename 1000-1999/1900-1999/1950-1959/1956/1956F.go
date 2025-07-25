package main

import (
	"bufio"
	"fmt"
	"os"
)

func find(parent []int, x int) int {
	for parent[x] != x {
		parent[x] = parent[parent[x]]
		x = parent[x]
	}
	return x
}

func remove(parent []int, x int) {
	parent[x] = find(parent, x+1)
}

func processRange(parent []int, visited []bool, l, r []int, n int, u int, L int, R int, queue *[]int) {
	if L < 1 {
		L = 1
	}
	if R > n {
		R = n
	}
	for L <= R {
		j := find(parent, L)
		if j > R {
			break
		}
		diff := u - j
		if diff < 0 {
			diff = -diff
		}
		if diff >= l[u]+l[j] && diff <= r[u]+r[j] {
			visited[j] = true
			remove(parent, j)
			*queue = append(*queue, j)
			L = j
		} else {
			L = j + 1
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		l := make([]int, n+1)
		r := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &l[i], &r[i])
		}
		parent := make([]int, n+2)
		for i := 1; i <= n+1; i++ {
			parent[i] = i
		}
		visited := make([]bool, n+2)
		components := 0
		queue := make([]int, 0)

		for i := 1; i <= n; i++ {
			if visited[i] {
				continue
			}
			components++
			visited[i] = true
			remove(parent, i)
			queue = append(queue, i)
			for len(queue) > 0 {
				u := queue[0]
				queue = queue[1:]

				// left range
				L1 := u - r[u]
				R1 := u - l[u]
				processRange(parent, visited, l, r, n, u, L1, R1, &queue)
				// right range
				L2 := u + l[u]
				R2 := u + r[u]
				processRange(parent, visited, l, r, n, u, L2, R2, &queue)
			}
		}
		fmt.Fprintln(out, components)
	}
}
