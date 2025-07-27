package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Edge struct {
	a int
	b int
}

type Query struct {
	typ int
	arg int
}

func runProg(prog string, args []string, input string) (string, error) {
	cmd := exec.Command(prog, args...)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(4)
	for t := 1; t <= 100; t++ {
		n := rand.Intn(4) + 1
		maxEdges := n * (n - 1) / 2
		m := rand.Intn(maxEdges) + 1
		perm := rand.Perm(n)
		p := make([]int, n+1)
		for i := 0; i < n; i++ {
			p[i+1] = perm[i] + 1
		}
		edges := make([]Edge, 0, m)
		seen := make(map[[2]int]bool)
		for len(edges) < m {
			u := rand.Intn(n) + 1
			v := rand.Intn(n) + 1
			if u == v {
				continue
			}
			if u > v {
				u, v = v, u
			}
			key := [2]int{u, v}
			if !seen[key] {
				seen[key] = true
				edges = append(edges, Edge{u, v})
			}
		}
		q := rand.Intn(m+3) + 1
		queries := make([]Query, 0, q)
		remaining := make([]int, m)
		for i := 0; i < m; i++ {
			remaining[i] = i + 1
		}
		for len(queries) < q {
			if len(remaining) == 0 || rand.Intn(2) == 0 {
				v := rand.Intn(n) + 1
				queries = append(queries, Query{1, v})
			} else {
				idx := rand.Intn(len(remaining))
				e := remaining[idx]
				remaining = append(remaining[:idx], remaining[idx+1:]...)
				queries = append(queries, Query{2, e})
			}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, q))
		for i := 1; i <= n; i++ {
			if i > 1 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", p[i]))
		}
		sb.WriteByte('\n')
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e.a, e.b))
		}
		for _, qq := range queries {
			sb.WriteString(fmt.Sprintf("%d %d\n", qq.typ, qq.arg))
		}
		input := sb.String()
		candOut, err := runProg(binary, nil, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\n", t, err)
			os.Exit(1)
		}
		refOut, err := runProg("go", []string{"run", "1416D.go"}, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", t, err)
			os.Exit(1)
		}
		if candOut != refOut {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\nexpected: %s\ngot: %s\n", t, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
