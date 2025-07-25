package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Test struct {
	n     int
	edges [][2]int
}

func genTree(n int) [][2]int {
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rand.Intn(i-1) + 1
		edges = append(edges, [2]int{p, i})
	}
	return edges
}

func generateTests() []Test {
	rand.Seed(43)
	tests := make([]Test, 0, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(8) + 2
		edges := genTree(n)
		tests = append(tests, Test{n: n, edges: edges})
	}
	// edge small trees
	tests = append(tests, Test{n: 2, edges: [][2]int{{1, 2}}})
	tests = append(tests, Test{n: 3, edges: [][2]int{{1, 2}, {1, 3}}})
	return tests
}

func solve(t Test) int64 {
	adj := make([][]int, t.n+1)
	for _, e := range t.edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	color := make([]int, t.n+1)
	for i := range color {
		color[i] = -1
	}
	q := []int{1}
	color[1] = 0
	for head := 0; head < len(q); head++ {
		u := q[head]
		for _, v := range adj[u] {
			if color[v] == -1 {
				color[v] = color[u] ^ 1
				q = append(q, v)
			}
		}
	}
	var c0, c1 int64
	for i := 1; i <= t.n; i++ {
		if color[i] == 0 {
			c0++
		} else {
			c1++
		}
	}
	return c0*c1 - int64(t.n-1)
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: verifierB <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	passed := 0
	for i, t := range tests {
		input := fmt.Sprintf("%d\n", t.n)
		for _, e := range t.edges {
			input += fmt.Sprintf("%d %d\n", e[0], e[1])
		}
		want := solve(t)
		output, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("Test %d: exec error: %v\n", i+1, err)
			continue
		}
		outStr := strings.TrimSpace(output)
		got, err := strconv.ParseInt(outStr, 10, 64)
		if err != nil || got != want {
			fmt.Printf("Test %d: expected %d got %s\n", i+1, want, outStr)
		} else {
			passed++
		}
	}
	fmt.Printf("Passed %d/%d tests\n", passed, len(tests))
}
