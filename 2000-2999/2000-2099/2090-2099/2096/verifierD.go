package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	refSource2096D = "2000-2999/2000-2099/2090-2099/2096/2096D.go"
	maxTotalN      = 200000
)

type testCase struct {
	points [][2]int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSource2096D)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := serializeTests(tests)

	// Sanity-check with reference solution.
	if err := runAndValidate(refBin, tests, input); err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}

	if err := runAndValidate(candidate, tests, input); err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference(source string) (string, error) {
	tmp, err := os.CreateTemp("", "2096D-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	srcPath, err := resolveSourcePath(source)
	if err != nil {
		os.Remove(tmp.Name())
		return "", err
	}

	cmd := exec.Command("go", "build", "-o", tmp.Name(), srcPath)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, buf.String())
	}
	return tmp.Name(), nil
}

func resolveSourcePath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, path), nil
}

func runAndValidate(target string, tests []testCase, input string) error {
	out, err := runProgram(target, input)
	if err != nil {
		return err
	}
	results, err := parseOutputs(out, len(tests))
	if err != nil {
		return err
	}
	for i, tc := range tests {
		if err := validateCase(tc, results[i][0], results[i][1]); err != nil {
			return fmt.Errorf("test %d: %v", i+1, err)
		}
	}
	return nil
}

func runProgram(target, input string) (string, error) {
	cmd := commandFor(target)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return stdout.String(), nil
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func parseOutputs(out string, t int) ([][2]int64, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	ans := make([][2]int64, t)
	for i := 0; i < t; i++ {
		if _, err := fmt.Fscan(reader, &ans[i][0], &ans[i][1]); err != nil {
			return nil, fmt.Errorf("expected %d answer pairs, got %d (%v)", t, i, err)
		}
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("extra output detected starting with %q", extra)
	}
	return ans, nil
}

func validateCase(tc testCase, s, t int64) error {
	if s < -1_000_000_000 || s > 1_000_000_000 || t < -1_000_000_000 || t > 1_000_000_000 {
		return fmt.Errorf("output coordinates out of range: (%d, %d)", s, t)
	}

	colParity := make(map[int64]byte)
	diagParity := make(map[int64]byte)
	for _, p := range tc.points {
		colParity[p[0]] ^= 1
		diagParity[p[0]+p[1]] ^= 1
	}

	oddCol, oddDiag := int64(0), int64(0)
	colFound, diagFound := false, false
	for k, v := range colParity {
		if v&1 == 1 {
			if colFound {
				return fmt.Errorf("invalid configuration: multiple columns with odd parity")
			}
			colFound = true
			oddCol = k
		}
	}
	for k, v := range diagParity {
		if v&1 == 1 {
			if diagFound {
				return fmt.Errorf("invalid configuration: multiple diagonals with odd parity")
			}
			diagFound = true
			oddDiag = k
		}
	}

	if !colFound || !diagFound {
		return fmt.Errorf("invalid configuration: missing odd parity in column or diagonal")
	}
	if s != oddCol {
		return fmt.Errorf("column parity mismatch: expected odd column %d, got %d", oddCol, s)
	}
	if s+t != oddDiag {
		return fmt.Errorf("diagonal parity mismatch: expected odd diag sum %d, got %d", oddDiag, s+t)
	}
	return nil
}

func serializeTests(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(len(tc.points)))
		sb.WriteByte('\n')
		for _, p := range tc.points {
			sb.WriteString(strconv.FormatInt(p[0], 10))
			sb.WriteByte(' ')
			sb.WriteString(strconv.FormatInt(p[1], 10))
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func buildTests() []testCase {
	tests := make([]testCase, 0)
	total := 0

	add := func(tc testCase) {
		tests = append(tests, tc)
		total += len(tc.points)
	}

	// Simple deterministic cases.
	add(fromOperations(0, 0, nil))                                   // single bulb
	add(fromOperations(0, 0, [][2]int64{{0, 0}}))                    // single toggle
	add(fromOperations(2, -3, [][2]int64{{1, 1}, {-2, 5}, {7, -4}})) // few toggles
	add(disjointCase(5000, 5, 0))                                    // small disjoint
	add(disjointCase(20000, 10, 3))                                  // mid-size

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Randomized cases.
	for total < maxTotalN && len(tests) < 30 {
		opCnt := rng.Intn(120) + 30
		baseX := int64(rng.Intn(2000) - 1000)
		baseY := int64(rng.Intn(2000) - 1000)
		tc := randomOperationsCase(rng, baseX, baseY, opCnt)
		if total+len(tc.points) > maxTotalN {
			break
		}
		add(tc)
	}

	// Large stress case.
	if total < maxTotalN {
		remaining := maxTotalN - total
		if remaining == 0 {
			return tests
		}
		k := (remaining - 1) / 4
		if k > 45000 {
			k = 45000
		}
		need := 1 + 4*k
		if need > remaining {
			need = remaining
		}
		add(disjointCase(need, 1000, 1000)) // ensures n = 1 + 4*k or trimmed by need
	}

	return tests
}

// fromOperations builds a test by starting with (s,t) lit and applying toggles at positions ops.
func fromOperations(s, t int64, ops [][2]int64) testCase {
	bulbs := make(map[[2]int64]bool)
	bulbs[[2]int64{s, t}] = true
	for _, op := range ops {
		toggle(bulbs, op[0], op[1])
	}
	return testCase{points: mapToSlice(bulbs)}
}

// disjointCase generates n bulbs using non-overlapping operations to reach a large size.
// targetN should be of form 1 + 4*k; if not, it will be rounded down appropriately.
func disjointCase(targetN int, startX, startY int64) testCase {
	if targetN < 1 {
		targetN = 1
	}
	k := (targetN - 1) / 4
	bulbs := make(map[[2]int64]bool)
	bulbs[[2]int64{startX, startY}] = true
	for i := 0; i < k; i++ {
		x := startX + int64(3*i+10)
		y := startY + int64(i%5)
		toggle(bulbs, x, y)
	}
	return testCase{points: mapToSlice(bulbs)}
}

func randomOperationsCase(rng *rand.Rand, s, t int64, ops int) testCase {
	bulbs := make(map[[2]int64]bool)
	bulbs[[2]int64{s, t}] = true
	for i := 0; i < ops; i++ {
		x := s + int64(rng.Intn(200)-100)
		y := t + int64(rng.Intn(200)-100)
		toggle(bulbs, x, y)
	}
	return testCase{points: mapToSlice(bulbs)}
}

func toggle(bulbs map[[2]int64]bool, x, y int64) {
	points := [][2]int64{
		{x, y},
		{x, y + 1},
		{x + 1, y - 1},
		{x + 1, y},
	}
	for _, p := range points {
		if bulbs[[2]int64{p[0], p[1]}] {
			delete(bulbs, [2]int64{p[0], p[1]})
		} else {
			bulbs[[2]int64{p[0], p[1]}] = true
		}
	}
}

func mapToSlice(bulbs map[[2]int64]bool) [][2]int64 {
	res := make([][2]int64, 0, len(bulbs))
	for p := range bulbs {
		res = append(res, p)
	}
	return res
}
