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

type testCase struct {
	n, m, k int
	types   []int
	edges   [][2]int
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2172M-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), "2172M.go")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func commandFor(target string) *exec.Cmd {
	switch filepath.Ext(target) {
	case ".go":
		return exec.Command("go", "run", target)
	case ".py":
		return exec.Command("python3", target)
	default:
		return exec.Command(target)
	}
}

func runProgram(target, input string) (string, error) {
	cmd := commandFor(target)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.k))
	for i, v := range tc.types {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for _, e := range tc.edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return sb.String()
}

func addTest(tests *[]string, tc testCase) {
	*tests = append(*tests, buildInput(tc))
}

func chainGraph(n int) testCase {
	types := make([]int, n)
	for i := 0; i < n; i++ {
		types[i] = i%n + 1
	}
	edges := make([][2]int, n-1)
	for i := 1; i < n; i++ {
		edges[i-1] = [2]int{i, i + 1}
	}
	return testCase{
		n:     n,
		m:     len(edges),
		k:     n,
		types: types,
		edges: edges,
	}
}

func starGraph(n int) testCase {
	types := make([]int, n)
	for i := 0; i < n; i++ {
		types[i] = 1
	}
	types[0] = 2
	edges := make([][2]int, n-1)
	for i := 2; i <= n; i++ {
		edges[i-2] = [2]int{1, i}
	}
	return testCase{
		n:     n,
		m:     len(edges),
		k:     2,
		types: types,
		edges: edges,
	}
}

func randomGraph(rng *rand.Rand, n, extraEdges int) testCase {
	if extraEdges < n-1 {
		extraEdges = n - 1
	}
	edges := make([][2]int, 0, extraEdges)
	for v := 2; v <= n; v++ {
		u := rng.Intn(v-1) + 1
		edges = append(edges, [2]int{u, v})
	}
	added := make(map[[2]int]struct{})
	for _, e := range edges {
		a, b := e[0], e[1]
		if a > b {
			a, b = b, a
		}
		added[[2]int{a, b}] = struct{}{}
	}
	for len(edges) < extraEdges {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		a, b := u, v
		if a > b {
			a, b = b, a
		}
		key := [2]int{a, b}
		if _, ok := added[key]; ok {
			continue
		}
		added[key] = struct{}{}
		edges = append(edges, [2]int{u, v})
	}
	k := rng.Intn(n) + 1
	types := make([]int, n)
	for i := range types {
		types[i] = rng.Intn(k) + 1
	}
	for missing := 1; missing <= k; missing++ {
		found := false
		for _, t := range types {
			if t == missing {
				found = true
				break
			}
		}
		if !found {
			types[rng.Intn(n)] = missing
		}
	}
	return testCase{
		n:     n,
		m:     len(edges),
		k:     k,
		types: types,
		edges: edges,
	}
}

func buildTests() []string {
	var tests []string
	addTest(&tests, testCase{
		n:     3,
		m:     3,
		k:     2,
		types: []int{2, 1, 1},
		edges: [][2]int{{1, 2}, {3, 1}, {3, 2}},
	})
	addTest(&tests, chainGraph(5))
	addTest(&tests, starGraph(6))

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 60; i++ {
		n := rng.Intn(15) + 3
		m := rng.Intn(n*(n-1)/2-(n-1)+1) + (n - 1)
		addTest(&tests, randomGraph(rng, n, m))
	}
	for i := 0; i < 20; i++ {
		n := rng.Intn(200) + 50
		m := n - 1 + rng.Intn(n)
		addTest(&tests, randomGraph(rng, n, m))
	}
	tests = append(tests, buildInput(randomGraph(rng, 200000, 200000)))
	return tests
}

func compareOutputs(expected, got string) error {
	expFields := strings.Fields(expected)
	gotFields := strings.Fields(got)
	if len(expFields) != len(gotFields) {
		return fmt.Errorf("expected %d outputs, got %d", len(expFields), len(gotFields))
	}
	for i := range expFields {
		if expFields[i] != gotFields[i] {
			return fmt.Errorf("mismatch at position %d: expected %s got %s", i+1, expFields[i], gotFields[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierM.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	for idx, input := range tests {
		expect, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		if err := compareOutputs(expect, got); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, err, input, expect, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}
