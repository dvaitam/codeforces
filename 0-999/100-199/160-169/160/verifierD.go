package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type Edge struct {
	u, v, w, id int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(45)
	var tests []string
	for t := 0; t < 100; t++ {
		n := rand.Intn(5) + 2
		m := rand.Intn(n*(n-1)/2-n+1) + n - 1
		// generate connected graph
		edges := make([]Edge, 0, m)
		used := make(map[[2]int]bool)
		// spanning tree
		for i := 2; i <= n; i++ {
			j := rand.Intn(i-1) + 1
			w := rand.Intn(10) + 1
			edges = append(edges, Edge{j, i, w, len(edges)})
			used[[2]int{j, i}] = true
			used[[2]int{i, j}] = true
		}
		for len(edges) < m {
			u := rand.Intn(n) + 1
			v := rand.Intn(n) + 1
			if u == v || used[[2]int{u, v}] {
				continue
			}
			w := rand.Intn(10) + 1
			edges = append(edges, Edge{u, v, w, len(edges)})
			used[[2]int{u, v}] = true
			used[[2]int{v, u}] = true
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", e.u, e.v, e.w))
		}
		tests = append(tests, sb.String())
	}
	for i, input := range tests {
		expect := solveD(strings.NewReader(input))
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: execution failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed: expected \n%s\n got \n%s\n", i+1, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

// DSU structure
type DSU struct {
	p, r []int
}

func NewDSU(n int) *DSU {
	p := make([]int, n+1)
	r := make([]int, n+1)
	for i := 1; i <= n; i++ {
		p[i] = i
	}
	return &DSU{p: p, r: r}
}

func (d *DSU) Find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.Find(d.p[x])
	}
	return d.p[x]
}

func (d *DSU) Union(a, b int) {
	a = d.Find(a)
	b = d.Find(b)
	if a == b {
		return
	}
	if d.r[a] < d.r[b] {
		a, b = b, a
	}
	d.p[b] = a
	if d.r[a] == d.r[b] {
		d.r[a]++
	}
}

func solveD(r io.Reader) string {
	in := bufio.NewReader(r)
	var n, m int
	fmt.Fscan(in, &n, &m)
	edges := make([]Edge, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &edges[i].u, &edges[i].v, &edges[i].w)
		edges[i].id = i
	}
	sort.Slice(edges, func(i, j int) bool { return edges[i].w < edges[j].w })
	ans := make([]string, m)
	dsu := NewDSU(n)
	for i := 0; i < m; {
		j := i
		for j < m && edges[j].w == edges[i].w {
			j++
		}
		type BE struct{ u, v, id int }
		var be []BE
		for k := i; k < j; k++ {
			u := dsu.Find(edges[k].u)
			v := dsu.Find(edges[k].v)
			if u == v {
				ans[edges[k].id] = "none"
			} else {
				be = append(be, BE{u, v, edges[k].id})
			}
		}
		if len(be) > 0 {
			comp := make(map[int]int)
			idx := 0
			for _, e := range be {
				if _, ok := comp[e.u]; !ok {
					comp[e.u] = idx
					idx++
				}
				if _, ok := comp[e.v]; !ok {
					comp[e.v] = idx
					idx++
				}
			}
			sz := idx
			type Adj struct{ to, id int }
			adj := make([][]Adj, sz)
			pairCnt := make(map[int]int)
			key := func(a, b int) int {
				if a > b {
					a, b = b, a
				}
				return a*sz + b
			}
			for _, e := range be {
				u := comp[e.u]
				v := comp[e.v]
				k := key(u, v)
				pairCnt[k]++
				adj[u] = append(adj[u], Adj{v, e.id})
				adj[v] = append(adj[v], Adj{u, e.id})
			}
			tin := make([]int, sz)
			low := make([]int, sz)
			vis := make([]bool, sz)
			timer := 0
			bridge := make(map[int]bool)
			var dfs func(u, pe int)
			dfs = func(u, pe int) {
				vis[u] = true
				timer++
				tin[u] = timer
				low[u] = timer
				for _, e := range adj[u] {
					v := e.to
					id := e.id
					if id == pe {
						continue
					}
					if vis[v] {
						if tin[v] < low[u] {
							low[u] = tin[v]
						}
					} else {
						dfs(v, id)
						if low[v] < low[u] {
							low[u] = low[v]
						}
						if low[v] > tin[u] {
							if pairCnt[key(u, v)] == 1 {
								bridge[id] = true
							}
						}
					}
				}
			}
			for u := 0; u < sz; u++ {
				if !vis[u] {
					dfs(u, -1)
				}
			}
			for _, e := range be {
				if bridge[e.id] {
					ans[e.id] = "any"
				} else {
					ans[e.id] = "at least one"
				}
			}
		}
		for k := i; k < j; k++ {
			dsu.Union(edges[k].u, edges[k].v)
		}
		i = j
	}
	var buf strings.Builder
	for i := 0; i < m; i++ {
		buf.WriteString(ans[i])
		buf.WriteByte('\n')
	}
	return buf.String()
}
