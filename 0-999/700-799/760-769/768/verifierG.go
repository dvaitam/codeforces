package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	n      int
	edges  [][2]int // each pair {parent, child}, parent=0 indicates root
	expect []int
}

func copyAdj(adj [][]int) [][]int {
	res := make([][]int, len(adj))
	for i := range adj {
		res[i] = append([]int(nil), adj[i]...)
	}
	return res
}

func removeEdge(adj [][]int, a, b int) {
	for i, v := range adj[a] {
		if v == b {
			adj[a] = append(adj[a][:i], adj[a][i+1:]...)
			break
		}
	}
	for i, v := range adj[b] {
		if v == a {
			adj[b] = append(adj[b][:i], adj[b][i+1:]...)
			break
		}
	}
}

func addEdge(adj [][]int, a, b int) {
	adj[a] = append(adj[a], b)
	adj[b] = append(adj[b], a)
}

func largestComponent(adj [][]int, remove int) int {
	n := len(adj) - 1
	visited := make([]bool, n+1)
	maxSize := 0
	for i := 1; i <= n; i++ {
		if i == remove || visited[i] {
			continue
		}
		size := 0
		stack := []int{i}
		visited[i] = true
		for len(stack) > 0 {
			v := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			size++
			for _, w := range adj[v] {
				if w == remove || visited[w] {
					continue
				}
				visited[w] = true
				stack = append(stack, w)
			}
		}
		if size > maxSize {
			maxSize = size
		}
	}
	return maxSize
}

func compute(n int, parent []int, remove int) int {
	adj := make([][]int, n+1)
	for v := 1; v <= n; v++ {
		if v == remove {
			continue
		}
		p := parent[v]
		if p == remove || p == 0 {
			continue
		}
		adj[v] = append(adj[v], p)
		adj[p] = append(adj[p], v)
	}
	best := largestComponent(adj, remove)
	for u := 1; u <= n; u++ {
		if u == remove {
			continue
		}
		p := parent[u]
		if p == 0 || p == remove {
			continue
		}
		adj2 := copyAdj(adj)
		removeEdge(adj2, u, p)
		for w := 1; w <= n; w++ {
			if w == remove || w == u {
				continue
			}
			adj3 := copyAdj(adj2)
			addEdge(adj3, u, w)
			size := largestComponent(adj3, remove)
			if size < best {
				best = size
			}
		}
	}
	return best
}

func bruteAnswer(n int, edges [][2]int) []int {
	parent := make([]int, n+1)
	for _, e := range edges {
		parent[e[1]] = e[0]
	}
	res := make([]int, n)
	for rm := 1; rm <= n; rm++ {
		res[rm-1] = compute(n, parent, rm)
	}
	return res
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 100)
	// small deterministic tree
	edges := [][2]int{{0, 1}, {1, 2}, {1, 3}}
	tests = append(tests, testCase{n: 3, edges: edges, expect: bruteAnswer(3, edges)})
	for len(tests) < 100 {
		n := rng.Intn(5) + 2 // 2..6
		edges := make([][2]int, n)
		edges[0] = [2]int{0, 1}
		for i := 2; i <= n; i++ {
			p := rng.Intn(i-1) + 1
			edges[i-1] = [2]int{p, i}
		}
		rng.Shuffle(n, func(i, j int) { edges[i], edges[j] = edges[j], edges[i] })
		tests = append(tests, testCase{n: n, edges: edges, expect: bruteAnswer(n, edges)})
	}
	return tests
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", tc.n))
		for _, e := range tc.edges {
			input.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		expectParts := make([]string, len(tc.expect))
		for j, v := range tc.expect {
			expectParts[j] = strconv.Itoa(v)
		}
		expected := strings.Join(expectParts, "\n")
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input.String())
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected:\n%s\n--- got:\n%s\ninput:\n%s", i+1, expected, got, input.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
