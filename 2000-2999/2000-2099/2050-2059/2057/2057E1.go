package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
	"sort"
)

type Edge struct {
	to int
	w  int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m, q int
		fmt.Fscan(in, &n, &m, &q)

		adj := make([][]Edge, n)
		weights := make([]int, 0, m)

		for i := 0; i < m; i++ {
			var u, v, w int
			fmt.Fscan(in, &u, &v, &w)
			u--
			v--
			adj[u] = append(adj[u], Edge{to: v, w: w})
			adj[v] = append(adj[v], Edge{to: u, w: w})
			weights = append(weights, w)
		}

		sort.Ints(weights)
		weights = unique(weights)

		mat := make([][]uint16, len(weights))
		for idx, limit := range weights {
			arr := make([]uint16, n*n)
			for s := 0; s < n; s++ {
				dist := make([]int, n)
				for i := range dist {
					dist[i] = 1 << 30
				}
				dist[s] = 0

				dq := list.New()
				dq.PushFront(s)

				for dq.Len() > 0 {
					v := dq.Remove(dq.Front()).(int)
					cur := dist[v]
					for _, e := range adj[v] {
						cost := 0
						if e.w > limit {
							cost = 1
						}
						if dist[e.to] > cur+cost {
							dist[e.to] = cur + cost
							if cost == 1 {
								dq.PushBack(e.to)
							} else {
								dq.PushFront(e.to)
							}
						}
					}
				}

				base := s * n
				for v := 0; v < n; v++ {
					arr[base+v] = uint16(dist[v])
				}
			}
			mat[idx] = arr
		}

		answers := make([]int, q)
		for i := 0; i < q; i++ {
			var a, b, k int
			fmt.Fscan(in, &a, &b, &k)
			a--
			b--
			target := uint16(k - 1)
			l, r := 0, len(weights)-1
			for l < r {
				mid := (l + r) / 2
				if mat[mid][a*n+b] <= target {
					r = mid
				} else {
					l = mid + 1
				}
			}
			answers[i] = weights[l]
		}
		for i, val := range answers {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, val)
		}
		fmt.Fprintln(out)
	}
}

func unique(arr []int) []int {
	if len(arr) == 0 {
		return arr
	}
	j := 1
	for i := 1; i < len(arr); i++ {
		if arr[i] != arr[i-1] {
			arr[j] = arr[i]
			j++
		}
	}
	return arr[:j]
}
