package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	refSourcePath = "2000-2999/2000-2099/2000-2009/2000/2000G.go"
	totalTests    = 30
)

type testInput struct {
	name  string
	input string
}

type edge struct {
	u, v int
	l1   int64
	l2   int64
}

type caseData struct {
	n, m    int
	t0, t1  int64
	t2      int64
	edges   []edge
	checked bool
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary-or-source")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := generateTests()
	for i, tc := range tests {
		fmt.Fprintf(os.Stderr, "Test %d/%d: %s\n", i+1, len(tests), tc.name)
		expect, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s\n", i+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		continue
		got, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%s\n", i+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		if normalizeOutput(expect) != normalizeOutput(got) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s)\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
				i+1, tc.name, tc.input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	tmp, err := os.CreateTemp("", "ref2000G-*")
	if err != nil {
		return "", nil, fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmp.Name()
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmpPath, refSourcePath)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmpPath)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, string(out))
	}
	cleanup := func() { os.Remove(tmpPath) }
	return tmpPath, cleanup, nil
}

func runProgram(path string, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func runCandidate(target string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		abs, err := filepath.Abs(target)
		if err != nil {
			return "", fmt.Errorf("resolve candidate path: %w", err)
		}
		cmd = exec.Command("go", "run", abs)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func normalizeOutput(out string) string {
	fields := strings.Fields(out)
	return strings.Join(fields, "\n")
}

func generateTests() []testInput {
	var tests []testInput
	tests = append(tests,
		buildSingleCaseTest("simple_path", caseData{
			n: 2, m: 1,
			t0: 20, t1: 5, t2: 10,
			edges: []edge{{
				u: 1, v: 2, l1: 7, l2: 12,
			}},
		}),
		buildSingleCaseTest("impossible", caseData{
			n: 3, m: 2,
			t0: 5, t1: 1, t2: 2,
			edges: []edge{
				{u: 1, v: 2, l1: 10, l2: 12},
				{u: 2, v: 3, l1: 10, l2: 12},
			},
		}),
		buildMultiCaseTest("multi_case_mix", []caseData{
			{
				n:  3,
				t0: 50, t1: 5, t2: 15,
				edges: []edge{
					{u: 1, v: 2, l1: 10, l2: 20},
					{u: 2, v: 3, l1: 10, l2: 20},
				},
			},
			{
				n:  4,
				t0: 200, t1: 30, t2: 60,
				edges: []edge{
					{u: 1, v: 2, l1: 15, l2: 30},
					{u: 2, v: 3, l1: 15, l2: 40},
					{u: 3, v: 4, l1: 15, l2: 50},
					{u: 1, v: 4, l1: 45, l2: 90},
				},
			},
		}),
	)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < totalTests {
		tc := randomCase(rng, len(tests))
		tests = append(tests, tc)
	}

	return tests
}

func randomCase(rng *rand.Rand, idx int) testInput {
	n := rng.Intn(25) + 2
	maxExtra := n + rng.Intn(n)
	totalEdges := n - 1 + rng.Intn(maxExtra+1)
	edges := make([]edge, 0, totalEdges)
	used := make(map[[2]int]struct{})
	addEdge := func(u, v int, l1, l2 int64) {
		if u > v {
			u, v = v, u
		}
		key := [2]int{u, v}
		if _, ok := used[key]; ok {
			return
		}
		used[key] = struct{}{}
		edges = append(edges, edge{u: u, v: v, l1: l1, l2: l2})
	}
	for v := 2; v <= n; v++ {
		u := rng.Intn(v-1) + 1
		l1 := int64(rng.Intn(90) + 1)
		l2 := l1 + int64(rng.Intn(90)+1)
		addEdge(u, v, l1, l2)
	}
	for len(edges) < totalEdges {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		l1 := int64(rng.Intn(900) + 1)
		l2 := l1 + int64(rng.Intn(900)+1)
		addEdge(u, v, l1, l2)
	}

	t1 := int64(rng.Intn(500) + 1)
	t2 := t1 + int64(rng.Intn(500)+1)
	t0 := t2 + int64(rng.Intn(500)+rng.Intn(1000)+1)

	tc := caseData{
		n:     n,
		m:     len(edges),
		t0:    t0,
		t1:    t1,
		t2:    t2,
		edges: edges,
	}
	return buildSingleCaseTest(fmt.Sprintf("rand_%d", idx), tc)
}

func buildSingleCaseTest(name string, cd caseData) testInput {
	m := cd.m
	if m == 0 {
		m = len(cd.edges)
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", cd.n, m))
	sb.WriteString(fmt.Sprintf("%d %d %d\n", cd.t0, cd.t1, cd.t2))
	for _, e := range cd.edges {
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", e.u, e.v, e.l1, e.l2))
	}
	ti := testInput{
		name:  name,
		input: sb.String(),
	}
	if err := validateTestInput(ti.input); err != nil {
		panic(fmt.Sprintf("generated invalid test %s: %v", name, err))
	}
	return ti
}

func buildMultiCaseTest(name string, cases []caseData) testInput {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, cd := range cases {
		m := cd.m
		if m == 0 {
			m = len(cd.edges)
		}
		sb.WriteString(fmt.Sprintf("%d %d\n", cd.n, m))
		sb.WriteString(fmt.Sprintf("%d %d %d\n", cd.t0, cd.t1, cd.t2))
		for _, e := range cd.edges {
			sb.WriteString(fmt.Sprintf("%d %d %d %d\n", e.u, e.v, e.l1, e.l2))
		}
	}
	ti := testInput{name: name, input: sb.String()}
	if err := validateTestInput(ti.input); err != nil {
		panic(fmt.Sprintf("generated invalid test %s: %v", name, err))
	}
	return ti
}

func validateTestInput(input string) error {
	reader := bufio.NewReader(strings.NewReader(input))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return err
	}
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		var n, m int
		if _, err := fmt.Fscan(reader, &n, &m); err != nil {
			return err
		}
		var t0, t1, t2 int64
		if _, err := fmt.Fscan(reader, &t0, &t1, &t2); err != nil {
			return err
		}
		for i := 0; i < m; i++ {
			var u, v int
			var l1, l2 int64
			if _, err := fmt.Fscan(reader, &u, &v, &l1, &l2); err != nil {
				return err
			}
		}
	}
	return nil
}
