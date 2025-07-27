package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveE(n int, edges [][2]int, queries [][]int) []string {
	g := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	parent := make([]int, n+1)
	depth := make([]int, n+1)
	tin := make([]int, n+1)
	tout := make([]int, n+1)
	timer := 0
	type frame struct{ v, p, idx int }
	stack := []frame{{1, 0, -1}}
	for len(stack) > 0 {
		f := &stack[len(stack)-1]
		if f.idx == -1 {
			parent[f.v] = f.p
			depth[f.v] = depth[f.p] + 1
			timer++
			tin[f.v] = timer
			f.idx = 0
		}
		if f.idx < len(g[f.v]) {
			to := g[f.v][f.idx]
			f.idx++
			if to == f.p {
				continue
			}
			stack = append(stack, frame{to, f.v, -1})
		} else {
			timer++
			tout[f.v] = timer
			stack = stack[:len(stack)-1]
		}
	}
	isAncestor := func(u, v int) bool {
		return tin[u] <= tin[v] && tout[v] <= tout[u]
	}
	res := make([]string, len(queries))
	for idx, q := range queries {
		deepest := 1
		nodes := make([]int, len(q))
		for i, v := range q {
			if v != 1 {
				v = parent[v]
			}
			nodes[i] = v
			if depth[v] > depth[deepest] {
				deepest = v
			}
		}
		ok := true
		for _, v := range nodes {
			if !isAncestor(v, deepest) {
				ok = false
				break
			}
		}
		if ok {
			res[idx] = "YES"
		} else {
			res[idx] = "NO"
		}
	}
	return res
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(5)
	for t := 1; t <= 100; t++ {
		n := rand.Intn(8) + 2
		edges := make([][2]int, n-1)
		for i := 2; i <= n; i++ {
			p := rand.Intn(i-1) + 1
			edges[i-2] = [2]int{p, i}
		}
		m := rand.Intn(5) + 1
		queries := make([][]int, m)
		for i := 0; i < m; i++ {
			k := rand.Intn(n) + 1
			used := make(map[int]bool)
			q := make([]int, k)
			for j := 0; j < k; j++ {
				v := rand.Intn(n) + 1
				for used[v] {
					v = rand.Intn(n) + 1
				}
				used[v] = true
				q[j] = v
			}
			queries[i] = q
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		for _, q := range queries {
			sb.WriteString(fmt.Sprintf("%d", len(q)))
			for _, v := range q {
				sb.WriteString(fmt.Sprintf(" %d", v))
			}
			sb.WriteByte('\n')
		}
		expectLines := solveE(n, edges, queries)
		expect := strings.Join(expectLines, "\n")
		out, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%s", t, err, sb.String())
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "test %d failed: expected\n%s\ngot\n%s\ninput:\n%s", t, expect, out, sb.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
