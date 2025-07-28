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

type Edge struct {
	u int
	v int
}

type TestCase struct {
	n     int
	edges []Edge
}

func genTree(n int) []Edge {
	edges := make([]Edge, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rand.Intn(i-1) + 1
		edges = append(edges, Edge{p, i})
	}
	return edges
}

func genTests() []TestCase {
	rand.Seed(time.Now().UnixNano())
	const T = 100
	tests := make([]TestCase, T)
	for i := 0; i < T; i++ {
		n := rand.Intn(4) + 5 // 5..8
		edges := genTree(n)
		tests[i] = TestCase{n: n, edges: edges}
	}
	return tests
}

func buildInput(tests []TestCase) []byte {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, len(tests))
	for _, tc := range tests {
		fmt.Fprintln(&buf, tc.n)
		for _, e := range tc.edges {
			fmt.Fprintf(&buf, "%d %d\n", e.u, e.v)
		}
	}
	return buf.Bytes()
}

func run(binary string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(binary, ".go") {
		cmd = exec.Command("go", "run", binary)
	} else {
		cmd = exec.Command(binary)
	}
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	tests := genTests()
	input := buildInput(tests)

	// compile reference solution
	refBin := "ref_solF"
	cmd := exec.Command("g++", "solF.cpp", "-O2", "-std=c++17", "-o", refBin)
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to compile reference: %v\n%s\n", err, string(out))
		os.Exit(1)
	}
	defer os.Remove(refBin)

	expected, err := run("./"+refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference execution error: %v\n", err)
		os.Exit(1)
	}

	actual, err := run(binary, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error executing binary: %v\n", err)
		os.Exit(1)
	}

	if strings.TrimSpace(actual) != strings.TrimSpace(expected) {
		fmt.Fprintln(os.Stderr, "output mismatch between reference and binary")
		os.Exit(1)
	}

	fmt.Println("all tests passed")
}
