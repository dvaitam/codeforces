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

const refSource = "2147I2.go"

type testCase struct {
	name string
	n    int
	m    int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierI2.go /path/to/candidate")
		os.Exit(1)
	}
	target := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference solution:", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()

	for idx, tc := range tests {
		input := fmt.Sprintf("%d %d\n", tc.n, tc.m)

		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s\n", idx+1, tc.name, err, input)
			os.Exit(1)
		}
		if err := validateOutput(tc, refOut); err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, input, candOut)
			os.Exit(1)
		}
		if err := validateOutput(tc, candOut); err != nil {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): %v\ninput:\n%sreference output (valid):\n%s\ncandidate output:\n%s\n", idx+1, tc.name, err, input, refOut, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier directory")
	}
	dir := filepath.Dir(file)

	tmp, err := os.CreateTemp("", "oracle-2147I2-")
	if err != nil {
		return "", nil, err
	}
	tmp.Close()
	binPath := tmp.Name()

	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.Remove(binPath)
		return "", nil, fmt.Errorf("reference build failed: %v\n%s", err, stderr.String())
	}

	cleanup := func() {
		_ = os.Remove(binPath)
	}
	return binPath, cleanup, nil
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

func runProgram(bin, input string) (string, error) {
	cmd := commandFor(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func validateOutput(tc testCase, output string) error {
	tokens := strings.Fields(output)
	if len(tokens) != tc.n {
		return fmt.Errorf("expected %d numbers, got %d", tc.n, len(tokens))
	}
	const limit = 1e18
	seq := make([]int64, tc.n)
	seen := make(map[int64]struct{})
	for i, tok := range tokens {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return fmt.Errorf("token %d is not a valid integer: %v", i+1, err)
		}
		if val < -limit || val > limit {
			return fmt.Errorf("value %d out of bounds [%g, %g]", i+1, -limit, limit)
		}
		seq[i] = val
		seen[val] = struct{}{}
	}
	if len(seen) > tc.m {
		return fmt.Errorf("uses %d distinct values which exceeds m=%d", len(seen), tc.m)
	}
	if tc.n >= 2 {
		prev := absDiff(seq[1], seq[0])
		for i := 2; i < tc.n; i++ {
			cur := absDiff(seq[i], seq[i-1])
			if prev >= cur {
				return fmt.Errorf("distance-convex violated at position %d: prev jump %d, current jump %d", i, prev, cur)
			}
			prev = cur
		}
	}
	return nil
}

func absDiff(a, b int64) int64 {
	if a > b {
		return a - b
	}
	return b - a
}

func buildTests() []testCase {
	tests := []testCase{
		{name: "sample_small", n: 8, m: 6},
		{name: "official_large", n: 300000, m: 15000},
		{name: "minimal_n2", n: 2, m: 1},
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 20; i++ {
		tests = append(tests, randomTest(rng, i+1))
	}
	return tests
}

func randomTest(rng *rand.Rand, idx int) testCase {
	n := rng.Intn(20) + 3
	m := rng.Intn(n) + 1
	// occasionally larger n
	if idx%5 == 0 {
		n = rng.Intn(200) + 50
		m = rng.Intn(n) + 1
	}
	return testCase{
		name: fmt.Sprintf("random_%d", idx),
		n:    n,
		m:    m,
	}
}
