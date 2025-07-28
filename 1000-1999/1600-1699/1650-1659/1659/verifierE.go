package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type edge struct{ u, v, w int }

func run(bin, in string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func runRef(in string) (string, error) {
	cmd := exec.Command("go", "run", "1659E.go")
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genCase() string {
	n := rand.Intn(7) + 2
	maxEdges := n * (n - 1) / 2
	m := rand.Intn(min(maxEdges, 10)-(n-1)+1) + (n - 1)
	edges := make([]edge, 0, m)
	parent := make([]int, n+1)
	for i := 1; i <= n; i++ {
		parent[i] = i
	}
	var find func(int) int
	find = func(x int) int {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}
	union := func(a, b int) {
		fa, fb := find(a), find(b)
		if fa != fb {
			parent[fb] = fa
		}
	}
	// Build tree edges
	for i := 2; i <= n; i++ {
		j := rand.Intn(i-1) + 1
		w := rand.Intn(16)
		edges = append(edges, edge{j, i, w})
		union(j, i)
	}
	used := make(map[string]bool)
	for _, e := range edges {
		used[fmt.Sprintf("%d-%d", min(e.u, e.v), max(e.u, e.v))] = true
	}
	for len(edges) < m {
		u := rand.Intn(n) + 1
		v := rand.Intn(n) + 1
		if u == v {
			continue
		}
		key := fmt.Sprintf("%d-%d", min(u, v), max(u, v))
		if used[key] {
			continue
		}
		used[key] = true
		w := rand.Intn(16)
		edges = append(edges, edge{u, v, w})
	}
	q := rand.Intn(5) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, len(edges))
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d %d\n", e.u, e.v, e.w)
	}
	fmt.Fprintf(&sb, "%d\n", q)
	for i := 0; i < q; i++ {
		u := rand.Intn(n) + 1
		v := rand.Intn(n) + 1
		fmt.Fprintf(&sb, "%d %d\n", u, v)
	}
	return sb.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(1)
	for t := 1; t <= 100; t++ {
		input := genCase()
		exp, err := runRef(input)
		if err != nil {
			fmt.Println("reference failed:", err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d exec failed: %v\n", t, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Printf("test %d failed\nexpected: %s\ngot: %s\n", t, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("ok")
}
