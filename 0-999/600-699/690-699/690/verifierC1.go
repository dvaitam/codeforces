package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type testCase struct {
	id    string
	input string
}

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC1.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[len(os.Args)-1]
	if candidate == "--" {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC1.go /path/to/candidate")
		os.Exit(1)
	}

	baseDir := currentDir()
	refBin, err := buildReference(baseDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for i, tc := range tests {
		expOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on %s: %v\ninput:\n%s", tc.id, err, tc.input)
			os.Exit(1)
		}
		gotOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on %s: %v\ninput:\n%s", tc.id, err, tc.input)
			os.Exit(1)
		}
		expVal, err := parseAnswer(expOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on %s: %v\noutput:\n%s", tc.id, err, expOut)
			os.Exit(1)
		}
		gotVal, err := parseAnswer(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on %s: %v\noutput:\n%s", tc.id, err, gotOut)
			os.Exit(1)
		}
		if expVal != gotVal {
			fmt.Fprintf(os.Stderr, "wrong answer on %s\nInput:\n%sExpected: %s\nGot: %s\n", tc.id, tc.input, formatBool(expVal), formatBool(gotVal))
			os.Exit(1)
		}
		if (i+1)%10 == 0 {
			fmt.Fprintf(os.Stderr, "validated %d/%d tests...\n", i+1, len(tests))
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func currentDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("cannot determine current file path")
	}
	return filepath.Dir(file)
}

func buildReference(dir string) (string, error) {
	out := filepath.Join(dir, "ref690C1.bin")
	cmd := exec.Command("go", "build", "-o", out, "690C1.go")
	cmd.Dir = dir
	if data, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("go build failed: %v\n%s", err, data)
	}
	return out, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseAnswer(out string) (bool, error) {
	reader := strings.NewReader(out)
	var token string
	if _, err := fmt.Fscan(reader, &token); err != nil {
		return false, fmt.Errorf("unable to read token: %v", err)
	}
	token = strings.ToLower(token)
	switch token {
	case "yes":
		return true, nil
	case "no":
		return false, nil
	default:
		return false, fmt.Errorf("unexpected token %q", token)
	}
}

func formatBool(v bool) string {
	if v {
		return "yes"
	}
	return "no"
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests,
		makeCase("single-node", 1, [][]int{}),
		makeCase("simple-edge", 2, [][]int{{1, 2}}),
		makeCase("disconnected", 3, [][]int{{1, 2}}),
		makeCase("extra-edge", 3, [][]int{{1, 2}, {2, 3}, {1, 3}}),
		makeCase("forest", 5, [][]int{{1, 2}, {2, 3}, {4, 5}}),
		makeCase("duplicate-edges", 4, [][]int{{1, 2}, {2, 3}, {3, 4}, {1, 2}}),
	)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 120; i++ {
		tests = append(tests, randomCase(rng, fmt.Sprintf("rand-%03d", i+1)))
	}
	return tests
}

func makeCase(id string, n int, edges [][]int) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, len(edges))
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	return testCase{id: id, input: sb.String()}
}

func randomCase(rng *rand.Rand, id string) testCase {
	n := rng.Intn(30) + 1
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(maxEdges + 1)
	edges := sampleEdges(rng, n, m)

	// occasionally force a tree or a disconnected graph
	switch rng.Intn(3) {
	case 0:
		edges = buildTreeEdges(rng, n)
	case 1:
		// ensure disconnected by splitting nodes into two sets and adding internal tree edges
		if n > 1 {
			cut := rng.Intn(n-1) + 1
			edges = append(buildTreeEdgesRange(rng, 1, cut), buildTreeEdgesRange(rng, cut+1, n)...)
		}
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, len(edges))
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	return testCase{id: id, input: sb.String()}
}

func sampleEdges(rng *rand.Rand, n, m int) [][]int {
	set := make(map[[2]int]struct{})
	edges := make([][]int, 0, m)
	for len(edges) < m {
		a := rng.Intn(n) + 1
		b := rng.Intn(n) + 1
		if a == b {
			continue
		}
		if a > b {
			a, b = b, a
		}
		key := [2]int{a, b}
		if _, ok := set[key]; ok {
			continue
		}
		set[key] = struct{}{}
		edges = append(edges, []int{a, b})
	}
	return edges
}

func buildTreeEdges(rng *rand.Rand, n int) [][]int {
	if n == 1 {
		return nil
	}
	edges := make([][]int, 0, n-1)
	for v := 2; v <= n; v++ {
		u := rng.Intn(v-1) + 1
		edges = append(edges, []int{u, v})
	}
	return edges
}

func buildTreeEdgesRange(rng *rand.Rand, l, r int) [][]int {
	if l >= r {
		return nil
	}
	edges := make([][]int, 0, r-l)
	for v := l + 1; v <= r; v++ {
		u := rng.Intn(v-l) + l
		edges = append(edges, []int{u, v})
	}
	return edges
}
