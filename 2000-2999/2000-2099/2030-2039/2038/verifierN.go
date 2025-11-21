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
)

const referenceSolutionRel = "2000-2999/2000-2099/2030-2039/2038/2038N.go"

var referenceSolutionPath string

func init() {
	referenceSolutionPath = referenceSolutionRel
	if _, file, _, ok := runtime.Caller(0); ok {
		dir := filepath.Dir(file)
		candidate := filepath.Join(dir, "2038N.go")
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

type testCase struct {
	s string
}

func deterministicTests() []testCase {
	return []testCase{
		{s: "3<7"},
		{s: "3>7"},
		{s: "8=9"},
		{s: "0=0"},
		{s: "5<3"},
		{s: "0<9"},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(20250305))
	var tests []testCase
	for len(tests) < 200 {
		a := byte('0' + rng.Intn(10))
		b := byte('0' + rng.Intn(10))
		op := []byte("<=>")[rng.Intn(3)]
		tests = append(tests, testCase{s: string([]byte{a, op, b})})
	}
	return tests
}

func formatInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(tc.s)
		sb.WriteByte('\n')
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
	tmpDir, err := os.MkdirTemp("", "2038N-ref")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref_2038N")
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

func parseOutputs(out string, count int) ([]string, error) {
	lines := strings.Fields(out)
	if len(lines) != count {
		return nil, fmt.Errorf("expected %d expressions, got %d", count, len(lines))
	}
	return lines, nil
}

func cost(orig, cand string) int {
	c := 0
	for i := 0; i < 3; i++ {
		if orig[i] != cand[i] {
			c++
		}
	}
	return c
}

func isValid(expr string) bool {
	if len(expr) != 3 {
		return false
	}
	a, op, b := expr[0], expr[1], expr[2]
	if a < '0' || a > '9' || b < '0' || b > '9' {
		return false
	}
	var truth bool
	switch op {
	case '<':
		truth = a < b
	case '>':
		truth = a > b
	case '=':
		truth = a == b
	default:
		return false
	}
	return truth
}

func bestCost(s string) int {
	best := 4
	for a := byte('0'); a <= '9'; a++ {
		for _, op := range []byte("<=>") {
			for b := byte('0'); b <= '9'; b++ {
				expr := []byte{a, op, b}
				if !isValid(string(expr)) {
					continue
				}
				c := cost(s, string(expr))
				if c < best {
					best = c
				}
			}
		}
	}
	return best
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierN.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests := append(deterministicTests(), randomTests()...)
	input := formatInput(tests)

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
	_, err = parseOutputs(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference output parse error: %v\n", err)
		os.Exit(1)
	}

	userOut, err := runProgram(bin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\noutput:\n%s\n", err, userOut)
		os.Exit(1)
	}
	userExprs, err := parseOutputs(userOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "participant output parse error: %v\n", err)
		os.Exit(1)
	}

	for i, expr := range userExprs {
		orig := tests[i].s
		if !isValid(expr) {
			fmt.Fprintf(os.Stderr, "test %d: expression %q is invalid\n", i+1, expr)
			os.Exit(1)
		}
		if cost(orig, expr) != bestCost(orig) {
			fmt.Fprintf(os.Stderr, "test %d: expression %q is not minimal for %q\n", i+1, expr, orig)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
