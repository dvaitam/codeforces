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
	outPath := filepath.Join(verifierDir, "ref1906G.bin")
	cmd := exec.Command("go", "build", "-o", outPath, "1906G.go")
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

func parseOutput(out string) (string, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return "", fmt.Errorf("empty output")
	}
	res := fields[0]
	if res != "FIRST" && res != "SECOND" {
		return "", fmt.Errorf("invalid output %q", res)
	}
	return res, nil
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
		return fmt.Errorf("expected %s, got %s\ncandidate output:\n%s", expected, got, candOut)
	}
	return nil
}

func manualTests() []testCase {
	return []testCase{
		{name: "single_cell", input: "1\n1 1\n"},
		{name: "two_cells_same_row", input: "2\n1 1\n1 2\n"},
		{name: "two_cells_same_col", input: "2\n1 1\n2 1\n"},
		{name: "far_cells", input: "2\n5 4\n3 7\n"},
	}
}

func randomTest(name string, rng *rand.Rand, maxN int) testCase {
	n := rng.Intn(maxN) + 1
	points := make(map[[2]int]bool)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for len(points) < n {
		r := rng.Intn(1_000_000_000) + 1
		c := rng.Intn(1_000_000_000) + 1
		key := [2]int{r, c}
		if points[key] {
			continue
		}
		points[key] = true
		fmt.Fprintf(&sb, "%d %d\n", r, c)
	}
	return testCase{name: name, input: sb.String()}
}

func deterministicTests() []testCase {
	tests := manualTests()
	seeds := []int64{1, 2, 3, 4, 5}
	for idx, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		tests = append(tests, randomTest(fmt.Sprintf("deterministic_%d", idx+1), rng, 5))
	}
	return tests
}

func generateTests() []testCase {
	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 80 {
		tests = append(tests, randomTest(fmt.Sprintf("random_%d", len(tests)+1), rng, 10))
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
