package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
)

type edge struct{ u, v int }

type testC struct {
	n     int
	edges []edge
}

func genTestsC() []testC {
	rand.Seed(3)
	tests := make([]testC, 100)
	for i := range tests {
		n := rand.Intn(10) + 3
		edges := make([]edge, n-1)
		for v := 2; v <= n; v++ {
			u := rand.Intn(v-1) + 1
			edges[v-2] = edge{u: u, v: v}
		}
		tests[i] = testC{n: n, edges: edges}
	}
	return tests
}

func buildGraph(n int, edges []edge) [][]int {
	g := make([][]int, n+1)
	for _, e := range edges {
		g[e.u] = append(g[e.u], e.v)
		g[e.v] = append(g[e.v], e.u)
	}
	return g
}

func findCentroids(n int, g [][]int) []int {
	size := make([]int, n+1)
	best := n + 1
	var cents []int
	var dfs func(int, int)
	dfs = func(u, p int) {
		size[u] = 1
		maxSub := 0
		for _, v := range g[u] {
			if v == p {
				continue
			}
			dfs(v, u)
			size[u] += size[v]
			if size[v] > maxSub {
				maxSub = size[v]
			}
		}
		if n-size[u] > maxSub {
			maxSub = n - size[u]
		}
		if maxSub < best {
			best = maxSub
			cents = []int{u}
		} else if maxSub == best {
			cents = append(cents, u)
		}
	}
	dfs(1, 0)
	return cents
}

func isTree(n int, edges []edge) bool {
	if len(edges) != n-1 {
		return false
	}
	g := buildGraph(n, edges)
	visited := make([]bool, n+1)
	var stack []int
	stack = append(stack, 1)
	visited[1] = true
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, v := range g[u] {
			if !visited[v] {
				visited[v] = true
				stack = append(stack, v)
			}
		}
	}
	for i := 1; i <= n; i++ {
		if !visited[i] {
			return false
		}
	}
	return true
}

func verify(tc testC, x1, y1, x2, y2 int) bool {
	// check edge exists
	ok := false
	for i, e := range tc.edges {
		if (e.u == x1 && e.v == y1) || (e.u == y1 && e.v == x1) {
			tc.edges[i] = tc.edges[len(tc.edges)-1]
			tc.edges = tc.edges[:len(tc.edges)-1]
			ok = true
			break
		}
	}
	if !ok {
		return false
	}
	if x2 < 1 || x2 > tc.n || y2 < 1 || y2 > tc.n || x1 < 1 || x1 > tc.n || y1 < 1 || y1 > tc.n {
		return false
	}
	tc.edges = append(tc.edges, edge{u: x2, v: y2})
	if !isTree(tc.n, tc.edges) {
		return false
	}
	cents := findCentroids(tc.n, buildGraph(tc.n, tc.edges))
	return len(cents) == 1
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsC()
	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintln(&input, tc.n)
		for _, e := range tc.edges {
			fmt.Fprintf(&input, "%d %d\n", e.u, e.v)
		}
	}

	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\noutput:\n%s\n", err, out.String())
		os.Exit(1)
	}

	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	scanner.Split(bufio.ScanWords)
	for i, tc := range tests {
		var vals [4]int
		for j := 0; j < 4; j++ {
			if !scanner.Scan() {
				fmt.Fprintf(os.Stderr, "wrong output format on test %d\n", i+1)
				os.Exit(1)
			}
			fmt.Sscan(scanner.Text(), &vals[j])
		}
		if !verify(tc, vals[0], vals[1], vals[2], vals[3]) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\n", i+1)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output")
		os.Exit(1)
	}
	fmt.Println("Accepted")
}
