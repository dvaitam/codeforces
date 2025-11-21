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
	outPath := filepath.Join(verifierDir, "ref1846F.bin")
	cmd := exec.Command("go", "build", "-o", outPath, "1846F.go")
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

func parseOutputs(out string, expectedCount int) ([]int, error) {
	reader := strings.NewReader(out)
	res := make([]int, 0, expectedCount)
	for {
		var val int
		if _, err := fmt.Fscan(reader, &val); err != nil {
			break
		}
		res = append(res, val)
	}
	if len(res) != expectedCount {
		return nil, fmt.Errorf("expected %d integers, got %d\noutput:\n%s", expectedCount, len(res), out)
	}
	return res, nil
}

func verifyCase(candidate, reference string, tc testCase) error {
	refOut, err := runProgram(reference, tc.input)
	if err != nil {
		return fmt.Errorf("reference error: %v", err)
	}
	var t int
	fmt.Sscan(tc.input, &t)
	refVals, err := parseOutputs(refOut, t)
	if err != nil {
		return fmt.Errorf("invalid reference output: %v", err)
	}

	candOut, err := runProgram(candidate, tc.input)
	if err != nil {
		return fmt.Errorf("candidate error: %v", err)
	}
	candVals, err := parseOutputs(candOut, t)
	if err != nil {
		return fmt.Errorf("invalid candidate output: %v", err)
	}
	for i := 0; i < t; i++ {
		if candVals[i] != refVals[i] {
			return fmt.Errorf("mismatch on test %d: expected %d got %d\ncandidate output:\n%s", i+1, refVals[i], candVals[i], candOut)
		}
	}
	return nil
}

func manualTests() []testCase {
	return []testCase{
		{name: "min_case", input: "1\n2 1\n1 2\n"},
		{name: "simple_case", input: "1\n5 3\n1 2 3 4 5\n"},
	}
}

func randomTest(name string, rng *rand.Rand, maxN int) testCase {
	t := rng.Intn(5) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		n := rng.Intn(maxN-1) + 2
		k := rng.Intn(n) + 1
		fmt.Fprintf(&sb, "%d %d\n", n, k)
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", rng.Intn(9)+1)
		}
		sb.WriteByte('\n')
	}
	return testCase{name: name, input: sb.String()}
}

func generateTests() []testCase {
	tests := manualTests()
	seeds := []int64{1, 2, 3, 4, 5}
	for idx, seed := range seeds {
		rng := rand.New(rand.NewSource(seed))
		tests = append(tests, randomTest(fmt.Sprintf("deterministic_%d", idx+1), rng, 10))
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 60 {
		tests = append(tests, randomTest(fmt.Sprintf("random_%d", len(tests)+1), rng, 50))
	}
	return tests
}

func main() {
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
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
