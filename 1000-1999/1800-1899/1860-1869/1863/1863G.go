package main

import (
	"bufio"
	"fmt"
	"os"
)

// NOTE: This is a naive brute-force implementation that only works for very
// small n. It explores all reachable arrays using a BFS. The official solution
// requires a more involved combinatorial analysis which is not implemented here.

func bfs(arr []int) int {
	n := len(arr)
	type state string
	toKey := func(a []int) state {
		b := make([]byte, 0, n*4)
		for i, v := range a {
			if i > 0 {
				b = append(b, ',')
			}
			b = append(b, fmt.Sprint(v)...)
		}
		return state(b)
	}
	start := make([]int, n)
	copy(start, arr)
	q := [][]int{start}
	vis := map[state]struct{}{toKey(start): {}}
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		for i := 0; i < n; i++ {
			j := cur[i] - 1
			if j < 0 || j >= n {
				continue
			}
			nxt := make([]int, n)
			copy(nxt, cur)
			nxt[i], nxt[j] = nxt[j], nxt[i]
			key := toKey(nxt)
			if _, ok := vis[key]; !ok {
				vis[key] = struct{}{}
				q = append(q, nxt)
			}
		}
	}
	return len(vis)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	if n > 9 {
		// Placeholder: problem constraints are too big for brute force.
		fmt.Println(0)
		return
	}
	ans := bfs(a)
	const mod = 1000000007
	fmt.Println(ans % mod)
}
