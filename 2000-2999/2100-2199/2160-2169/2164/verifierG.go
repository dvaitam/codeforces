package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type edge struct {
	u, v int
}

type testCase struct {
	n     int
	edges map[edge]struct{}
}

func canonicalEdge(u, v int) edge {
	if u > v {
		u, v = v, u
	}
	return edge{u: u, v: v}
}

func buildReferenceBinary() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to determine current file path")
	}
	dir := filepath.Dir(file)
	refPath := filepath.Join(dir, "2164G_ref.bin")

	cmd := exec.Command("go", "build", "-o", refPath, "2164G.go")
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_ = os.Remove(refPath)
		return "", fmt.Errorf("failed to build reference binary: %v\n%s", err, stderr.String())
	}
	return refPath, nil
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if stderr.Len() > 0 {
			return "", fmt.Errorf("%v: %s", err, strings.TrimSpace(stderr.String()))
		}
		return "", err
	}
	return stdout.String(), nil
}

func parseInput(input string) ([]testCase, error) {
	reader := strings.NewReader(input)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, fmt.Errorf("failed to read t: %v", err)
	}
	cases := make([]testCase, t)
	for i := 0; i < t; i++ {
		var n int
		if _, err := fmt.Fscan(reader, &n); err != nil {
			return nil, fmt.Errorf("test %d: failed to read n: %v", i+1, err)
		}
		if n < 3 {
			return nil, fmt.Errorf("test %d: invalid n %d", i+1, n)
		}
		edges := make(map[edge]struct{}, n-1)
		for j := 0; j < n-1; j++ {
			var u, v int
			if _, err := fmt.Fscan(reader, &u, &v); err != nil {
				return nil, fmt.Errorf("test %d: failed to read edge %d: %v", i+1, j+1, err)
			}
			if u < 1 || u > n || v < 1 || v > n || u == v {
				return nil, fmt.Errorf("test %d: invalid edge %d %d", i+1, u, v)
			}
			key := canonicalEdge(u, v)
			if _, exists := edges[key]; exists {
				return nil, fmt.Errorf("test %d: duplicate edge %d %d", i+1, u, v)
			}
			edges[key] = struct{}{}
		}
		cases[i] = testCase{n: n, edges: edges}
	}
	return cases, nil
}

func validateOutput(cases []testCase, out string) error {
	reader := strings.NewReader(out)
	for idx, tc := range cases {
		if tc.n == 0 {
			continue
		}
		seen := make(map[edge]struct{}, len(tc.edges))
		for e := 0; e < tc.n-1; e++ {
			var u, v int
			if _, err := fmt.Fscan(reader, &u, &v); err != nil {
				if err == io.EOF {
					return fmt.Errorf("test %d: expected %d edges, got %d", idx+1, tc.n-1, e)
				}
				return fmt.Errorf("test %d: failed to read edge %d: %v", idx+1, e+1, err)
			}
			if u < 1 || u > tc.n || v < 1 || v > tc.n {
				return fmt.Errorf("test %d: edge %d %d out of range", idx+1, u, v)
			}
			if u == v {
				return fmt.Errorf("test %d: self-loop %d %d", idx+1, u, v)
			}
			key := canonicalEdge(u, v)
			if _, ok := tc.edges[key]; !ok {
				return fmt.Errorf("test %d: edge %d %d not in tree", idx+1, u, v)
			}
			if _, used := seen[key]; used {
				return fmt.Errorf("test %d: duplicate edge %d %d", idx+1, u, v)
			}
			seen[key] = struct{}{}
		}
		if len(seen) != len(tc.edges) {
			return fmt.Errorf("test %d: missing edges")
		}
	}
	var extra int
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return fmt.Errorf("extra output detected, starts with %d", extra)
	} else if err != io.EOF {
		return fmt.Errorf("failed when checking for extra output: %v", err)
	}
	return nil
}

type treeCase struct {
	n     int
	edges [][2]int
}

func buildInput(cases []treeCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(cases)))
	sb.WriteByte('\n')
	for _, tc := range cases {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for _, e := range tc.edges {
			sb.WriteString(strconv.Itoa(e[0]))
			sb.WriteByte(' ')
			sb.WriteString(strconv.Itoa(e[1]))
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func randomTree(n int, rng *rand.Rand) [][2]int {
	edges := make([][2]int, n-1)
	for v := 2; v <= n; v++ {
		u := rng.Intn(v-1) + 1
		edges[v-2] = [2]int{u, v}
	}
	for i := range edges {
		j := rng.Intn(i + 1)
		edges[i], edges[j] = edges[j], edges[i]
	}
	return edges
}

func chainTree(n int) [][2]int {
	edges := make([][2]int, n-1)
	for i := 2; i <= n; i++ {
		edges[i-2] = [2]int{i - 1, i}
	}
	return edges
}

func starTree(n int) [][2]int {
	edges := make([][2]int, n-1)
	for i := 2; i <= n; i++ {
		edges[i-2] = [2]int{1, i}
	}
	return edges
}

func generateTests() []string {
	rng := rand.New(rand.NewSource(2024))
	var tests []string

	tests = append(tests, buildInput([]treeCase{
		{n: 3, edges: chainTree(3)},
	}))

	tests = append(tests, buildInput([]treeCase{
		{n: 4, edges: starTree(4)},
		{n: 5, edges: chainTree(5)},
	}))

	var randomCases []treeCase
	for i := 0; i < 5; i++ {
		n := rng.Intn(20) + 3
		randomCases = append(randomCases, treeCase{n: n, edges: randomTree(n, rng)})
	}
	tests = append(tests, buildInput(randomCases))

	randomCases = randomCases[:0]
	for i := 0; i < 10; i++ {
		n := rng.Intn(200) + 50
		randomCases = append(randomCases, treeCase{n: n, edges: randomTree(n, rng)})
	}
	tests = append(tests, buildInput(randomCases))

	tests = append(tests, buildInput([]treeCase{
		{n: 50000, edges: randomTree(50000, rng)},
	}))

	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	refPath, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refPath)

	tests := generateTests()
	for idx, input := range tests {
		cases, err := parseInput(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal error parsing test %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		refOut, err := runProgram(refPath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if err := validateOutput(cases, refOut); err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		userOut, err := runProgram(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		if err := validateOutput(cases, userOut); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%s\noutput:\n%s\n", idx+1, err, input, userOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
