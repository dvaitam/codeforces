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
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
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
		want, err := parseOutput(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d: %v\nOutput:\n%s\n", idx+1, err, refOut)
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
		if want != got {
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %d, got %d\nInput:\n%s\nCandidate output:\n%s\n", idx+1, want, got, tc.input, gotOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	src := os.Getenv("REFERENCE_SOURCE_PATH")
	if src == "" {
		return "", fmt.Errorf("REFERENCE_SOURCE_PATH not set")
	}
	return src, nil
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref1487G_%d.bin", time.Now().UnixNano()))
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

func parseOutput(out string) (int64, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q", fields[0])
	}
	return val, nil
}

func generateTests() []testCase {
	var tests []testCase
	// Fixed small test cases respecting constraints: n >= 3, c_i > n/3 for all i
	tests = append(tests,
		buildTest(3, validCounts(3, 2)),  // c_i > 1 for all, so c_i >= 2
		buildTest(3, validCounts(3, 3)),  // c_i = 3 for all
		buildTest(6, validCounts(6, 3)),  // c_i > 2, so c_i >= 3
	)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 50 {
		n := rng.Intn(48) + 3 // n in [3, 50]
		minC := n/3 + 1       // c_i > n/3, so c_i >= n/3 + 1
		c := make([]int, 26)
		for i := 0; i < 26; i++ {
			c[i] = minC + rng.Intn(n-minC+1) // c_i in [minC, n]
		}
		tests = append(tests, buildTest(n, c))
	}
	tests = append(tests, buildTest(200, randomValidCounts(200)))
	return tests
}

func validCounts(n, val int) []int {
	c := make([]int, 26)
	for i := range c {
		c[i] = val
	}
	return c
}

func buildTest(n int, counts []int) testCase {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < 26; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		val := 0
		if i < len(counts) {
			val = counts[i]
		}
		b.WriteString(strconv.Itoa(val))
	}
	b.WriteByte('\n')
	return testCase{input: b.String()}
}

func randomValidCounts(n int) []int {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	minC := n/3 + 1
	c := make([]int, 26)
	for i := range c {
		c[i] = minC + rng.Intn(n-minC+1)
	}
	return c
}
