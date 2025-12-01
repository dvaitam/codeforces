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

const refSource = "./1639A.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for i, input := range tests {
		wantOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		gotOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		wantTokens := tokenizeOutput(wantOut)
		gotTokens := tokenizeOutput(gotOut)
		if len(wantTokens) != len(gotTokens) {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %d tokens, got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				i+1, len(wantTokens), len(gotTokens), input, wantOut, gotOut)
			os.Exit(1)
		}
		for j := range wantTokens {
			if wantTokens[j] != gotTokens[j] {
				fmt.Fprintf(os.Stderr, "test %d failed at token %d: expected %q got %q\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
					i+1, j+1, wantTokens[j], gotTokens[j], input, wantOut, gotOut)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1639A-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("go build failed: %v\n%s", err, stderr.String())
	}
	return tmp.Name(), nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		abs, err := filepath.Abs(bin)
		if err != nil {
			return "", err
		}
		cmd = exec.Command("go", "run", abs)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func tokenizeOutput(out string) []string {
	return strings.Fields(strings.ReplaceAll(out, "\r", ""))
}

type edge struct{ u, v int }
type query struct {
	d     int
	deg   []int
	flags []int
}

type testCase struct {
	n, m, start, base int
	edges             []edge
	queries           []query
	finalToken        string
}

func generateTests() []string {
	var tests []string
	tests = append(tests, buildTest([]testCase{{
		n: 1, m: 0, start: 1, base: 0, finalToken: "AC",
	}}))

	tests = append(tests, buildTest([]testCase{{
		n:     3,
		m:     2,
		start: 2,
		base:  5,
		edges: []edge{{1, 2}, {2, 3}},
		queries: []query{
			{d: 3, deg: []int{2, 3, 1}, flags: []int{1, 0, 1}},
			{d: 2, deg: []int{2, 2}, flags: []int{1, 1}},
			{d: 4, deg: []int{3, 2, 2, 1}, flags: []int{0, 0, 1, 1}},
		},
		finalToken: "AC",
	}}))

	tests = append(tests, buildTest([]testCase{
		randomCase(5, 4, 3),
		randomCase(4, 3, 2),
	}))

	for i := 0; i < 3; i++ {
		tests = append(tests, buildTest([]testCase{randomCase(8+i, 5+i, 4)}))
	}

	return tests
}

func buildTest(cases []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, cs := range cases {
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", cs.n, cs.m, cs.start, cs.base))
		for _, e := range cs.edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
		}
		for _, q := range cs.queries {
			sb.WriteString("R\n")
			sb.WriteString(fmt.Sprintf("%d\n", q.d))
			for i := 0; i < q.d; i++ {
				sb.WriteString(fmt.Sprintf("%d %d\n", q.deg[i], q.flags[i]))
			}
		}
		sb.WriteString(cs.finalToken)
		sb.WriteByte('\n')
	}
	return sb.String()
}

var globalRNG = rand.New(rand.NewSource(1337))

func randomCase(maxN, maxM, maxQueries int) testCase {
	rng := globalRNG
	n := rng.Intn(maxN-1) + 2
	if n > 50 {
		n = 50
	}
	m := rng.Intn(maxM) + 1
	if m < n-1 {
		m = n - 1
	}
	edges := make([]edge, 0, m)
	for i := 1; i < n; i++ {
		edges = append(edges, edge{i, i + 1})
	}
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		edges = append(edges, edge{u, v})
	}
	queryCount := rng.Intn(maxQueries) + 1
	queries := make([]query, 0, queryCount)
	for q := 0; q < queryCount; q++ {
		d := rng.Intn(4) + 1
		deg := make([]int, d)
		flags := make([]int, d)
		for i := 0; i < d; i++ {
			deg[i] = rng.Intn(4) + 1
			flags[i] = rng.Intn(2)
		}
		queries = append(queries, query{d: d, deg: deg, flags: flags})
	}
	final := "AC"
	if rng.Intn(5) == 0 {
		final = "F"
	}
	return testCase{
		n:          n,
		m:          len(edges),
		start:      rng.Intn(n) + 1,
		base:       rng.Intn(10) + 1,
		edges:      edges,
		queries:    queries,
		finalToken: final,
	}
}
