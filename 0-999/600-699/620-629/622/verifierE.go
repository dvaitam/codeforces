package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func solveE(n int, edges [][2]int) int {
	adj := make([][]int, n+1)
	for _, e := range edges {
		x, y := e[0], e[1]
		adj[x] = append(adj[x], y)
		adj[y] = append(adj[y], x)
	}

	// For each child of root 1, collect leaf depths, sort, then greedily
	// compute arrival times. Ants within the same subtree of root all pass
	// through the same child of root at distinct times, so p = max(p+1, depth).
	ans := 0
	for _, child := range adj[1] {
		var leaves []int
		var dfs func(u, par, d int)
		dfs = func(u, par, d int) {
			isLeaf := true
			for _, v := range adj[u] {
				if v == par {
					continue
				}
				dfs(v, u, d+1)
				isLeaf = false
			}
			if isLeaf {
				leaves = append(leaves, d)
			}
		}
		dfs(child, 1, 1)
		sort.Ints(leaves)
		p := 0
		for _, d := range leaves {
			if p+1 > d {
				p = p + 1
			} else {
				p = d
			}
		}
		if p > ans {
			ans = p
		}
	}
	return ans
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func randTree(n int) [][2]int {
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rand.Intn(i-1) + 1
		edges = append(edges, [2]int{p, i})
	}
	return edges
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(5)
	type Test struct {
		n     int
		edges [][2]int
	}
	var tests []Test
	tests = append(tests, Test{2, [][2]int{{1, 2}}})
	tests = append(tests, Test{3, [][2]int{{1, 2}, {1, 3}}})
	for len(tests) < 100 {
		n := rand.Intn(20) + 2
		tests = append(tests, Test{n, randTree(n)})
	}
	for idx, t := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", t.n))
		for _, e := range t.edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		input := sb.String()
		expected := fmt.Sprintf("%d", solveE(t.n, t.edges))
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("test %d failed:\ninput:\n%sexpected %s got %s\n", idx+1, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
