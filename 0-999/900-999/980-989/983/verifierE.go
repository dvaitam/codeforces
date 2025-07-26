package main

import (
	"bytes"
	"container/list"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const maxLog = 20

type pair struct{ a, b int }

func solveE(n int, parents []int, routes []pair, query pair) int {
	g := make([][]int, n)
	parent := make([][maxLog]int, n)
	depth := make([]int, n)
	for i := 1; i < n; i++ {
		p := parents[i-1]
		g[p] = append(g[p], i)
		g[i] = append(g[i], p)
		parent[i][0] = p
	}
	stack := []int{0}
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, to := range g[v] {
			if to == parent[v][0] && v != 0 {
				continue
			}
			parent[to][0] = v
			depth[to] = depth[v] + 1
			stack = append(stack, to)
		}
	}
	for j := 1; j < maxLog; j++ {
		for i := 0; i < n; i++ {
			parent[i][j] = parent[parent[i][j-1]][j-1]
		}
	}
	lca := func(a, b int) int {
		if depth[a] < depth[b] {
			a, b = b, a
		}
		diff := depth[a] - depth[b]
		for i := 0; i < maxLog; i++ {
			if diff>>i&1 == 1 {
				a = parent[a][i]
			}
		}
		if a == b {
			return a
		}
		for i := maxLog - 1; i >= 0; i-- {
			if parent[a][i] != parent[b][i] {
				a = parent[a][i]
				b = parent[b][i]
			}
		}
		return parent[a][0]
	}

	m := len(routes)
	routeNodes := make([][]int, m)
	cityRoutes := make([][]int, n)
	for i, r := range routes {
		a := r.a
		b := r.b
		gNode := lca(a, b)
		path := []int{}
		x := a
		for x != gNode {
			path = append(path, x)
			x = parent[x][0]
		}
		tmp := []int{}
		y := b
		for y != gNode {
			tmp = append(tmp, y)
			y = parent[y][0]
		}
		path = append(path, gNode)
		for j := len(tmp) - 1; j >= 0; j-- {
			path = append(path, tmp[j])
		}
		routeNodes[i] = path
		for _, v := range path {
			cityRoutes[v] = append(cityRoutes[v], i)
		}
	}

	u := query.a
	v := query.b
	tot := n + m
	dist := make([]int, tot)
	for i := range dist {
		dist[i] = -1
	}
	dq := list.New()
	dist[u] = 0
	dq.PushBack(u)
	for dq.Len() > 0 {
		e := dq.Front()
		dq.Remove(e)
		x := e.Value.(int)
		if x == v {
			break
		}
		if x < n {
			for _, r := range cityRoutes[x] {
				y := n + r
				if dist[y] == -1 || dist[y] > dist[x]+1 {
					dist[y] = dist[x] + 1
					dq.PushBack(y)
				}
			}
		} else {
			r := x - n
			for _, y := range routeNodes[r] {
				if dist[y] == -1 || dist[y] > dist[x] {
					dist[y] = dist[x]
					dq.PushFront(y)
				}
			}
		}
	}
	return dist[v]
}

func generateE(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 2
	parents := make([]int, n-1)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 1; i < n; i++ {
		p := rng.Intn(i)
		parents[i-1] = p
		fmt.Fprintf(&sb, "%d\n", p+1)
	}
	m := rng.Intn(3) + 1
	routes := make([]pair, m)
	for i := 0; i < m; i++ {
		a := rng.Intn(n)
		b := rng.Intn(n)
		fmt.Fprintf(&sb, "%d %d\n", a+1, b+1)
		routes[i] = pair{a, b}
	}
	fmt.Fprintf(&sb, "1\n")
	u := rng.Intn(n)
	v := rng.Intn(n)
	fmt.Fprintf(&sb, "%d %d\n", u+1, v+1)
	res := solveE(n, parents, routes, pair{u, v})
	return sb.String(), fmt.Sprint(res)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(46))
	for i := 0; i < 100; i++ {
		in, exp := generateE(rng)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(in)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			fmt.Printf("case %d runtime error: %v\n%s", i+1, err, out.String())
			return
		}
		got := strings.TrimSpace(out.String())
		if got != exp {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, in, exp, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
