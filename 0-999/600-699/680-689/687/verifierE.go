package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleE")
	cmd := exec.Command("go", "build", "-o", oracle, "687E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// genCase generates a random directed graph with n vertices and up to maxEdges edges.
// The graph has no self-loops and at most one edge between any unordered pair.
func genCase(rng *rand.Rand, n, maxEdges int) string {
	type edge struct{ u, v int }
	edgeSet := make(map[[2]int]bool)
	var edges []edge
	attempts := maxEdges * 4
	for len(edges) < maxEdges && attempts > 0 {
		attempts--
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		key := [2]int{u, v}
		if edgeSet[key] {
			continue
		}
		edgeSet[key] = true
		edges = append(edges, edge{u, v})
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, len(edges))
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e.u, e.v)
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Fixed cases from the problem (if any known samples existed).
	// n=1,m=0 edge case.
	fixed := []string{
		"1 0\n",
		"2 1\n1 2\n",
		"2 2\n1 2\n2 1\n",
		"3 3\n1 2\n2 3\n3 1\n",
		"4 4\n1 2\n2 3\n3 4\n4 1\n",
	}

	idx := 0
	for _, input := range fixed {
		idx++
		exp, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on fixed test %d: %v\n", idx, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fixed test %d failed: %v\ninput:\n%s", idx, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("fixed test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx, input, exp, got)
			os.Exit(1)
		}
	}

	// Random small tests.
	for idx < 300 {
		idx++
		n := rng.Intn(8) + 1
		maxEdges := rng.Intn(n*(n-1)/2+1) + 1
		if n == 1 {
			maxEdges = 0
		}
		input := genCase(rng, n, maxEdges)
		exp, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on test %d: %v\ninput:\n%s", idx, err, input)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%s", idx, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx, input, exp, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", idx)
}
