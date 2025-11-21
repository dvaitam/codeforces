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
	n       int
	edges   [][2]int
	queries [][2]int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refSrc, err := locateReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	refBin, err := buildReference(refSrc)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, input := range tests {
		want, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if normalize(got) != normalize(want) {
			fmt.Fprintf(os.Stderr, "mismatch on test %d\nInput:\n%s\nExpected:\n%s\nGot:\n%s\n", idx+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	candidates := []string{
		"97E.go",
		filepath.Join("0-999", "0-99", "90-99", "97", "97E.go"),
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("could not find 97E.go relative to the working directory")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref97E_%d.bin", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", outPath, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return outPath, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr:\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func normalize(out string) string {
	out = strings.TrimSpace(out)
	lines := strings.Split(out, "\n")
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	return strings.Join(lines, "\n")
}

func generateTests() []string {
	var tests []string
	tests = append(tests, formatTestCase(testCase{
		n:       1,
		queries: [][2]int{{1, 1}},
	}))
	tests = append(tests, formatTestCase(testCase{
		n:       2,
		edges:   [][2]int{{1, 2}},
		queries: [][2]int{{1, 2}, {2, 1}},
	}))
	tests = append(tests, formatTestCase(testCase{
		n:       4,
		edges:   [][2]int{{1, 2}, {2, 3}, {3, 4}},
		queries: [][2]int{{1, 3}, {2, 4}, {1, 4}},
	}))
	tests = append(tests, formatTestCase(testCase{
		n:       5,
		edges:   [][2]int{{1, 2}, {2, 3}, {3, 1}, {3, 4}, {4, 5}},
		queries: [][2]int{{1, 4}, {2, 5}, {4, 5}, {5, 1}},
	}))
	tests = append(tests, formatTestCase(testCase{
		n:       6,
		edges:   [][2]int{{1, 2}, {2, 3}, {4, 5}, {5, 6}, {6, 4}},
		queries: [][2]int{{1, 6}, {2, 3}, {4, 6}, {1, 1}, {4, 5}},
	}))
	tests = append(tests, formatTestCase(buildLineTest(50)))
	tests = append(tests, formatTestCase(completeGraphTest(4)))

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	target := 70
	for len(tests) < target {
		n := rng.Intn(30) + 1
		maxEdges := n * (n - 1) / 2
		m := 0
		if maxEdges > 0 {
			m = rng.Intn(maxEdges + 1)
		}
		edges := make([][2]int, 0, m)
		used := make(map[[2]int]struct{})
		for len(edges) < m {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			if u == v {
				continue
			}
			if u > v {
				u, v = v, u
			}
			key := [2]int{u, v}
			if _, ok := used[key]; ok {
				continue
			}
			used[key] = struct{}{}
			edges = append(edges, key)
		}
		q := rng.Intn(n*n + n + 1)
		if q == 0 {
			q = 1
		}
		queries := make([][2]int, q)
		for i := 0; i < q; i++ {
			queries[i] = [2]int{rng.Intn(n) + 1, rng.Intn(n) + 1}
		}
		tests = append(tests, formatTestCase(testCase{
			n:       n,
			edges:   edges,
			queries: queries,
		}))
	}
	return tests
}

func buildLineTest(n int) testCase {
	edges := make([][2]int, 0, n-1)
	for i := 1; i < n; i++ {
		edges = append(edges, [2]int{i, i + 1})
	}
	queries := make([][2]int, 0, n+1)
	for i := 1; i <= n; i++ {
		queries = append(queries, [2]int{1, i})
	}
	queries = append(queries, [2]int{n, n})
	return testCase{n: n, edges: edges, queries: queries}
}

func completeGraphTest(n int) testCase {
	var edges [][2]int
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			edges = append(edges, [2]int{i, j})
		}
	}
	var queries [][2]int
	for i := 1; i <= n; i++ {
		for j := 1; j <= n; j++ {
			queries = append(queries, [2]int{i, j})
		}
	}
	return testCase{n: n, edges: edges, queries: queries}
}

func formatTestCase(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, len(tc.edges)))
	for _, e := range tc.edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.queries)))
	for _, q := range tc.queries {
		sb.WriteString(fmt.Sprintf("%d %d\n", q[0], q[1]))
	}
	return sb.String()
}
