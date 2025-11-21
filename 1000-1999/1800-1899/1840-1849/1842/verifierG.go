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

const mod = 1000000007

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
	outPath := filepath.Join(verifierDir, "ref1842G.bin")
	cmd := exec.Command("go", "build", "-o", outPath, "1842G.go")
	cmd.Dir = verifierDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
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
		return 0, fmt.Errorf("failed to parse integer: %v\noutput:\n%s", err, out)
	}
	val %= mod
	if val < 0 {
		val += mod
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
		return fmt.Errorf("expected %d, got %d", expected, got)
	}
	return nil
}

func formatInput(n int, m, v int64, a []int64) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, v)
	for i, val := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", val)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func manualTests() []testCase {
	return []testCase{
		{name: "single_small", input: formatInput(1, 1, 1, []int64{1})},
		{name: "sample_like", input: formatInput(2, 2, 5, []int64{2, 12})},
		{name: "no_update", input: formatInput(3, 1, 0, []int64{3, 4, 5})},
	}
}

func randomCase(name string, rng *rand.Rand, maxN int) testCase {
	n := rng.Intn(maxN) + 1
	m := rng.Int63n(1_000_000_000) + 1
	v := rng.Int63n(1_000_000_000) + 1
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Int63n(1_000_000_000) + 1
	}
	return testCase{
		name:  name,
		input: formatInput(n, m, v, a),
	}
}

func generateTests() []testCase {
	tests := manualTests()
	seeds := []int64{1, 2, 3, 4, 5}
	for idx, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		tests = append(tests, randomCase(fmt.Sprintf("deterministic_%d", idx+1), rng, 8))
	}
	tests = append(tests,
		testCase{name: "large_v", input: formatInput(5, 1000000000, 1000000000, []int64{1, 2, 3, 4, 5})},
		testCase{name: "max_n_small", input: formatInput(10, 1000000, 999999937, []int64{7, 3, 5, 9, 11, 2, 4, 8, 6, 1})},
	)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 80 {
		tests = append(tests, randomCase(fmt.Sprintf("random_%d", len(tests)+1), rng, 12))
	}
	return tests
}

func main() {
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
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
