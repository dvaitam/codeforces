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

const referenceSolutionRel = "1000-1999/1800-1899/1810-1819/1812/1812B.go"

var referenceSolutionPath string

func init() {
	referenceSolutionPath = referenceSolutionRel
	if _, file, _, ok := runtime.Caller(0); ok {
		dir := filepath.Dir(file)
		candidate := filepath.Join(dir, "1812B.go")
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
	n int
}

func generateTests() []testCase {
	tests := []testCase{
		{n: 1},
		{n: 2},
		{n: 25},
		{n: 13},
	}
	rng := rand.New(rand.NewSource(20250305))
	for len(tests) < 50 {
		tests = append(tests, testCase{n: rng.Intn(25) + 1})
	}
	return tests
}

func formatInput(tc testCase) string {
	return fmt.Sprintf("%d\n", tc.n)
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
	tmpDir, err := os.MkdirTemp("", "1812B-ref")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref_1812B")
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

func parseOutput(out string) (int, error) {
	out = strings.TrimSpace(out)
	if out == "" {
		return 0, fmt.Errorf("empty output")
	}
	var val int
	if _, err := fmt.Sscanf(out, "%d", &val); err != nil {
		return 0, fmt.Errorf("failed to parse integer from %q", out)
	}
	return val, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	refBin, cleanup, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := generateTests()
	for idx, tc := range tests {
		input := formatInput(tc)
		expectedStr, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", idx+1, err, input, expectedStr)
			os.Exit(1)
		}
		expected, err := parseOutput(expectedStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference parse error on test %d: %v\ninput:\n%soutput:\n%s\n", idx+1, err, input, expectedStr)
			os.Exit(1)
		}

		outStr, runErr := runProgram(bin, input)
		if runErr != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", idx+1, runErr, input, outStr)
			os.Exit(1)
		}
		got, err := parseOutput(outStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse error on test %d: %v\ninput:\n%soutput:\n%s\n", idx+1, err, input, outStr)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%sexpected: %d\ngot: %d\nreference output:\n%s\nparticipant output:\n%s\n", idx+1, input, expected, got, expectedStr, outStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
