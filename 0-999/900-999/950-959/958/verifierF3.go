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

const mod = 1009

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
	outPath := filepath.Join(verifierDir, "ref958F3.bin")
	cmd := exec.Command("go", "build", "-o", outPath, "958F3.go")
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
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(out string) (int, error) {
	reader := strings.NewReader(out)
	var val int
	if _, err := fmt.Fscan(reader, &val); err != nil {
		return 0, fmt.Errorf("failed to read integer: %v\noutput:\n%s", err, out)
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
		return fmt.Errorf("reference runtime error: %v\n%s", err, refOut)
	}
	expected, err := parseOutput(refOut)
	if err != nil {
		return fmt.Errorf("invalid reference output: %v", err)
	}

	candOut, err := runProgram(candidate, tc.input)
	if err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, candOut)
	}
	got, err := parseOutput(candOut)
	if err != nil {
		return err
	}
	if got%mod != expected%mod {
		return fmt.Errorf("expected %d, got %d\ncandidate output:\n%s", expected%mod, got%mod, candOut)
	}
	return nil
}

func formatInput(n, m, k int, colors []int) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)
	for i, c := range colors {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(c))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func manualTests() []testCase {
	return []testCase{
		{name: "single_knight", input: formatInput(1, 1, 1, []int{1})},
		{name: "impossible_k", input: formatInput(5, 3, 5, []int{1, 1, 2, 2, 3})},
		{name: "all_same_color", input: formatInput(6, 1, 3, []int{1, 1, 1, 1, 1, 1})},
		{name: "each_unique", input: formatInput(4, 4, 2, []int{1, 2, 3, 4})},
	}
}

func randomCase(name string, seed int64, maxN int) testCase {
	rng := rand.New(rand.NewSource(seed))
	n := rng.Intn(maxN) + 1
	m := rng.Intn(n) + 1
	k := rng.Intn(n) + 1
	colors := make([]int, n)
	for i := 0; i < n; i++ {
		colors[i] = rng.Intn(m) + 1
	}
	return testCase{
		name:  name,
		input: formatInput(n, m, k, colors),
	}
}

func largeCase(name string, n, m, k int, seed int64) testCase {
	rng := rand.New(rand.NewSource(seed))
	colors := make([]int, n)
	for i := 0; i < n; i++ {
		colors[i] = rng.Intn(m) + 1
	}
	return testCase{
		name:  name,
		input: formatInput(n, m, k, colors),
	}
}

func generateTests() []testCase {
	tests := manualTests()
	deterministicSeeds := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for idx, seed := range deterministicSeeds {
		tests = append(tests, randomCase(fmt.Sprintf("deterministic_%d", idx+1), seed, 50))
	}
	tests = append(tests,
		largeCase("large_equal", 2000, 2000, 1500, 123),
		largeCase("large_sparse", 5000, 700, 2500, 456),
		largeCase("max_like", 200000, 200000, 100000, 789),
	)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 80 {
		maxN := 200
		if len(tests)%15 == 0 {
			maxN = 2000
		}
		tests = append(tests, randomCase(fmt.Sprintf("random_%d", len(tests)+1), rng.Int63(), maxN))
	}
	return tests
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF3.go /path/to/binary")
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
