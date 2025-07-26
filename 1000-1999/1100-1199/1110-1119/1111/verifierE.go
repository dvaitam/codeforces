package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type edge struct{ u, v int }

type query struct {
	k     int
	m     int
	r     int
	nodes []int
}

type test struct {
	n       int
	edges   []edge
	queries []query
}

func buildReference() (string, error) {
	exe := filepath.Join(os.TempDir(), "refE_bin")
	cmd := exec.Command("go", "build", "-o", exe, "1111E.go")
	cmd.Dir = "."
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return exe, nil
}

func run(binary string, input string) (string, error) {
	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(46))
	tests := make([]test, 0, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(5) + 2
		edges := make([]edge, n-1)
		for j := 1; j < n; j++ {
			u := rng.Intn(j) + 1
			edges[j-1] = edge{u, j + 1}
		}
		qn := rng.Intn(2) + 1
		queries := make([]query, qn)
		for j := 0; j < qn; j++ {
			k := rng.Intn(n) + 1
			m := rng.Intn(k) + 1
			r := rng.Intn(n) + 1
			nodes := make([]int, k)
			used := map[int]bool{}
			for x := 0; x < k; x++ {
				v := rng.Intn(n) + 1
				for used[v] {
					v = rng.Intn(n) + 1
				}
				used[v] = true
				nodes[x] = v
			}
			queries[j] = query{k, m, r, nodes}
		}
		tests = append(tests, test{n, edges, queries})
	}
	return tests
}

func formatInput(t test) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", t.n, len(t.queries))
	for _, e := range t.edges {
		fmt.Fprintf(&sb, "%d %d\n", e.u, e.v)
	}
	for _, q := range t.queries {
		fmt.Fprintf(&sb, "%d %d %d", q.k, q.m, q.r)
		for _, v := range q.nodes {
			fmt.Fprintf(&sb, " %d", v)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	ref, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference solution:", err)
		os.Exit(1)
	}
	tests := generateTests()
	for i, t := range tests {
		input := formatInput(t)
		want, err := run(ref, input)
		if err != nil {
			fmt.Printf("reference error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := run(binary, input)
		if err != nil {
			fmt.Printf("Test %d: error running binary: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected: %s\nGot: %s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
