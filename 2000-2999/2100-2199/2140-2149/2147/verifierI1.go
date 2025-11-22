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

type testCase struct {
	n, m int
}

const limit int64 = 1_000_000_000_000_000_000 // 1e18

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierI1.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[len(os.Args)-1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := generateTests()

	for i, tc := range tests {
		input := fmt.Sprintf("%d %d\n", tc.n, tc.m)

		// Ensure reference produces a valid output on this test.
		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n%s", i+1, err, refOut)
			os.Exit(1)
		}
		if err := validateOutput(refOut, tc); err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d: %v\n", i+1, err)
			os.Exit(1)
		}

		candOut, err := runProgramWithExt(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\n%s", i+1, err, candOut)
			os.Exit(1)
		}
		if err := validateOutput(candOut, tc); err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier location")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "ref-2147I1-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracle2147I1")
	cmd := exec.Command("go", "build", "-o", outPath, "2147I1.go")
	cmd.Dir = dir
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("%v\n%s", err, buf.String())
	}
	cleanup := func() { os.RemoveAll(tmpDir) }
	return outPath, cleanup, nil
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runProgramWithExt(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return errBuf.String(), err
	}
	if errBuf.Len() > 0 {
		return errBuf.String(), fmt.Errorf("unexpected stderr output")
	}
	return out.String(), nil
}

func validateOutput(output string, tc testCase) error {
	tokens := strings.Fields(output)
	if len(tokens) != tc.n {
		return fmt.Errorf("expected %d numbers, got %d", tc.n, len(tokens))
	}
	seq := make([]int64, tc.n)
	for i, tok := range tokens {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return fmt.Errorf("token %q is not an integer", tok)
		}
		if val < -limit || val > limit {
			return fmt.Errorf("value %d out of allowed range [-1e18,1e18]", val)
		}
		seq[i] = val
	}

	// distinct values
	seen := make(map[int64]struct{})
	for _, v := range seq {
		seen[v] = struct{}{}
	}
	if len(seen) > tc.m {
		return fmt.Errorf("too many distinct values: %d > %d", len(seen), tc.m)
	}

	// distance-convex condition
	for i := 1; i+1 < tc.n; i++ {
		prev := absDiff(seq[i], seq[i-1])
		next := absDiff(seq[i+1], seq[i])
		if prev >= next {
			return fmt.Errorf("not distance-convex at position %d: |a[i]-a[i-1]|=%d, |a[i+1]-a[i]|=%d", i+1, prev, next)
		}
	}

	return nil
}

func absDiff(a, b int64) int64 {
	d := a - b
	if d < 0 {
		return -d
	}
	return d
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase

	// Deterministic small cases
	tests = append(tests,
		testCase{n: 1, m: 1},
		testCase{n: 2, m: 2},
		testCase{n: 4, m: 4},
		testCase{n: 6, m: 6},
		testCase{n: 30, m: 40},
		testCase{n: 50, m: 60},
	)

	for len(tests) < 25 {
		n := rng.Intn(40) + 3 // 3..42
		m := n + rng.Intn(10) // ensure m >= n
		if n > 60 {           // keep values from reference within range
			n = 60
			m = 70
		}
		tests = append(tests, testCase{n: n, m: m})
	}
	return tests
}
