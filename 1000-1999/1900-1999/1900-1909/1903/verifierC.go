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
	outPath := filepath.Join(verifierDir, "ref1903C.bin")
	cmd := exec.Command("go", "build", "-o", outPath, "1903C.go")
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

func parseOutputs(out string, expected int) ([]int64, error) {
	reader := strings.NewReader(out)
	res := make([]int64, 0, expected)
	for {
		var val int64
		if _, err := fmt.Fscan(reader, &val); err != nil {
			break
		}
		res = append(res, val)
	}
	if len(res) != expected {
		return nil, fmt.Errorf("expected %d outputs, got %d\noutput:\n%s", expected, len(res), out)
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
	expectedVals, err := parseOutputs(refOut, t)
	if err != nil {
		return fmt.Errorf("invalid reference output: %v", err)
	}

	candOut, err := runProgram(candidate, tc.input)
	if err != nil {
		return fmt.Errorf("candidate error: %v", err)
	}
	gotVals, err := parseOutputs(candOut, t)
	if err != nil {
		return fmt.Errorf("invalid candidate output: %v", err)
	}
	for i := 0; i < t; i++ {
		if expectedVals[i] != gotVals[i] {
			return fmt.Errorf("test %d: expected %d got %d\ncandidate output:\n%s", i+1, expectedVals[i], gotVals[i], candOut)
		}
	}
	return nil
}

func formatInput(testCases []testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(testCases))
	for _, tc := range testCases {
		reader := strings.NewReader(tc.input)
		var t int
		fmt.Fscan(reader, &t)
		for caseIdx := 0; caseIdx < t; caseIdx++ {
			var n int
			fmt.Fscan(reader, &n)
			fmt.Fprintf(&sb, "%d\n", n)
			for i := 0; i < n; i++ {
				var val int64
				fmt.Fscan(reader, &val)
				if i > 0 {
					sb.WriteByte(' ')
				}
				fmt.Fprintf(&sb, "%d", val)
			}
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func manualTests() []testCase {
	return []testCase{
		{name: "single_positive", input: "1\n1\n5\n"},
		{name: "single_negative", input: "1\n1\n-3\n"},
		{name: "mixed_small", input: "1\n5\n1 -3 7 -6 2\n"},
		{name: "all_zero", input: "1\n4\n0 0 0 0\n"},
	}
}

func randomTest(name string, rng *rand.Rand, maxN int, maxVal int64) testCase {
	t := rng.Intn(3) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		n := rng.Intn(maxN) + 1
		fmt.Fprintf(&sb, "%d\n", n)
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			val := rng.Int63n(2*maxVal+1) - maxVal
			fmt.Fprintf(&sb, "%d", val)
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
		tests = append(tests, randomTest(fmt.Sprintf("deterministic_%d", idx+1), rng, 8, 20))
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 100 {
		tests = append(tests, randomTest(fmt.Sprintf("random_%d", len(tests)+1), rng, 50, 1_000_000_000))
	}
	return tests
}

func main() {
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
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
