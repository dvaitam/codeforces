package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min3(a, b, c int) int {
	return min(a, min(b, c))
}

func ceil_div2(a int) int {
	if a <= 0 {
		return 0
	}
	return (a + 1) / 2
}

var m, n int
var parent [3][2][100005]int
var visited [3][2][100005]bool

func find(t, b, k int) int {
	if parent[t][b][k] == k {
		return k
	}
	parent[t][b][k] = find(t, b, parent[t][b][k])
	return parent[t][b][k]
}

type TK struct{ t, k int }

func get_types(G, W int) []TK {
	var res []TK
	if G == 0 {
		res = append(res, TK{0, W})
	}
	if G == m {
		res = append(res, TK{2, W})
	}
	if G == W {
		res = append(res, TK{1, W})
	}
	return res
}

type Range struct{ t, L, R int }

func get_b0_transitions(t, k int) []Range {
	var res []Range
	if t == 0 {
		if k >= 1 {
			res = append(res, Range{0, max(0, k-n), k - 1})
		}
	} else if t == 2 {
		if k >= 1 {
			res = append(res, Range{2, max(0, k-n), k - 1})
		}
		L := max(0, ceil_div2(m+k-n))
		R := min(k, (m+k-1)/2)
		if L <= R {
			res = append(res, Range{1, L, R})
		}
		if n >= m {
			max_y := min3(k, n-m, m)
			if max_y >= 0 {
				res = append(res, Range{0, k - max_y, k})
			}
		}
	} else if t == 1 {
		if k >= 1 && n >= 2 {
			res = append(res, Range{1, max(0, k-n/2), k - 1})
		}
		if k >= 1 && n >= k {
			res = append(res, Range{0, max(0, k-(n-k)), k})
		}
	}
	return res
}

func get_GW(t, k int) (int, int) {
	if t == 0 {
		return 0, k
	}
	if t == 1 {
		return k, k
	}
	return m, k
}

func mark_visited(G, W, b int) {
	for _, tk := range get_types(G, W) {
		if !visited[tk.t][b][tk.k] {
			visited[tk.t][b][tk.k] = true
			parent[tk.t][b][tk.k] = find(tk.t, b, tk.k+1)
		}
	}
}

type State struct {
	G, W, b, dist int
}

func solve() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fscan(reader, &m, &n)

	for t := 0; t < 3; t++ {
		for b := 0; b < 2; b++ {
			for k := 0; k <= m+1; k++ {
				parent[t][b][k] = k
			}
		}
	}

	queue := make([]State, 0, 1000000)
	queue = append(queue, State{m, m, 0, 0})
	mark_visited(m, m, 0)

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr.G == 0 && curr.W == 0 && curr.b == 1 {
			fmt.Println(curr.dist)
			return
		}

		var targets []Range
		if curr.b == 0 {
			for _, tk := range get_types(curr.G, curr.W) {
				targets = append(targets, get_b0_transitions(tk.t, tk.k)...)
			}
		} else {
			G_d := m - curr.G
			W_d := m - curr.W
			for _, tk_d := range get_types(G_d, W_d) {
				targets_d := get_b0_transitions(tk_d.t, tk_d.k)
				for _, r := range targets_d {
					var st int
					if r.t == 0 {
						st = 2
					} else if r.t == 1 {
						st = 1
					} else if r.t == 2 {
						st = 0
					}
					targets = append(targets, Range{st, m - r.R, m - r.L})
				}
			}
		}

		for _, r := range targets {
			k := find(r.t, 1-curr.b, r.L)
			for k <= r.R {
				G_new, W_new := get_GW(r.t, k)
				queue = append(queue, State{G_new, W_new, 1 - curr.b, curr.dist + 1})
				mark_visited(G_new, W_new, 1-curr.b)
				k = find(r.t, 1-curr.b, k+1)
			}
		}
	}

	fmt.Println("-1")
}

func main() {
	solve()
}
