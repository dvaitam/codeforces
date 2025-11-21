package main

import (
	"bufio"
	"fmt"
	"os"
)

type stackItem struct {
	node   int
	parent int
	depth  int
}

const negInf = -1 << 60

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		adj := make([][]int, n)
		for i := 0; i < n-1; i++ {
			var a, b int
			fmt.Fscan(in, &a, &b)
			a--
			b--
			adj[a] = append(adj[a], b)
			adj[b] = append(adj[b], a)
		}
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		// v is equal to u in this version, ignore

		parent := make([]int, n)
		for i := range parent {
			parent[i] = -1
		}
		stack := []int{0}
		parent[0] = -2
		for len(stack) > 0 {
			vtx := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			for _, nb := range adj[vtx] {
				if parent[nb] == -1 {
					parent[nb] = vtx
					stack = append(stack, nb)
				}
			}
		}

		path := make([]int, 0)
		cur := u
		for cur != -2 {
			path = append(path, cur)
			cur = parent[cur]
		}
		for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
			path[i], path[j] = path[j], path[i]
		}
		k := len(path) - 1

		inPath := make([]bool, n)
		for _, node := range path {
			inPath[node] = true
		}

		down := make([]int, k+1)
		for idx, node := range path {
			best := 0
			for _, nb := range adj[node] {
				if idx > 0 && nb == path[idx-1] {
					continue
				}
				if idx < k && nb == path[idx+1] {
					continue
				}
				if inPath[nb] {
					continue
				}
				st := []stackItem{{nb, node, 1}}
				for len(st) > 0 {
					it := st[len(st)-1]
					st = st[:len(st)-1]
					if it.depth > best {
						best = it.depth
					}
					for _, nxt := range adj[it.node] {
						if nxt == it.parent || inPath[nxt] {
							continue
						}
						st = append(st, stackItem{nxt, it.node, it.depth + 1})
					}
				}
			}
			down[idx] = best
		}

		suffix := make([]int, k+2)
		suffix[k+1] = negInf
		for idx := k; idx >= 0; idx-- {
			val := (k - idx) + down[idx]
			if val > suffix[idx+1] {
				suffix[idx] = val
			} else {
				suffix[idx] = suffix[idx+1]
			}
		}

		mid := k / 2
		limit := k + 1
		for j := mid + 1; j <= k; j++ {
			bobLen := (k - j) + down[j]
			laPrev := (j - 1) + down[j-1]
			if bobLen >= laPrev {
				limit = j
				break
			}
		}
		upper := limit - 1
		if upper > k-1 {
			upper = k - 1
		}
		if upper < 0 {
			upper = -1
		}
		winner := "Bob"
		for i := 0; i <= upper; i++ {
			la := i + down[i]
			lb := suffix[i+1]
			if la > lb {
				winner = "Alice"
				break
			}
		}
		fmt.Fprintln(out, winner)
	}
}
