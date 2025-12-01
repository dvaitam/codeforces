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

type node struct {
	parent int
	typ    int
}

type query struct {
	t int
	u int
	v int
}

type testCase struct {
	n       int
	nodes   []node
	queries []query
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-855D-")
	if err != nil {
		return "", nil, err
	}
	path := filepath.Join(tmpDir, "oracleD")
	cmd := exec.Command("go", "build", "-o", path, "855D.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, string(out))
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return path, cleanup, nil
}

func runBinary(bin, input string) (string, error) {
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

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i := 0; i < tc.n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.nodes[i].parent, tc.nodes[i].typ))
	}
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.queries)))
	for _, q := range tc.queries {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", q.t, q.u, q.v))
	}
	return sb.String()
}

func normalizeAnswers(out string, expectedCount int) ([]string, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	filtered := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		line = strings.ToUpper(line)
		filtered = append(filtered, line)
	}
	if len(filtered) != expectedCount {
		return nil, fmt.Errorf("expected %d answers, got %d", expectedCount, len(filtered))
	}
	for i, v := range filtered {
		if v != "YES" && v != "NO" {
			return nil, fmt.Errorf("invalid answer at line %d: %q", i+1, v)
		}
	}
	return filtered, nil
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n: 1,
			nodes: []node{
				{parent: -1, typ: -1},
			},
			queries: []query{
				{t: 1, u: 1, v: 1},
			},
		},
		{
			n: 3,
			nodes: []node{
				{-1, -1},
				{1, 0},
				{2, 0},
			},
			queries: []query{
				{1, 1, 3},
				{2, 1, 3},
			},
		},
		{
			n: 4,
			nodes: []node{
				{-1, -1},
				{1, 0},
				{2, 1},
				{3, 0},
			},
			queries: []query{
				{1, 1, 4},
				{2, 1, 4},
				{2, 1, 3},
			},
		},
	}
}

func randomTest(rng *rand.Rand) testCase {
	n := rng.Intn(50) + 1
	nodes := make([]node, n)
	for i := 0; i < n; i++ {
		if i == 0 || rng.Intn(5) == 0 {
			nodes[i] = node{parent: -1, typ: -1}
			continue
		}
		parent := rng.Intn(i) + 1
		typ := rng.Intn(2)
		nodes[i] = node{parent: parent, typ: typ}
	}
	q := rng.Intn(200) + 1
	queries := make([]query, q)
	for i := 0; i < q; i++ {
		t := rng.Intn(2) + 1
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		queries[i] = query{t: t, u: u, v: v}
	}
	return testCase{n: n, nodes: nodes, queries: queries}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng))
	}

	for idx, tc := range tests {
		input := buildInput(tc)
		expOut, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		expAns, err := normalizeAnswers(expOut, len(tc.queries))
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid oracle output on test %d: %v\noutput:\n%s\n", idx+1, err, expOut)
			os.Exit(1)
		}

		gotOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		gotAns, err := normalizeAnswers(gotOut, len(tc.queries))
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\noutput:\n%s\ninput:\n%s\n", idx+1, err, gotOut, input)
			os.Exit(1)
		}

		for i := range expAns {
			if gotAns[i] != expAns[i] {
				fmt.Fprintf(os.Stderr, "mismatch on test %d query %d: expected %s got %s\ninput:\n%s\n", idx+1, i+1, expAns[i], gotAns[i], input)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
