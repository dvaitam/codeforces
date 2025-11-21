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

const referenceSolutionRel = "2000-2999/2000-2099/2020-2029/2025/2025F.go"

var referenceSolutionPath string

func init() {
	referenceSolutionPath = referenceSolutionRel
	if _, file, _, ok := runtime.Caller(0); ok {
		dir := filepath.Dir(file)
		candidate := filepath.Join(dir, "2025F.go")
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

type query struct{ x, y int }
type testCase struct {
	n, q    int
	queries []query
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n: 3,
			q: 4,
			queries: []query{
				{1, 2},
				{3, 2},
				{3, 1},
				{1, 2},
			},
		},
		{
			n: 4,
			q: 4,
			queries: []query{
				{1, 2},
				{2, 3},
				{3, 4},
				{3, 2},
			},
		},
		{
			n: 4,
			q: 2,
			queries: []query{
				{2, 1},
				{4, 3},
			},
		},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(20250305))
	var tests []testCase
	totalQ := 0
	for len(tests) < 80 && totalQ < 5000 {
		n := rng.Intn(20) + 2
		q := rng.Intn(40) + 1
		if totalQ+q > 5000 {
			q = 5000 - totalQ
		}
		if q <= 0 {
			break
		}
		qs := make([]query, q)
		for i := 0; i < q; i++ {
			x := rng.Intn(n) + 1
			y := rng.Intn(n-1) + 1
			if y >= x {
				y++
			}
			qs[i] = query{x, y}
		}
		tests = append(tests, testCase{n: n, q: q, queries: qs})
		totalQ += q
	}
	return tests
}

func formatInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.q))
		for _, qu := range tc.queries {
			sb.WriteString(fmt.Sprintf("%d %d\n", qu.x, qu.y))
		}
	}
	return sb.String()
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
	tmpDir, err := os.MkdirTemp("", "2025F-ref")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref_2025F")
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

func parseOperations(output string, total int) ([]string, error) {
	scanner := bufio.NewScanner(strings.NewReader(output))
	scanner.Split(bufio.ScanWords)
	ops := make([]string, 0, total)
	for scanner.Scan() {
		ops = append(ops, scanner.Text())
	}
	if len(ops) != total {
		return nil, fmt.Errorf("expected %d operations, got %d", total, len(ops))
	}
	return ops, nil
}

func simulate(n int, queries []query, ops []string) (int64, error) {
	val := make([]int64, n+1)
	if len(ops) != len(queries) {
		return 0, fmt.Errorf("operation count mismatch")
	}
	for i, op := range ops {
		if len(op) == 0 {
			return 0, fmt.Errorf("empty operation at %d", i+1)
		}
		act := op
		if len(op) > 2 {
			act = string([]byte{op[0], op[len(op)-1]})
		}
		var p int
		switch act[0] {
		case 'x', 'X':
			p = queries[i].x
		case 'y', 'Y':
			p = queries[i].y
		default:
			return 0, fmt.Errorf("invalid choice %q at step %d", act[0], i+1)
		}
		var delta int64
		switch act[1] {
		case '+':
			delta = 1
		case '-':
			delta = -1
		default:
			return 0, fmt.Errorf("invalid sign %q at step %d", act[1], i+1)
		}
		if val[p]+delta < 0 {
			return 0, fmt.Errorf("negative value at position %d after step %d", p, i+1)
		}
		val[p] += delta
	}
	var sum int64
	for i := 1; i <= n; i++ {
		sum += val[i]
	}
	return sum, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests := append(deterministicTests(), randomTests()...)
	input := formatInput(tests)

	totalOps := 0
	for _, tc := range tests {
		totalOps += tc.q
	}

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
	refOps, err := parseOperations(refOut, totalOps)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}

	targetSums := make([]int64, len(tests))
	idx := 0
	for ti, tc := range tests {
		sum, err := simulate(tc.n, tc.queries, refOps[idx:idx+tc.q])
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference simulation error on test %d: %v\n", ti+1, err)
			os.Exit(1)
		}
		targetSums[ti] = sum
		idx += tc.q
	}

	userOut, err := runProgram(bin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\noutput:\n%s\n", err, userOut)
		os.Exit(1)
	}
	userOps, err := parseOperations(userOut, totalOps)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse participant output: %v\noutput:\n%s\n", err, userOut)
		os.Exit(1)
	}

	idx = 0
	for ti, tc := range tests {
		sum, err := simulate(tc.n, tc.queries, userOps[idx:idx+tc.q])
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d invalid: %v\n", ti+1, err)
			os.Exit(1)
		}
		if sum != targetSums[ti] {
			fmt.Fprintf(os.Stderr, "test %d failed: expected final sum %d, got %d\n", ti+1, targetSums[ti], sum)
			os.Exit(1)
		}
		idx += tc.q
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
