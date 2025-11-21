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
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
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
			fmt.Fprintf(os.Stderr, "test %d: expected %d answers, got %d\n", idx+1, len(wantVals), len(gotVals))
			os.Exit(1)
		}
		for i := range wantVals {
			if wantVals[i] != gotVals[i] {
				fmt.Fprintf(os.Stderr, "test %d: mismatch on query %d. expected %d, got %d\nInput:\n%s\n", idx+1, i+1, wantVals[i], gotVals[i], tc.input)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	candidates := []string{
		"1221D.go",
		filepath.Join("1000-1999", "1200-1299", "1220-1229", "1221", "1221D.go"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("could not find 1221D.go")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref1221D_%d.bin", time.Now().UnixNano()))
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
		newTestCase([]board{{1, 1}}),
		newTestCase([]board{{1, 1}, {1, 1}}),
		newTestCase([]board{{3, 2}, {3, 5}, {3, 4}}),
	)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 60 {
		tests = append(tests, randomTest(rng, rng.Intn(20)+1, rng.Intn(4)+1))
	}
	tests = append(tests, randomTest(rng, 1000, 3))
	tests = append(tests, randomTest(rng, 2000, 2))
	return tests
}

type board struct {
	a, b int64
}

func newTestCase(boards []board) testCase {
	var b strings.Builder
	b.WriteString("1\n")
	b.WriteString(fmt.Sprintf("%d\n", len(boards)))
	for _, bo := range boards {
		b.WriteString(fmt.Sprintf("%d %d\n", bo.a, bo.b))
	}
	return testCase{input: b.String()}
}

func randomTest(rng *rand.Rand, n, q int) testCase {
	if q < 1 {
		q = 1
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d\n", q))
	for i := 0; i < q; i++ {
		if n < 1 {
			n = 1
		}
		b.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			a := rng.Int63n(5_000_000) + 1
			c := rng.Int63n(5_000_000) + 1
			b.WriteString(fmt.Sprintf("%d %d\n", a, c))
		}
	}
	return testCase{input: b.String()}
}
