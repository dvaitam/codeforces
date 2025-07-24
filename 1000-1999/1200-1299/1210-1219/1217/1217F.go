package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

type Query struct {
	t int
	x int
	y int
}

type DSU struct {
	parent []int
	size   []int
}

func NewDSU(n int) *DSU {
	d := &DSU{parent: make([]int, n+1), size: make([]int, n+1)}
	for i := 1; i <= n; i++ {
		d.parent[i] = i
		d.size[i] = 1
	}
	return d
}

func (d *DSU) Find(x int) int {
	for d.parent[x] != x {
		d.parent[x] = d.parent[d.parent[x]]
		x = d.parent[x]
	}
	return x
}

func (d *DSU) Union(a, b int) {
	a = d.Find(a)
	b = d.Find(b)
	if a == b {
		return
	}
	if d.size[a] < d.size[b] {
		a, b = b, a
	}
	d.parent[b] = a
	d.size[a] += d.size[b]
}

func normEdge(u, v int) [2]int {
	if u > v {
		u, v = v, u
	}
	return [2]int{u, v}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	fmt.Fscan(in, &n, &m)
	queries := make([]Query, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &queries[i].t, &queries[i].x, &queries[i].y)
	}

	block := int(math.Sqrt(float64(m))) + 1
	active := make(map[[2]int]bool) // currently active edges
	res := make([]byte, 0)
	last := 0
	vis := make([]int, n+1)
	stamp := 0

	for s := 0; s < m; s += block {
		e := s + block
		if e > m {
			e = m
		}
		// collect candidate edges in this block
		cand := make(map[[2]int]struct{})
		for i := s; i < e; i++ {
			if queries[i].t == 1 {
				x1, y1 := queries[i].x, queries[i].y
				cand[normEdge(x1, y1)] = struct{}{}
				x2 := x1 + 1
				if x2 > n {
					x2 = 1
				}
				y2 := y1 + 1
				if y2 > n {
					y2 = 1
				}
				cand[normEdge(x2, y2)] = struct{}{}
			}
		}

		// build DSU without candidate edges
		dsu := NewDSU(n)
		for eKey := range active {
			if _, ok := cand[eKey]; !ok {
				dsu.Union(eKey[0], eKey[1])
			}
		}

		root := make([]int, n+1)
		for i := 1; i <= n; i++ {
			root[i] = dsu.Find(i)
		}

		adj := make(map[int]map[int]int)
		specActive := make(map[[2]int]bool)

		// initialize adjacency with active candidate edges
		for eKey := range cand {
			if active[eKey] {
				specActive[eKey] = true
				u := root[eKey[0]]
				v := root[eKey[1]]
				if u != v {
					if adj[u] == nil {
						adj[u] = make(map[int]int)
					}
					if adj[v] == nil {
						adj[v] = make(map[int]int)
					}
					adj[u][v]++
					adj[v][u]++
				}
			} else {
				specActive[eKey] = false
			}
		}

		// process queries in block
		for i := s; i < e; i++ {
			q := queries[i]
			u, v := q.x, q.y
			if last == 1 {
				u++
				v++
				if u > n {
					u = 1
				}
				if v > n {
					v = 1
				}
			}
			if q.t == 1 {
				key := normEdge(u, v)
				if active[key] {
					active[key] = false
				} else {
					active[key] = true
				}
				if _, ok := cand[key]; ok {
					ru := root[u]
					rv := root[v]
					if ru != rv {
						if specActive[key] {
							// remove
							if adj[ru] != nil {
								adj[ru][rv]--
								if adj[ru][rv] == 0 {
									delete(adj[ru], rv)
								}
							}
							if adj[rv] != nil {
								adj[rv][ru]--
								if adj[rv][ru] == 0 {
									delete(adj[rv], ru)
								}
							}
						} else {
							if adj[ru] == nil {
								adj[ru] = make(map[int]int)
							}
							if adj[rv] == nil {
								adj[rv] = make(map[int]int)
							}
							adj[ru][rv]++
							adj[rv][ru]++
						}
					}
					specActive[key] = !specActive[key]
				}
			} else {
				ru := root[u]
				rv := root[v]
				ans := false
				if ru == rv {
					ans = true
				} else {
					stamp++
					queue := []int{ru}
					vis[ru] = stamp
					for len(queue) > 0 && !ans {
						cur := queue[0]
						queue = queue[1:]
						for to := range adj[cur] {
							if adj[cur][to] > 0 && vis[to] != stamp {
								vis[to] = stamp
								if to == rv {
									ans = true
									break
								}
								queue = append(queue, to)
							}
						}
					}
				}
				if ans {
					last = 1
					res = append(res, '1')
				} else {
					last = 0
					res = append(res, '0')
				}
			}
		}
	}

	fmt.Println(string(res))
}
