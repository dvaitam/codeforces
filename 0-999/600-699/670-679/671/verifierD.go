package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type worker struct {
	u, v int
	c    int64
}

type testCaseD struct {
	n       int
	m       int
	edges   [][2]int
	workers []worker
}

func genTree(n int) [][2]int {
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rand.Intn(i-1) + 1
		edges = append(edges, [2]int{p, i})
	}
	return edges
}

func genTestsD() []testCaseD {
	rand.Seed(4)
	tests := make([]testCaseD, 100)
	for i := range tests {
		n := rand.Intn(5) + 2 //2..6
		m := rand.Intn(5) + 1
		edges := genTree(n)
		// build parent map
		parent := make([]int, n+1)
		for _, e := range edges {
			x, y := e[0], e[1]
			parent[y] = x
		}
		workers := make([]worker, m)
		for j := 0; j < m; j++ {
			u := rand.Intn(n-1) + 2 // 2..n
			// choose ancestor from path u->1
			v := u
			for v != 1 && rand.Intn(2) == 0 {
				v = parent[v]
			}
			if v == 0 {
				v = 1
			}
			c := int64(rand.Intn(5) + 1)
			workers[j] = worker{u, v, c}
		}
		tests[i] = testCaseD{n, m, edges, workers}
	}
	return tests
}

// solver copied from 671D.go

type Node struct {
	val         int64
	to          int
	add         int64
	left, right *Node
}

func apply(n *Node, d int64) {
	if n != nil {
		n.val += d
		n.add += d
	}
}

func push(n *Node) {
	if n != nil && n.add != 0 {
		apply(n.left, n.add)
		apply(n.right, n.add)
		n.add = 0
	}
}

func merge(a, b *Node) *Node {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if a.val > b.val {
		a, b = b, a
	}
	push(a)
	a.right = merge(a.right, b)
	a.left, a.right = a.right, a.left
	return a
}

var (
	g       [][]int
	roots   []*Node
	visited []bool
)

func dfs(u, p int, top []int, weight []int64, ans *int64) bool {
	for _, v := range g[u] {
		if v != p {
			if !dfs(v, u, top, weight, ans) {
				return false
			}
			roots[u] = merge(roots[u], roots[v])
		}
	}
	visited[u] = true
	if u == 1 {
		return true
	}
	for roots[u] != nil {
		push(roots[u])
		if visited[roots[u].to] {
			roots[u] = merge(roots[u].left, roots[u].right)
		} else {
			break
		}
	}
	if roots[u] == nil {
		return false
	}
	*ans += roots[u].val
	apply(roots[u], -roots[u].val)
	return true
}

func solveD(tc testCaseD) int64 {
	n, m := tc.n, tc.m
	g = make([][]int, n+1)
	for _, e := range tc.edges {
		x, y := e[0], e[1]
		g[x] = append(g[x], y)
		g[y] = append(g[y], x)
	}
	roots = make([]*Node, n+1)
	top := make([]int, m+1)
	weight := make([]int64, m+1)
	for i, w := range tc.workers {
		idx := i + 1
		top[idx] = w.v
		weight[idx] = w.c
		node := &Node{val: w.c, to: w.v}
		roots[w.u] = merge(roots[w.u], node)
	}
	visited = make([]bool, n+1)
	var ans int64
	if !dfs(1, 0, top, weight, &ans) {
		return -1
	}
	return ans
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsD()
	for i, tc := range tests {
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.m)
		for _, e := range tc.edges {
			fmt.Fprintf(&input, "%d %d\n", e[0], e[1])
		}
		for _, w := range tc.workers {
			fmt.Fprintf(&input, "%d %d %d\n", w.u, w.v, w.c)
		}
		expect := solveD(tc)
		out, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		val, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: invalid output %q\n", i+1, out)
			os.Exit(1)
		}
		if val != expect {
			fmt.Fprintf(os.Stderr, "test %d: expected %d got %d\n", i+1, expect, val)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
