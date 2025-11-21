package main

import (
	"bufio"
	"fmt"
	"os"
)

type state struct {
	i, j int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var s, t string
		fmt.Fscan(in, &s)
		fmt.Fscan(in, &t)
		fmt.Fprintln(out, bfs(s, t))
	}
}

func bfs(s, t string) int {
	n := len(s)
	m := len(t)
	eq := make([]bool, min(n, m)+1)
	eq[0] = true
	for i := 1; i <= min(n, m); i++ {
		if s[i-1] == t[i-1] && eq[i-1] {
			eq[i] = true
		} else {
			break
		}
	}

	dist := make([]int, (n+1)*(m+1))
	for i := range dist {
		dist[i] = -1
	}
	idx := func(i, j int) int {
		return i*(m+1) + j
	}

	queue := make([]state, 0, (n+1)*(m+1))
	queue = append(queue, state{0, 0})
	dist[idx(0, 0)] = 0
	for head := 0; head < len(queue); head++ {
		cur := queue[head]
		i, j := cur.i, cur.j
		d := dist[idx(i, j)]
		if i == n && j == m {
			return d
		}
		// append to s
		if i < n {
			next := state{i + 1, j}
			if dist[idx(next.i, next.j)] == -1 {
				dist[idx(next.i, next.j)] = d + 1
				queue = append(queue, next)
			}
		}
		// append to t
		if j < m {
			next := state{i, j + 1}
			if dist[idx(next.i, next.j)] == -1 {
				dist[idx(next.i, next.j)] = d + 1
				queue = append(queue, next)
			}
		}
		// copy s -> t
		if i <= m && (i == 0 || (i <= len(eq)-1 && eq[i])) {
			next := state{i, i}
			if dist[idx(next.i, next.j)] == -1 {
				dist[idx(next.i, next.j)] = d + 1
				queue = append(queue, next)
			}
		}
		// copy t -> s
		if j <= n && (j == 0 || (j <= len(eq)-1 && eq[j])) {
			next := state{j, j}
			if dist[idx(next.i, next.j)] == -1 {
				dist[idx(next.i, next.j)] = d + 1
				queue = append(queue, next)
			}
		}
	}
	return -1
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
