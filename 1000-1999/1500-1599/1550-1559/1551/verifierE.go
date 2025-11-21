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
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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
		"1551E.go",
		filepath.Join("1000-1999", "1500-1599", "1550-1559", "1551", "1551E.go"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("could not find 1551E.go")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref1551E_%d.bin", time.Now().UnixNano()))
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
	fields := strings.Fields(out)
	res := make([]int, len(fields))
	for i, tok := range fields {
		val, err := strconv.Atoi(tok)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", tok)
		}
		res[i] = val
	}
	return res, nil
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests,
		buildTest([][]int{{1, 1}}),
		buildTest([][]int{{1, 1}, {1, 1}}),
		buildTest([][]int{{2, 1, 1, 2}}),
	)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 60 {
		tests = append(tests, randomTest(rng, rng.Intn(5)+1, rng.Intn(5)+1))
	}
	tests = append(tests, randomTest(rng, 50, 30))
	return tests
}

func buildTest(cases [][]int) testCase {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, cs := range cases {
		n := cs[0]
		k := 1
		if len(cs) > 1 {
			k = cs[1]
		}
		if k > n {
			k = n
		}
		b.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for i := 0; i < n; i++ {
			if i > 0 {
				b.WriteByte(' ')
			}
			if 2+i < len(cs) {
				b.WriteString(strconv.Itoa(cs[2+i]))
			} else {
				b.WriteString(strconv.Itoa((i % n) + 1))
			}
		}
		b.WriteByte('\n')
	}
	return testCase{input: b.String()}
}

func randomTest(rng *rand.Rand, t, maxN int) testCase {
	if t < 1 {
		t = 1
	}
	if maxN < 1 {
		maxN = 1
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := rng.Intn(maxN) + 1
		k := rng.Intn(n) + 1
		b.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for j := 0; j < n; j++ {
			if j > 0 {
				b.WriteByte(' ')
			}
			val := rng.Intn(n) + 1
			b.WriteString(strconv.Itoa(val))
		}
		b.WriteByte('\n')
	}
	return testCase{input: b.String()}
}
