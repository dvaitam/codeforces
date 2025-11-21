package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const referenceSolutionRel = "2000-2999/2000-2099/2030-2039/2036/2036E.go"

var referenceSolutionPath string

func init() {
	referenceSolutionPath = referenceSolutionRel
	if _, file, _, ok := runtime.Caller(0); ok {
		dir := filepath.Dir(file)
		candidate := filepath.Join(dir, "2036E.go")
		if _, err := os.Stat(candidate); err == nil {
			referenceSolutionPath = candidate
			return
		}
	}
	if abs, err := filepath.Abs(referenceSolutionRel); err == nil {
		if _, err := os.Stat(abs); err == nil {
			referenceSolutionPath = abs
		}
	}
}

type requirement struct {
	r  int
	op string
	c  int
}

type query struct {
	m   int
	req []requirement
}

type testCase struct {
	n, k, q int
	a       [][]int
	queries []query
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n: 3, k: 4, q: 4,
			a: [][]int{
				{1, 3, 5, 9},
				{4, 6, 5, 3},
				{2, 1, 2, 7},
			},
			queries: []query{
				{m: 3, req: []requirement{{1, ">", 4}, {1, "<", 6}, {2, "<", 8}}},
				{m: 1, req: []requirement{{1, ">", 6}}},
				{m: 2, req: []requirement{{2, "<", 8}, {3, ">", 8}}},
				{m: 1, req: []requirement{{4, "<", 9}}},
			},
		},
		{
			n: 1, k: 1, q: 2,
			a: [][]int{{5}},
			queries: []query{
				{m: 1, req: []requirement{{1, ">", 4}}},
				{m: 1, req: []requirement{{1, "<", 4}}},
			},
		},
	}
}

func randomMatrix(rng *rand.Rand, n, k int) [][]int {
	a := make([][]int, n)
	for i := 0; i < n; i++ {
		a[i] = make([]int, k)
		for j := 0; j < k; j++ {
			a[i][j] = rng.Intn(50) + 1
		}
	}
	return a
}

func randomQuery(rng *rand.Rand, k int) query {
	m := rng.Intn(5) + 1
	req := make([]requirement, m)
	for i := 0; i < m; i++ {
		op := "<"
		if rng.Intn(2) == 0 {
			op = ">"
		}
		req[i] = requirement{
			r:  rng.Intn(k) + 1,
			op: op,
			c:  rng.Intn(20),
		}
	}
	return query{m: m, req: req}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(20250305))
	var tests []testCase
	totalNK := 0
	totalM := 0
	for len(tests) < 60 && totalNK < 200000 && totalM < 200000 {
		n := rng.Intn(20) + 1
		k := rng.Intn(20) + 1
		if totalNK+n*k > 200000 {
			break
		}
		q := rng.Intn(20) + 1
		queries := make([]query, q)
		for i := 0; i < q; i++ {
			queries[i] = randomQuery(rng, k)
			totalM += queries[i].m
			if totalM > 200000 {
				break
			}
		}
		if totalM > 200000 {
			break
		}
		tests = append(tests, testCase{
			n: n, k: k, q: q,
			a:       randomMatrix(rng, n, k),
			queries: queries,
		})
		totalNK += n * k
	}
	return tests
}

func formatInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", tests[0].n, tests[0].k, tests[0].q))
	for _, row := range tests[0].a {
		for j, val := range row {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", val))
		}
		sb.WriteByte('\n')
	}
	for _, q := range tests[0].queries {
		sb.WriteString(fmt.Sprintf("%d\n", q.m))
		for _, req := range q.req {
			sb.WriteString(fmt.Sprintf("%d %s %d\n", req.r, req.op, req.c))
		}
	}
	return sb.String()
}

func formatMultiTestInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", tests[0].n, tests[0].k, tests[0].q))
	for _, row := range tests[0].a {
		for j, val := range row {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", val))
		}
		sb.WriteByte('\n')
	}
	for _, q := range tests[0].queries {
		sb.WriteString(fmt.Sprintf("%d\n", q.m))
		for _, req := range q.req {
			sb.WriteString(fmt.Sprintf("%d %s %d\n", req.r, req.op, req.c))
		}
	}
	return sb.String()
}

func formatFullInput(tests []testCase) string {
	var sb strings.Builder
	totalTests := len(tests)
	sb.WriteString(fmt.Sprintf("%d %d %d\n", tests[0].n, tests[0].k, tests[0].q))
	for _, row := range tests[0].a {
		for j, val := range row {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", val))
		}
		sb.WriteByte('\n')
	}
	for _, q := range tests[0].queries {
		sb.WriteString(fmt.Sprintf("%d\n", q.m))
		for _, req := range q.req {
			sb.WriteString(fmt.Sprintf("%d %s %d\n", req.r, req.op, req.c))
		}
	}
	return fmt.Sprintf("%d\n%s", totalTests, sb.String())
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildReferenceBinary() (string, func(), error) {
	if referenceSolutionPath == "" {
		return "", nil, fmt.Errorf("reference solution path not set")
	}
	if _, err := os.Stat(referenceSolutionPath); err != nil {
		return "", nil, fmt.Errorf("reference solution not found: %v", err)
	}
	tmpDir, err := os.MkdirTemp("", "2036E-ref")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref_2036E")
	cmd := exec.Command("go", "build", "-o", binPath, referenceSolutionPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return binPath, cleanup, nil
}

func parseOutputs(out string, total int) ([]int, error) {
	scanner := bufio.NewScanner(strings.NewReader(out))
	res := make([]int, 0, total)
	for scanner.Scan() {
		var val int
		if _, err := fmt.Sscan(scanner.Text(), &val); err != nil {
			return nil, fmt.Errorf("failed to parse %q: %v", scanner.Text(), err)
		}
		res = append(res, val)
	}
	if len(res) != total {
		return nil, fmt.Errorf("expected %d outputs, got %d", total, len(res))
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests := append(deterministicTests(), randomTests()...)
	input := formatFullInput(tests)

	refBin, cleanup, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}

	totalAnswers := 0
	for _, tc := range tests {
		totalAnswers += tc.q
	}

	expected, err := parseOutputs(refOut, totalAnswers)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}

	userOut, err := runProgram(bin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\noutput:\n%s\n", err, userOut)
		os.Exit(1)
	}
	got, err := parseOutputs(userOut, totalAnswers)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse participant output: %v\noutput:\n%s\n", err, userOut)
		os.Exit(1)
	}

	for i := range expected {
		if expected[i] != got[i] {
			fmt.Fprintf(os.Stderr, "query %d mismatch: expected %d got %d\n", i+1, expected[i], got[i])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
