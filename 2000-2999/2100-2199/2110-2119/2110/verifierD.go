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

const refSource = "2000-2999/2100-2199/2110-2119/2110/2110D.go"

type edge struct {
	s int
	t int
	w int64
}

type testCase struct {
	n     int
	b     []int64
	edges []edge
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[len(os.Args)-1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	input := renderInput(tests)

	refOut, err := runWithInput(exec.Command(refBin), input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference solution failed: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}

	candOut, err := runWithInput(commandFor(candidate), input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	expect, err := parseOutputs(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse reference output: %v\n", err)
		os.Exit(1)
	}
	got, err := parseOutputs(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse candidate output: %v\n", err)
		os.Exit(1)
	}

	for i := range tests {
		if expect[i] != got[i] {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d got %d\ninput:\n%s", i+1, expect[i], got[i], formatSingleInput(tests[i]))
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2110D-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
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

func parseOutputs(output string, expected int) ([]int64, error) {
	fields := strings.Fields(output)
	if len(fields) < expected {
		return nil, fmt.Errorf("expected %d numbers, got %d", expected, len(fields))
	}
	if len(fields) > expected {
		return nil, fmt.Errorf("extra output detected after %d numbers", expected)
	}
	res := make([]int64, expected)
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse number %d (%q): %v", i+1, f, err)
		}
		res[i] = val
	}
	return res, nil
}

func renderInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, len(tc.edges)))
		for i, v := range tc.b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		for _, e := range tc.edges {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", e.s, e.t, e.w))
		}
	}
	return sb.String()
}

func formatSingleInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, len(tc.edges)))
	for i, v := range tc.b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	for _, e := range tc.edges {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e.s, e.t, e.w))
	}
	return sb.String()
}

func generateTests() []testCase {
	tests := []testCase{
		// Samples
		{
			n: 3,
			b: []int64{2, 0, 0},
			edges: []edge{
				{1, 2, 1}, {2, 3, 1}, {1, 3, 2},
			},
		},
		{
			n: 5,
			b: []int64{2, 2, 5, 0, 1},
			edges: []edge{
				{1, 2, 2}, {1, 3, 1}, {1, 4, 3}, {3, 5, 5}, {2, 4, 4}, {4, 5, 3},
			},
		},
		{
			n: 3,
			b: []int64{1, 1, 0},
			edges: []edge{
				{1, 2, 1},
			},
		},
		{
			n: 4,
			b: []int64{1, 9, 0, 0},
			edges: []edge{
				{1, 2, 1}, {1, 3, 3}, {2, 3, 10}, {3, 4, 5},
			},
		},
		// Simple cases
		{
			n: 2,
			b: []int64{0, 0},
			edges: []edge{
				{1, 2, 1},
			},
		},
		{
			n: 2,
			b: []int64{0, 5},
			edges: []edge{
				{1, 2, 10},
			},
		},
		{
			n:     4,
			b:     []int64{0, 0, 0, 0},
			edges: []edge{},
		},
	}

	totalN := 0
	for _, tc := range tests {
		totalN += tc.n
	}

	const maxTotalN = 15000
	rng := rand.New(rand.NewSource(2110))

	for totalN < maxTotalN {
		n := rng.Intn(400) + 2
		if totalN+n > maxTotalN {
			n = maxTotalN - totalN
		}
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			b[i] = rng.Int63n(1_000_000_000 + 1)
		}

		// Ensure at least a backbone path.
		backbone := make([]edge, 0, n-1)
		for i := 1; i < n; i++ {
			backbone = append(backbone, edge{s: i, t: i + 1, w: rng.Int63n(5) + 1})
		}

		edgeSet := make(map[[2]int]struct{})
		for _, e := range backbone {
			edgeSet[[2]int{e.s, e.t}] = struct{}{}
		}

		maxExtra := n * (n - 1) / 2
		available := maxExtra - len(edgeSet)
		extraLimit := minInt(available, 6*n)
		desired := len(edgeSet)
		if extraLimit > 0 {
			desired += rng.Intn(extraLimit) + 1
		}
		for len(edgeSet) < desired {
			s := rng.Intn(n-1) + 1
			t := rng.Intn(n-s) + s + 1
			key := [2]int{s, t}
			if _, ok := edgeSet[key]; ok {
				continue
			}
			edgeSet[key] = struct{}{}
		}

		edges := make([]edge, 0, len(edgeSet))
		for key := range edgeSet {
			edges = append(edges, edge{s: key[0], t: key[1], w: rng.Int63n(1_000_000_000) + 1})
		}

		tests = append(tests, testCase{n: n, b: b, edges: edges})
		totalN += n
		if len(tests) > 300 {
			break
		}
	}

	return tests
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
