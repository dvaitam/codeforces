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
	to  int
	val int
}

type TestCase struct {
	input  string
	output string
}

func dfs(v, p int, val int, adj [][]Edge, parent, depth, edgeVal []int) {
	parent[v] = p
	edgeVal[v] = val
	for _, e := range adj[v] {
		if e.to != p {
			depth[e.to] = depth[v] + 1
			dfs(e.to, v, e.val, adj, parent, depth, edgeVal)
		}
	}
}

func pathValues(u, v int, parent, depth, edgeVal []int) []int {
	var vals []int
	for depth[u] > depth[v] {
		vals = append(vals, edgeVal[u])
		u = parent[u]
	}
	for depth[v] > depth[u] {
		vals = append(vals, edgeVal[v])
		v = parent[v]
	}
	for u != v {
		vals = append(vals, edgeVal[u])
		vals = append(vals, edgeVal[v])
		u = parent[u]
		v = parent[v]
	}
	return vals
}

func solveCaseF(n int, edges [][3]int) string {
	adj := make([][]Edge, n+1)
	for _, e := range edges {
		v, u, x := e[0], e[1], e[2]
		adj[v] = append(adj[v], Edge{u, x})
		adj[u] = append(adj[u], Edge{v, x})
	}
	parent := make([]int, n+1)
	depth := make([]int, n+1)
	edgeVal := make([]int, n+1)
	dfs(1, 0, 0, adj, parent, depth, edgeVal)
	var ans int64
	for v := 1; v <= n; v++ {
		for u := v + 1; u <= n; u++ {
			vals := pathValues(v, u, parent, depth, edgeVal)
			freq := make(map[int]int)
			for _, val := range vals {
				freq[val]++
			}
			for _, f := range freq {
				if f == 1 {
					ans++
				}
			}
		}
	}
	return fmt.Sprintf("%d\n", ans)
}

func generateTests() []TestCase {
	rand.Seed(6)
	tests := make([]TestCase, 0, 20)
	for t := 0; t < 20; t++ {
		n := rand.Intn(4) + 2
		edges := make([][3]int, n-1)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n-1; i++ {
			v := i + 1
			u := i + 2
			x := rand.Intn(3) + 1
			edges[i] = [3]int{v, u, x}
			sb.WriteString(fmt.Sprintf("%d %d %d\n", v, u, x))
		}
		out := solveCaseF(n, edges)
		tests = append(tests, TestCase{sb.String(), out})
	}
	return tests
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := generateTests()
	passed := 0
	for i, tc := range tests {
		got, err := runBinary(bin, tc.input)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			continue
		}
		g := strings.TrimSpace(got)
		e := strings.TrimSpace(tc.output)
		if g != e {
			fmt.Printf("Test %d failed. Expected %q got %q\n", i+1, e, g)
		} else {
			passed++
		}
	}
	fmt.Printf("%d/%d tests passed\n", passed, len(tests))
}
