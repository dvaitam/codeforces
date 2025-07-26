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

type edge struct{ a, b int }

func dfsInfo(root int, adj [][]int) ([]int, []int) {
	n := len(adj) - 1
	parent := make([]int, n+1)
	size := make([]int, n+1)
	order := make([]int, 0, n)
	stack := []int{root}
	parent[root] = 0
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, v)
		for _, to := range adj[v] {
			if to == parent[v] {
				continue
			}
			parent[to] = v
			stack = append(stack, to)
		}
	}
	for i := len(order) - 1; i >= 0; i-- {
		v := order[i]
		size[v] = 1
		for _, to := range adj[v] {
			if to == parent[v] {
				continue
			}
			size[v] += size[to]
		}
	}
	return parent, size
}

func solve(n, x, y int, edges []edge) string {
	adj := make([][]int, n+1)
	for _, e := range edges {
		adj[e.a] = append(adj[e.a], e.b)
		adj[e.b] = append(adj[e.b], e.a)
	}
	parentX, sizeX := dfsInfo(x, adj)
	parentY, sizeY := dfsInfo(y, adj)
	xp := y
	for parentX[xp] != x {
		xp = parentX[xp]
	}
	yq := x
	for parentY[yq] != y {
		yq = parentY[yq]
	}
	sizeA := n - sizeX[xp]
	sizeB := n - sizeY[yq]
	totalPairs := int64(n) * int64(n-1)
	invalid := int64(sizeA) * int64(sizeB)
	ans := totalPairs - invalid
	return fmt.Sprintf("%d", ans)
}

func generateTree(rng *rand.Rand, n int) []edge {
	edges := make([]edge, 0, n-1)
	parents := make([]int, n+1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, edge{p, i})
		parents[i] = p
	}
	return edges
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(15) + 2
	edges := generateTree(rng, n)
	x := rng.Intn(n) + 1
	y := rng.Intn(n) + 1
	for y == x {
		y = rng.Intn(n) + 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, x, y)
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e.a, e.b)
	}
	input := sb.String()
	expected := solve(n, x, y, edges)
	return input, expected
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, in, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
