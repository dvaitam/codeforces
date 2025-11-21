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
		fmt.Fprintln(os.Stderr, "usage: go run verifierA2.go /path/to/binary")
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
		wantVals, err := parseOutput(wantOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d: %v\nOutput:\n%s\n", idx+1, err, wantOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\nInput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		gotVals, err := parseOutput(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d: %v\nOutput:\n%s\n", idx+1, err, gotOut)
			os.Exit(1)
		}

		if len(gotVals) != len(wantVals) {
			fmt.Fprintf(os.Stderr, "test %d: expected %d values, got %d\n", idx+1, len(wantVals), len(gotVals))
			os.Exit(1)
		}
		for i := range wantVals {
			if gotVals[i] != wantVals[i] {
				fmt.Fprintf(os.Stderr, "test %d: mismatch at position %d. expected %d, got %d\nInput:\n%s\n", idx+1, i+1, wantVals[i], gotVals[i], tc.input)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	candidates := []string{
		"1129A2.go",
		filepath.Join("1000-1999", "1100-1199", "1120-1129", "1129", "1129A2.go"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("could not find 1129A2.go")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref1129A2_%d.bin", time.Now().UnixNano()))
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

func parseOutput(out string) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return nil, fmt.Errorf("empty output")
	}
	vals := make([]int64, len(fields))
	for i, tok := range fields {
		v, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", tok)
		}
		vals[i] = v
	}
	return vals, nil
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests,
		buildTest(2, [][2]int{{1, 2}}),
		buildTest(3, [][2]int{{1, 2}, {2, 3}, {3, 1}}),
		buildTest(4, [][2]int{{1, 3}, {2, 4}, {3, 1}, {4, 2}}),
	)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 60 {
		n := rng.Intn(40) + 2
		m := rng.Intn(200) + 1
		tests = append(tests, randomTest(rng, n, m))
	}
	tests = append(tests, randomTest(rng, 5000, 20000))
	return tests
}

func buildTest(n int, edges [][2]int) testCase {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d %d\n", n, len(edges)))
	for _, e := range edges {
		b.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return testCase{input: b.String()}
}

func randomTest(rng *rand.Rand, n, m int) testCase {
	var edges [][2]int
	for i := 0; i < m; i++ {
		a := rng.Intn(n) + 1
		b := rng.Intn(n-1) + 1
		if b >= a {
			b++
		}
		edges = append(edges, [2]int{a, b})
	}
	return buildTest(n, edges)
}
