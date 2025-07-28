package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type dsu struct {
	parent []int
}

func newDSU(n int) *dsu {
	p := make([]int, n+1)
	for i := 0; i <= n; i++ {
		p[i] = i
	}
	return &dsu{parent: p}
}

func (d *dsu) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *dsu) union(x, y int) bool {
	fx := d.find(x)
	fy := d.find(y)
	if fx != fy {
		d.parent[fy] = fx
		return true
	}
	return false
}

func solveCase(n, m int, edges [][2]int) string {
	adj := make([][]int, n+1)
	uf := newDSU(n)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
		uf.union(u, v)
	}
	root := uf.find(1)
	for i := 1; i <= n; i++ {
		if uf.find(i) != root {
			return "NO\n"
		}
	}
	vis := make([]int, n+1)
	for i := 1; i <= n; i++ {
		vis[i] = -1
	}
	vis[1] = 1
	q := make([]int, 0, n)
	for _, v := range adj[1] {
		if vis[v] == -1 {
			vis[v] = 0
			q = append(q, v)
		}
	}
	for head := 0; head < len(q); head++ {
		u := q[head]
		for _, v := range adj[u] {
			if vis[v] == -1 {
				vis[v] = 1
				for _, w := range adj[v] {
					if vis[w] == -1 {
						vis[w] = 0
						q = append(q, w)
					}
				}
			}
		}
	}
	res := []int{}
	for i := 1; i <= n; i++ {
		if vis[i] == 1 {
			res = append(res, i)
		}
	}
	var sb strings.Builder
	sb.WriteString("YES\n")
	sb.WriteString(fmt.Sprintf("%d\n", len(res)))
	for i, v := range res {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 2
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(maxEdges + 1)
	edgeSet := make(map[[2]int]struct{})
	edges := make([][2]int, 0, m)
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		e := [2]int{u, v}
		if _, ok := edgeSet[e]; ok {
			continue
		}
		edgeSet[e] = struct{}{}
		edges = append(edges, e)
	}
	input := fmt.Sprintf("1\n%d %d\n", n, m)
	for _, e := range edges {
		input += fmt.Sprintf("%d %d\n", e[0], e[1])
	}
	ans := solveCase(n, m, edges)
	return input, strings.TrimSpace(ans)
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
