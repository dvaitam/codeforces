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
	name  string
	input string
}

var verifierDir string

func init() {
	if _, file, _, ok := runtime.Caller(0); ok {
		verifierDir = filepath.Dir(file)
	} else {
		verifierDir = "."
	}
}

func buildReference() (string, error) {
	outPath := filepath.Join(verifierDir, "ref1906M.bin")
	cmd := exec.Command("go", "build", "-o", outPath, "1906M.go")
	cmd.Dir = verifierDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return outPath, nil
}

func runProgram(target, input string) (string, error) {
	if !filepath.IsAbs(target) {
		if abs, err := filepath.Abs(target); err == nil {
			target = abs
		}
	}
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
	return out.String(), nil
}

func parseOutput(out string) (int64, error) {
	var val int64
	if _, err := fmt.Fscan(strings.NewReader(out), &val); err != nil {
		return 0, fmt.Errorf("failed to parse output: %v\n%s", err, out)
	}
	return val, nil
}

func verifyCase(candidate, reference string, tc testCase) error {
	refOut, err := runProgram(reference, tc.input)
	if err != nil {
		return fmt.Errorf("reference error: %v", err)
	}
	expected, err := parseOutput(refOut)
	if err != nil {
		return fmt.Errorf("invalid reference output: %v", err)
	}

	candOut, err := runProgram(candidate, tc.input)
	if err != nil {
		return fmt.Errorf("candidate error: %v", err)
	}
	got, err := parseOutput(candOut)
	if err != nil {
		return fmt.Errorf("invalid candidate output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d, got %d\ncandidate output:\n%s", expected, got, candOut)
	}
	return nil
}

func manualTests() []testCase {
	return []testCase{
		{name: "triangle_from_three_sides", input: "3\n1 1 1\n"},
		{name: "large_equal", input: "4\n1000000000 1000000000 1000000000 1000000000\n"},
		{name: "mixed_small", input: "5\n1 2 3 4 5\n"},
	}
}

func randomTest(name string, rng *rand.Rand, maxN int, maxVal int64) testCase {
	n := rng.Intn(maxN-2) + 3
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		val := rng.Int63n(maxVal) + 1
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", val)
	}
	sb.WriteByte('\n')
	return testCase{name: name, input: sb.String()}
}

func generateTests() []testCase {
	tests := manualTests()
	seeds := []int64{1, 2, 3, 4, 5}
	for idx, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		tests = append(tests, randomTest(fmt.Sprintf("deterministic_%d", idx+1), rng, 8, 20))
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 100 {
		tests = append(tests, randomTest(fmt.Sprintf("random_%d", len(tests)+1), rng, 200000, 2_000_000_000))
	}
	return tests
}

func main() {
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierM.go /path/to/binary")
		os.Exit(1)
	}
	candidate := args[0]

	ref, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := generateTests()
	for i, tc := range tests {
		if err := verifyCase(candidate, ref, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d (%s) failed: %v\ninput:\n%s", i+1, tc.name, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
