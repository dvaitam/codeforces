package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const tol = 1e-6

type testCase struct {
	input string
	n     int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC2.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refSrc, err := locateReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	refBin, err := buildReference(refSrc)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\nInput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		refAmp, refStates, err := parseOutput(refOut, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d: %v\nOutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\nInput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		candAmp, candStates, err := parseOutput(candOut, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d: %v\nOutput:\n%s\n", idx+1, err, candOut)
			os.Exit(1)
		}

		if math.Abs(candAmp-refAmp) > tol && math.Abs((candAmp-refAmp)/refAmp) > tol {
			fmt.Fprintf(os.Stderr, "test %d: amplitude mismatch (expected %.6f, got %.6f)\nInput:\n%s\n", idx+1, refAmp, candAmp, tc.input)
			os.Exit(1)
		}
		if len(candStates) != len(refStates) {
			fmt.Fprintf(os.Stderr, "test %d: expected %d states, got %d\nInput:\n%s\nCandidate output:\n%s\n", idx+1, len(refStates), len(candStates), tc.input, candOut)
			os.Exit(1)
		}
		for state := range candStates {
			if _, ok := refStates[state]; !ok {
				fmt.Fprintf(os.Stderr, "test %d: invalid state %s\nInput:\n%s\nCandidate output:\n%s\n", idx+1, state, tc.input, candOut)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	candidates := []string{
		"1357C2.go",
		filepath.Join("1000-1999", "1300-1399", "1350-1359", "1357", "1357C2.go"),
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("could not locate 1357C2.go")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref1357C2_%d.bin", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", outPath, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return outPath, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr:\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseOutput(out string, n int) (float64, map[string]struct{}, error) {
	tokens := strings.Fields(out)
	if len(tokens) == 0 {
		return 0, nil, fmt.Errorf("empty output")
	}
	var amp float64
	foundAmp := false
	states := make(map[string]struct{})
	for _, tok := range tokens {
		if !foundAmp {
			if val, err := strconv.ParseFloat(tok, 64); err == nil {
				amp = val
				foundAmp = true
				continue
			}
		}
		if len(tok) == n && isBinary(tok) {
			states[tok] = struct{}{}
		}
	}
	if !foundAmp {
		return 0, nil, fmt.Errorf("amplitude not found")
	}
	if len(states) == 0 {
		return 0, nil, fmt.Errorf("no basis states detected")
	}
	return amp, states, nil
}

func isBinary(s string) bool {
	for _, ch := range s {
		if ch != '0' && ch != '1' {
			return false
		}
	}
	return true
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests,
		buildTest(1, 0),
		buildTest(1, 1),
		buildTest(2, 0),
		buildTest(2, 1),
		buildTest(3, 0),
	)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 40 {
		n := rng.Intn(6) + 1 // limit to keep output manageable
		parity := rng.Intn(2)
		tests = append(tests, buildTest(n, parity))
	}
	tests = append(tests, buildTest(8, 0))
	tests = append(tests, buildTest(8, 1))
	return tests
}

func buildTest(n, parity int) testCase {
	if parity != 0 {
		parity = 1
	}
	return testCase{
		input: fmt.Sprintf("%d %d\n", n, parity),
		n:     n,
	}
}
