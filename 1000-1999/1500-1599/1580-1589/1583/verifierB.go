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

type restriction struct{ a, b, c int }

type caseB struct {
	n int
	m int
	r []restriction
}

func genCase(rng *rand.Rand) caseB {
	n := rng.Intn(8) + 3
	m := rng.Intn(n-1) + 1
	r := make([]restriction, m)
	for i := 0; i < m; i++ {
		var a, b, c int
		for {
			a = rng.Intn(n) + 1
			b = rng.Intn(n) + 1
			c = rng.Intn(n) + 1
			if a != b && b != c && a != c {
				break
			}
		}
		r[i] = restriction{a, b, c}
	}
	return caseB{n, m, r}
}

func runCase(bin string, tc caseB) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for _, rr := range tc.r {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", rr.a, rr.b, rr.c))
	}
	input := sb.String()
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) != 2*(tc.n-1) {
		return fmt.Errorf("expected %d numbers got %d", 2*(tc.n-1), len(fields))
	}
	edges := make([][2]int, tc.n-1)
	idx := 0
	for i := 0; i < tc.n-1; i++ {
		var u, v int
		fmt.Sscan(fields[idx], &u)
		fmt.Sscan(fields[idx+1], &v)
		idx += 2
		if u < 1 || u > tc.n || v < 1 || v > tc.n || u == v {
			return fmt.Errorf("invalid edge %d %d", u, v)
		}
		edges[i] = [2]int{u, v}
	}
	// check tree using DSU
	parent := make([]int, tc.n+1)
	for i := 1; i <= tc.n; i++ {
		parent[i] = i
	}
	var find func(int) int
	find = func(x int) int {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}
	unite := func(a, b int) bool {
		ra, rb := find(a), find(b)
		if ra == rb {
			return false
		}
		parent[rb] = ra
		return true
	}
	for _, e := range edges {
		if !unite(e[0], e[1]) {
			return fmt.Errorf("edges form cycle")
		}
	}
	root := find(1)
	for i := 2; i <= tc.n; i++ {
		if find(i) != root {
			return fmt.Errorf("graph not connected")
		}
	}

	// build adjacency
	g := make([][]int, tc.n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	// function to check if b is on path from a to c
	isOnPath := func(a, b, c int) bool {
		prev := make([]int, tc.n+1)
		for i := 1; i <= tc.n; i++ {
			prev[i] = -1
		}
		queue := []int{a}
		prev[a] = 0
		for len(queue) > 0 {
			v := queue[0]
			queue = queue[1:]
			if v == c {
				break
			}
			for _, to := range g[v] {
				if prev[to] == -1 {
					prev[to] = v
					queue = append(queue, to)
				}
			}
		}
		if prev[c] == -1 {
			return false
		}
		x := c
		for x != a {
			if x == b {
				return true
			}
			x = prev[x]
		}
		if a == b {
			return true
		}
		return false
	}
	for _, rr := range tc.r {
		if isOnPath(rr.a, rr.b, rr.c) {
			return fmt.Errorf("restriction violated a=%d b=%d c=%d", rr.a, rr.b, rr.c)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
