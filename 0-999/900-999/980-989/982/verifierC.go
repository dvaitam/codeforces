package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	input  string
	output string
}

func solve(n int, edges [][2]int) string {
	if n%2 == 1 {
		return "-1"
	}
	adj := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	parent := make([]int, n+1)
	order := []int{1}
	for i := 0; i < len(order); i++ {
		v := order[i]
		for _, to := range adj[v] {
			if to != parent[v] {
				parent[to] = v
				order = append(order, to)
			}
		}
	}
	size := make([]int, n+1)
	ans := 0
	for i := len(order) - 1; i >= 0; i-- {
		v := order[i]
		size[v] = 1
		for _, to := range adj[v] {
			if to != parent[v] {
				size[v] += size[to]
			}
		}
		if v != 1 && size[v]%2 == 0 {
			ans++
		}
	}
	return fmt.Sprintf("%d", ans)
}

func generateTests() []testCase {
	rand.Seed(3)
	var tests []testCase
	// simple case n=2
	tests = append(tests, testCase{input: "2\n1 2\n", output: "1"})
	// random cases
	for len(tests) < 120 {
		n := rand.Intn(10) + 1
		edges := make([][2]int, 0, n-1)
		for i := 2; i <= n; i++ {
			p := rand.Intn(i-1) + 1
			edges = append(edges, [2]int{p, i})
		}
		var b strings.Builder
		b.WriteString(fmt.Sprintf("%d\n", n))
		for _, e := range edges {
			fmt.Fprintf(&b, "%d %d\n", e[0], e[1])
		}
		tests = append(tests, testCase{input: b.String(), output: solve(n, edges)})
	}
	return tests
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		got, err := runBinary(binary, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != tc.output {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %q got %q\n", i+1, tc.output, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
