package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
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
		wantOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\nInput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		want, err := parseOutput(wantOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d: %v\nOutput:\n%s\n", idx+1, err, wantOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\nInput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		got, err := parseOutput(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d: %v\nOutput:\n%s\n", idx+1, err, gotOut)
			os.Exit(1)
		}

		if len(got) != len(want) {
			fmt.Fprintf(os.Stderr, "test %d: expected %d values, got %d\nInput:\n%s\n", idx+1, len(want), len(got), tc.input)
			os.Exit(1)
		}
		for i := range want {
			if want[i] != got[i] {
				fmt.Fprintf(os.Stderr, "test %d mismatch at position %d: expected %d, got %d\nInput:\n%s\nCandidate output:\n%s\n", idx+1, i+1, want[i], got[i], tc.input, gotOut)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	candidates := []string{
		"1774C.go",
		filepath.Join("1000-1999", "1700-1799", "1770-1779", "1774", "1774C.go"),
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("could not find 1774C.go")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref1774C_%d.bin", time.Now().UnixNano()))
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

func parseOutput(out string) ([]int, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) == 0 || len(lines[0]) == 0 {
		return nil, fmt.Errorf("empty output")
	}
	var vals []int
	for _, line := range lines {
		fields := strings.Fields(line)
		for _, tok := range fields {
			v, err := strconv.Atoi(tok)
			if err != nil {
				return nil, fmt.Errorf("invalid integer %q", tok)
			}
			vals = append(vals, v)
		}
	}
	return vals, nil
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests,
		buildTest(2, "0"),
		buildTest(3, "101"),
		buildTest(4, "0011"),
		buildTest(5, "11111"),
	)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 60 {
		n := rng.Intn(50) + 2
		tests = append(tests, randomTest(rng, n))
	}
	tests = append(tests, randomTest(rng, 2000))
	return tests
}

func buildTest(n int, s string) testCase {
	if len(s) != n-1 {
		sb := strings.Builder{}
		for len(sb.String()) < n-1 {
			sb.WriteString("01")
		}
		s = sb.String()[:n-1]
	}
	return testCase{input: fmt.Sprintf("1\n%d\n%s\n", n, s)}
}

func randomTest(rng *rand.Rand, n int) testCase {
	var sb strings.Builder
	for i := 0; i < n-1; i++ {
		if rng.Intn(2) == 0 {
			sb.WriteByte('0')
		} else {
			sb.WriteByte('1')
		}
	}
	return testCase{input: fmt.Sprintf("1\n%d\n%s\n", n, sb.String())}
}
