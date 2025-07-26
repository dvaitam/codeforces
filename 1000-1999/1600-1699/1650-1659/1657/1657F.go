package main

import (
	"bufio"
	"fmt"
	"os"
)

type Occ struct {
	q   int
	pos int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}
	tree := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
	}

	parent := make([]int, n+1)
	depth := make([]int, n+1)
	bfs := []int{1}
	parent[1] = 0
	for idx := 0; idx < len(bfs); idx++ {
		u := bfs[idx]
		for _, v := range tree[u] {
			if v == parent[u] {
				continue
			}
			parent[v] = u
			depth[v] = depth[u] + 1
			bfs = append(bfs, v)
		}
	}

	type Query struct {
		x, y   int
		s      []byte
		path   []int
		orient int
	}

	queries := make([]Query, q)
	adj := make([][]Occ, n+1)
	letters := make([]byte, n+1)
	queueV := make([]int, 0)
	queueQ := make([]int, 0)

	getPath := func(x, y int) []int {
		px := make([]int, 0)
		py := make([]int, 0)
		for x != y {
			if depth[x] > depth[y] {
				px = append(px, x)
				x = parent[x]
			} else if depth[y] > depth[x] {
				py = append(py, y)
				y = parent[y]
			} else {
				px = append(px, x)
				py = append(py, y)
				x = parent[x]
				y = parent[y]
			}
		}
		px = append(px, x)
		for i := len(py) - 1; i >= 0; i-- {
			px = append(px, py[i])
		}
		return px
	}

	assignLetter := func(v int, c byte) bool {
		if letters[v] == 0 {
			letters[v] = c
			queueV = append(queueV, v)
			return true
		}
		return letters[v] == c
	}

	assignOrient := func(idx int, o int) bool {
		if queries[idx].orient == -1 {
			queries[idx].orient = o
			queueQ = append(queueQ, idx)
			return true
		}
		return queries[idx].orient == o
	}

	for i := 0; i < q; i++ {
		var x, y int
		var s string
		fmt.Fscan(in, &x, &y, &s)
		p := getPath(x, y)
		queries[i] = Query{x: x, y: y, s: []byte(s), path: p, orient: -1}
		m := len(p)
		for j := 0; j < m; j++ {
			v := p[j]
			adj[v] = append(adj[v], Occ{q: i, pos: j})
		}
		for j := 0; j < m; j++ {
			a := s[j]
			b := s[m-1-j]
			if a == b {
				if !assignLetter(p[j], a) {
					fmt.Fprintln(out, "NO")
					return
				}
			}
		}
	}

	for len(queueV) > 0 || len(queueQ) > 0 {
		for len(queueV) > 0 {
			v := queueV[0]
			queueV = queueV[1:]
			for _, oc := range adj[v] {
				qi := oc.q
				pos := oc.pos
				q := &queries[qi]
				a := q.s[pos]
				b := q.s[len(q.s)-1-pos]
				if q.orient == -1 {
					if letters[v] == a && letters[v] != b {
						if !assignOrient(qi, 0) {
							fmt.Fprintln(out, "NO")
							return
						}
					} else if letters[v] != a && letters[v] == b {
						if !assignOrient(qi, 1) {
							fmt.Fprintln(out, "NO")
							return
						}
					} else if letters[v] != a && letters[v] != b {
						fmt.Fprintln(out, "NO")
						return
					}
				} else {
					expected := a
					if q.orient == 1 {
						expected = b
					}
					if letters[v] != expected {
						fmt.Fprintln(out, "NO")
						return
					}
				}
			}
		}
		for len(queueQ) > 0 {
			qi := queueQ[0]
			queueQ = queueQ[1:]
			q := &queries[qi]
			m := len(q.path)
			for j := 0; j < m; j++ {
				v := q.path[j]
				c := q.s[j]
				if q.orient == 1 {
					c = q.s[m-1-j]
				}
				if !assignLetter(v, c) {
					fmt.Fprintln(out, "NO")
					return
				}
			}
		}
	}

	// assign remaining queries and vertices arbitrarily
	for i := 0; i < q; i++ {
		if queries[i].orient == -1 {
			if !assignOrient(i, 0) {
				fmt.Fprintln(out, "NO")
				return
			}
			for len(queueQ) > 0 {
				qi := queueQ[0]
				queueQ = queueQ[1:]
				q := &queries[qi]
				m := len(q.path)
				for j := 0; j < m; j++ {
					v := q.path[j]
					c := q.s[j]
					if q.orient == 1 {
						c = q.s[m-1-j]
					}
					if !assignLetter(v, c) {
						fmt.Fprintln(out, "NO")
						return
					}
				}
			}
		}
	}

	for i := 1; i <= n; i++ {
		if letters[i] == 0 {
			letters[i] = 'a'
		}
	}
	fmt.Fprintln(out, "YES")
	fmt.Fprintln(out, string(letters[1:]))
}
