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

type edge struct {
	u, v int
	w    int64
}

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleD")
	cmd := exec.Command("go", "build", "-o", oracle, "536D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v\n%s", err, string(out))
	}
	return oracle, nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func buildInput(n int, s, t int, values []int64, edges []edge, queries []int64) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, len(edges))
	fmt.Fprintf(&sb, "%d %d\n", s, t)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", values[i])
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d %d\n", e.u, e.v, e.w)
	}
	fmt.Fprintf(&sb, "%d\n", len(queries))
	for i, q := range queries {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", q)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func deterministicTests() []string {
	var tests []string
	// Small simple graph
	tests = append(tests, buildInput(
		2, 1, 2,
		[]int64{5, -3},
		[]edge{
			{1, 2, 1},
		},
		[]int64{1, 5},
	))
	// Triangle with different weights
	tests = append(tests, buildInput(
		3, 1, 3,
		[]int64{10, 20, 30},
		[]edge{
			{1, 2, 2},
			{2, 3, 2},
			{1, 3, 5},
		},
		[]int64{0, 3, 10},
	))
	// Graph with self-loop and zero edge cost
	tests = append(tests, buildInput(
		4, 2, 4,
		[]int64{1, 2, 3, 4},
		[]edge{
			{1, 1, 0},
			{1, 2, 1},
			{2, 3, 2},
			{3, 4, 3},
			{4, 2, 1},
		},
		[]int64{1, 2, 4, 8},
	))
	return tests
}

func randInt(rnd *rand.Rand, l, r int) int {
	return rnd.Intn(r-l+1) + l
}

func randomValues(n int, rnd *rand.Rand) []int64 {
	vals := make([]int64, n)
	for i := range vals {
		vals[i] = rnd.Int63n(2_000_000_001) - 1_000_000_000
	}
	return vals
}

func randomGraph(n int, rnd *rand.Rand) []edge {
	maxEdges := 100000
	edges := make([]edge, 0, n-1)
	// build tree to ensure connectivity
	for i := 2; i <= n; i++ {
		parent := randInt(rnd, 1, i-1)
		w := rnd.Int63n(1_000_000_000 + 1)
		edges = append(edges, edge{parent, i, w})
	}
	remaining := maxEdges - len(edges)
	if remaining < 0 {
		remaining = 0
	}
	additionalLimit := min(remaining, n*4)
	extraCount := 0
	if additionalLimit > 0 {
		extraCount = rnd.Intn(additionalLimit + 1)
	}
	for i := 0; i < extraCount; i++ {
		if rnd.Intn(10) == 0 {
			u := randInt(rnd, 1, n)
			w := rnd.Int63n(1_000_000_000 + 1)
			edges = append(edges, edge{u, u, w})
			continue
		}
		u := randInt(rnd, 1, n)
		v := randInt(rnd, 1, n)
		w := rnd.Int63n(1_000_000_000 + 1)
		edges = append(edges, edge{u, v, w})
	}
	return edges
}

func randomQueries(rnd *rand.Rand) []int64 {
	q := randInt(rnd, 1, 50)
	qs := make([]int64, q)
	for i := range qs {
		qs[i] = rnd.Int63n(1_000_000_000_000) + 1
	}
	return qs
}

func randomTests(count int) []string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]string, 0, count)
	for i := 0; i < count; i++ {
		var n int
		switch {
		case i%25 == 0:
			n = randInt(rnd, 1000, 2000)
		case i%5 == 0:
			n = randInt(rnd, 100, 500)
		default:
			n = randInt(rnd, 2, 80)
		}
		values := randomValues(n, rnd)
		edges := randomGraph(n, rnd)
		qs := randomQueries(rnd)
		s := randInt(rnd, 1, n)
		t := randInt(rnd, 1, n)
		for t == s {
			t = randInt(rnd, 1, n)
		}
		tests = append(tests, buildInput(n, s, t, values, edges, qs))
	}
	// Add one near-max test
	n := 2000
	values := randomValues(n, rnd)
	edges := randomGraph(n, rnd)
	qs := randomQueries(rnd)
	s := 1
	t := n
	tests = append(tests, buildInput(n, s, t, values, edges, qs))
	return tests
}

func normalizeOutput(out string) (string, error) {
	out = strings.TrimSpace(out)
	if out == "" {
		return "", fmt.Errorf("empty output")
	}
	switch out {
	case "Break a heart", "Cry", "Flowers":
		return out, nil
	default:
		return "", fmt.Errorf("invalid outcome %q", out)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	tests := deterministicTests()
	tests = append(tests, randomTests(150)...)

	for idx, input := range tests {
		expOut, err := runProgram(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle error: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotOut, err := runProgram(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		exp, err := normalizeOutput(expOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle output invalid: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := normalizeOutput(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if exp != got {
			fmt.Fprintf(os.Stderr, "case %d mismatch: expected %q got %q\n", idx+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
