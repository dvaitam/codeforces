package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	adj := make([]map[int]struct{}, n+1)
	for i := 0; i < m; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		if adj[x] == nil {
			adj[x] = make(map[int]struct{})
		}
		if adj[y] == nil {
			adj[y] = make(map[int]struct{})
		}
		adj[x][y] = struct{}{}
		adj[y][x] = struct{}{}
	}

	next := make([]int, n+1)
	prev := make([]int, n+1)
	for i := 1; i <= n; i++ {
		next[i] = i + 1
		prev[i] = i - 1
	}
	next[n] = 0
	head := 1

	inSet := make([]bool, n+1)
	for i := 1; i <= n; i++ {
		inSet[i] = true
	}

	remove := func(v int) {
		if prev[v] != 0 {
			next[prev[v]] = next[v]
		} else {
			head = next[v]
		}
		if next[v] != 0 {
			prev[next[v]] = prev[v]
		}
		next[v] = 0
		prev[v] = 0
		inSet[v] = false
	}

	comps := []int{}
	for head != 0 {
		start := head
		remove(start)
		queue := []int{start}
		size := 0
		for len(queue) > 0 {
			v := queue[0]
			queue = queue[1:]
			size++
			u := head
			for u != 0 {
				if _, ok := adj[v][u]; ok {
					u = next[u]
					continue
				}
				nextU := next[u]
				remove(u)
				queue = append(queue, u)
				u = nextU
			}
		}
		comps = append(comps, size)
	}

	sort.Ints(comps)
	fmt.Fprintln(out, len(comps))
	for i, v := range comps {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, v)
	}
	fmt.Fprintln(out)
}
