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
		fmt.Fprintln(os.Stderr, "usage: go run verifierA3.go /path/to/binary")
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
		wantP, wantG, err := parseOutput(wantOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d: %v\nOutput:\n%s\n", idx+1, err, wantOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\nInput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		gotP, gotG, err := parseOutput(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d: %v\nOutput:\n%s\n", idx+1, err, gotOut)
			os.Exit(1)
		}
		if gotP != wantP || gotG != wantG {
			fmt.Fprintf(os.Stderr, "test %d mismatch.\nExpected: %d %d\nGot: %d %d\nInput:\n%s\n", idx+1, wantP, wantG, gotP, gotG, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	candidates := []string{
		"1184A3.go",
		filepath.Join("1000-1999", "1100-1199", "1180-1189", "1184", "1184A3.go"),
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("could not find 1184A3.go")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref1184A3_%d.bin", time.Now().UnixNano()))
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

func parseOutput(out string) (int, int, error) {
	fields := strings.Fields(out)
	if len(fields) < 2 {
		return 0, 0, fmt.Errorf("expected two integers, got %q", out)
	}
	p, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid prime: %v", err)
	}
	g, err := strconv.Atoi(fields[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid generator: %v", err)
	}
	return p, g, nil
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests,
		buildTest(1, 5, "0", "1"),
		buildTest(2, 6, "01", "10"),
		buildTest(3, 7, "abc", "def"),
	)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 40 {
		n := rng.Intn(6) + 1
		m := rng.Intn(30) + n + 2
		tests = append(tests, randomTest(rng, n, m))
	}
	for len(tests) < 70 {
		n := rng.Intn(15) + 1
		m := rng.Intn(100) + n + 2
		tests = append(tests, randomTest(rng, n, m))
	}
	tests = append(tests, randomTest(rng, 50, 200))
	tests = append(tests, randomTest(rng, 100, 500))
	return tests
}

func buildTest(n, m int, s1, s2 string) testCase {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d %d\n", n, m))
	b.WriteString(fmt.Sprintf("%s\n%s\n", s1, s2))
	return testCase{input: b.String()}
}

func randomTest(rng *rand.Rand, n, m int) testCase {
	if m < n+2 {
		m = n + 2
	}
	if m < 5 {
		m = 5
	}
	s1 := randomString(rng, n)
	s2 := randomString(rng, n)
	return buildTest(n, m, s1, s2)
}

func randomString(rng *rand.Rand, n int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var b strings.Builder
	for i := 0; i < n; i++ {
		b.WriteByte(chars[rng.Intn(len(chars))])
	}
	return b.String()
}
