package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const referenceSolutionRel = "2000-2999/2100-2199/2110-2119/2112/2112E.go"

var referenceSolutionPath string

func init() {
	referenceSolutionPath = referenceSolutionRel
	if _, file, _, ok := runtime.Caller(0); ok {
		dir := filepath.Dir(file)
		candidate := filepath.Join(dir, "2112E.go")
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
	name string
	m    int
}

func buildInput(cases []testCase) string {
	var sb strings.Builder
	sb.Grow(len(cases) * 16)
	sb.WriteString(strconv.Itoa(len(cases)))
	sb.WriteByte('\n')
	for _, tc := range cases {
		sb.WriteString(strconv.Itoa(tc.m))
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
	tmpDir, err := os.MkdirTemp("", "2112E-ref-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref_2112E")
	cmd := exec.Command("go", "build", "-o", binPath, referenceSolutionPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	cleanup := func() { _ = os.RemoveAll(tmpDir) }
	return binPath, cleanup, nil
}

func parseOutputs(out string, t int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d integers, got %d", t, len(fields))
	}
	res := make([]int64, t)
	for i, f := range fields {
		v, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q at position %d", f, i+1)
		}
		res[i] = v
	}
	return res, nil
}

func deterministicTests() []testCase {
	values := []int{1, 3, 5, 7, 9, 2, 4, 6, 8, 10, 12, 18, 27}
	cases := make([]testCase, len(values))
	for i, v := range values {
		cases[i] = testCase{name: fmt.Sprintf("det_%d", v), m: v}
	}
	return cases
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	for i := 0; i < 200; i++ {
		val := rng.Intn(500_000) + 1
		tests = append(tests, testCase{
			name: fmt.Sprintf("random_%d", i+1),
			m:    val,
		})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	refBin, cleanup, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := append(deterministicTests(), randomTests()...)

	// We batch all cases in one run to reuse preprocessing of solutions.
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\ninput:\n%soutput:\n%s", err, input, refOut)
		os.Exit(1)
	}
	refVals, err := parseOutputs(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid reference output: %v\noutput:\n%s", err, refOut)
		os.Exit(1)
	}

	out, err := runProgram(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\ninput:\n%soutput:\n%s", err, input, out)
		os.Exit(1)
	}
	vals, err := parseOutputs(out, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid candidate output: %v\ninput:\n%soutput:\n%s", err, input, out)
		os.Exit(1)
	}

	for i := range vals {
		if vals[i] != refVals[i] {
			fmt.Fprintf(os.Stderr, "mismatch on test %d (%s): expected %d got %d\nm=%d\nreference output:\n%s\ncandidate output:\n%s",
				i+1, tests[i].name, refVals[i], vals[i], tests[i].m, refOut, out)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed.")
}
