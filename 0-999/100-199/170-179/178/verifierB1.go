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

const refSource = "./178B1.go"

type testCase struct {
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB1.go /path/to/binary")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	candidate := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		expect, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		got, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, err, tc.input, got)
			os.Exit(1)
		}
		if !equalTokens(expect, got) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, tc.input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "178B1-ref-*")
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
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func equalTokens(a, b string) bool {
	ta := strings.Fields(a)
	tb := strings.Fields(b)
	if len(ta) != len(tb) {
		return false
	}
	for i := range ta {
		if ta[i] != tb[i] {
			return false
		}
	}
	return true
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(17801781))
	var tests []testCase

	tests = append(tests, manualCase(
		5,
		[][2]int{{1, 2}, {2, 3}, {3, 4}, {4, 5}, {2, 5}},
		[][2]int{{1, 5}, {1, 3}, {4, 4}},
	))

	tests = append(tests, manualCase(
		6,
		[][2]int{{1, 2}, {2, 3}, {3, 1}, {3, 4}, {4, 5}, {5, 6}},
		[][2]int{{1, 2}, {1, 6}, {4, 6}, {2, 5}},
	))

	tests = append(tests, manualCase(
		4,
		[][2]int{{1, 2}, {2, 3}, {3, 4}, {4, 1}},
		[][2]int{{1, 3}, {2, 4}, {1, 1}},
	))

	tests = append(tests, denseTest(12))
	tests = append(tests, lineTest(15))

	for i := 0; i < 25; i++ {
		n := rng.Intn(40) + 2
		tests = append(tests, randomTest(rng, n))
	}

	for i := 0; i < 10; i++ {
		n := rng.Intn(60) + 40
		tests = append(tests, randomTest(rng, n))
	}

	tests = append(tests, limitLikeTest())
	return tests
}

func manualCase(n int, edges [][2]int, queries [][2]int) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, len(edges))
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	fmt.Fprintf(&sb, "%d\n", len(queries))
	for _, q := range queries {
		fmt.Fprintf(&sb, "%d %d\n", q[0], q[1])
	}
	return testCase{input: sb.String()}
}

func denseTest(n int) testCase {
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
			if len(queries) >= 3*n {
				break
			}
		}
		if len(queries) >= 3*n {
			break
		}
	}
	return manualCase(n, edges, queries)
}

func lineTest(n int) testCase {
	var edges [][2]int
	for i := 1; i < n; i++ {
		edges = append(edges, [2]int{i, i + 1})
	}
	var queries [][2]int
	for i := 1; i <= n; i++ {
		queries = append(queries, [2]int{1, i})
	}
	return manualCase(n, edges, queries)
}

func randomTest(rng *rand.Rand, n int) testCase {
	maxEdges := n * (n - 1) / 2
	m := n - 1
	if extra := maxEdges - (n - 1); extra > 0 {
		m += rng.Intn(extra + 1)
	}
	edges := make([][2]int, 0, m)
	seen := make(map[[2]int]struct{}, m)
	for v := 2; v <= n; v++ {
		p := rng.Intn(v-1) + 1
		addEdge(p, v, &edges, seen)
	}
	for len(edges) < m {
		a := rng.Intn(n) + 1
		b := rng.Intn(n) + 1
		if a == b {
			continue
		}
		addEdge(a, b, &edges, seen)
	}
	k := rng.Intn(n) + 1
	var queries [][2]int
	for i := 0; i < k; i++ {
		s := rng.Intn(n) + 1
		l := rng.Intn(n) + 1
		queries = append(queries, [2]int{s, l})
	}
	return manualCase(n, edges, queries)
}

func addEdge(a, b int, edges *[][2]int, seen map[[2]int]struct{}) {
	if a > b {
		a, b = b, a
	}
	key := [2]int{a, b}
	if _, ok := seen[key]; ok {
		return
	}
	seen[key] = struct{}{}
	*edges = append(*edges, [2]int{a, b})
}

func limitLikeTest() testCase {
	n := 100
	var edges [][2]int
	for i := 1; i < n; i++ {
		edges = append(edges, [2]int{i, i + 1})
	}
	for i := 1; i+10 <= n; i += 10 {
		edges = append(edges, [2]int{i, i + 10})
	}
	var queries [][2]int
	for i := 1; i <= n; i++ {
		queries = append(queries, [2]int{i, n + 1 - i})
	}
	return manualCase(n, edges, queries)
}
