package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const refSourceE = "2000-2999/2000-2099/2060-2069/2068/2068E.go"

func main() {
	var candidate string
	if len(os.Args) == 2 {
		candidate = os.Args[1]
	} else if len(os.Args) == 3 && os.Args[1] == "--" {
		candidate = os.Args[2]
	} else {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/candidate")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	for idx, input := range tests {
		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		refVal, err := parseSingleInt(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s\nstdout/stderr:\n%s\n", idx+1, err, input, candOut)
			os.Exit(1)
		}
		candVal, err := parseSingleInt(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid candidate output on test %d: %v\noutput:\n%s\n", idx+1, err, candOut)
			os.Exit(1)
		}

		if refVal != candVal {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%sreference: %d\ncandidate: %d\n", idx+1, input, refVal, candVal)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2068E-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSourceE))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseSingleInt(out string) (int64, error) {
	fields := strings.Fields(out)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected 1 integer, got %d tokens", len(fields))
	}
	return strconv.ParseInt(fields[0], 10, 64)
}

func buildTests() []string {
	var tests []string

	// Samples from statement.
	tests = append(tests, "5 5\n1 2\n1 3\n2 5\n3 4\n4 5\n")
	tests = append(tests, "11 12\n1 2\n2 3\n3 4\n4 5\n5 11\n3 6\n6 7\n7 11\n1 8\n8 9\n9 10\n10 11\n")
	tests = append(tests, "8 10\n1 2\n2 3\n3 4\n3 8\n4 8\n1 5\n5 6\n6 7\n6 8\n7 8\n")

	// Minimal graph where police can trap.
	tests = append(tests, "2 1\n1 2\n")
	tests = append(tests, "3 2\n1 2\n2 3\n") // bridge to destination -> -1

	// Two parallel disjoint 1- n paths of different lengths.
	tests = append(tests, formatGraph(6, [][2]int{
		{1, 2}, {2, 6}, // short path
		{1, 3}, {3, 4}, {4, 5}, {5, 6}, // longer path
	}))

	rng := rand.New(rand.NewSource(2068))

	// Random medium graphs.
	for i := 0; i < 5; i++ {
		n := rng.Intn(25) + 6
		m := n - 1 + rng.Intn(n)
		tests = append(tests, randomGraphInput(rng, n, m))
	}

	// Dense-ish random graph.
	tests = append(tests, randomGraphInput(rng, 120, 250))

	// Stress graph near limits but still reasonable for a verifier run.
	tests = append(tests, randomGraphInput(rng, 3000, 6000))

	return tests
}

func formatGraph(n int, edges [][2]int) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d\n", n, len(edges))
	for _, e := range edges {
		fmt.Fprintf(&b, "%d %d\n", e[0], e[1])
	}
	return b.String()
}

func randomGraphInput(rng *rand.Rand, n, m int) string {
	if m < n-1 {
		m = n - 1
	}
	// Build a random tree first to ensure connectivity.
	type edge struct{ u, v int }
	edges := make([]edge, 0, m)
	for v := 2; v <= n; v++ {
		u := rng.Intn(v-1) + 1
		edges = append(edges, edge{u: u, v: v})
	}
	exists := make(map[[2]int]struct{}, m*2)
	for _, e := range edges {
		a, b := e.u, e.v
		if a > b {
			a, b = b, a
		}
		exists[[2]int{a, b}] = struct{}{}
	}
	// Add extra random edges.
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n-1) + 1
		if v >= u {
			v++
		}
		a, b := u, v
		if a > b {
			a, b = b, a
		}
		key := [2]int{a, b}
		if _, ok := exists[key]; ok {
			continue
		}
		exists[key] = struct{}{}
		edges = append(edges, edge{u: u, v: v})
	}
	// Shuffle edges for variability.
	for i := len(edges) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		edges[i], edges[j] = edges[j], edges[i]
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d\n", n, len(edges))
	for _, e := range edges {
		fmt.Fprintf(&b, "%d %d\n", e.u, e.v)
	}
	return b.String()
}
