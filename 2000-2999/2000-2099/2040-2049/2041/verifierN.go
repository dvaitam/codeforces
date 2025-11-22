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
	input    string
	outCount int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierN.go /path/to/binary")
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
		expectOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\nInput:\n%s", idx+1, err, tc.input)
			os.Exit(1)
		}
		expectVals, err := parseOutput(expectOut, tc.outCount)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d: %v\nInput:\n%sOutput:\n%s", idx+1, err, tc.input, expectOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\nInput:\n%s", idx+1, err, tc.input)
			os.Exit(1)
		}
		gotVals, err := parseOutput(gotOut, tc.outCount)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d: %v\nInput:\n%sOutput:\n%s", idx+1, err, tc.input, gotOut)
			os.Exit(1)
		}

		if len(expectVals) != len(gotVals) {
			fmt.Fprintf(os.Stderr, "test %d: output count mismatch, expected %d got %d\nInput:\n%sOutput:\n%s", idx+1, len(expectVals), len(gotVals), tc.input, gotOut)
			os.Exit(1)
		}
		for i := range expectVals {
			if expectVals[i] != gotVals[i] {
				fmt.Fprintf(os.Stderr, "test %d: mismatch at position %d, expected %d got %d\nInput:\n%sOutput:\n%s", idx+1, i+1, expectVals[i], gotVals[i], tc.input, gotOut)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	candidates := []string{
		"2041N.go",
		filepath.Join("2000-2999", "2000-2099", "2040-2049", "2041", "2041N.go"),
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("could not find 2041N.go")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref2041N_%d.bin", time.Now().UnixNano()))
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

func parseOutput(out string, expected int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d integers, got %d", expected, len(fields))
	}
	res := make([]int64, expected)
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", f)
		}
		res[i] = val
	}
	return res, nil
}

func generateTests() []testCase {
	tests := make([]testCase, 0, 50)
	tests = append(tests, sampleTests()...)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 35; i++ {
		tests = append(tests, randomTest(rng))
	}
	return tests
}

func sampleTests() []testCase {
	cases := []string{
		"5 3\n1 2 1 1 1\n1 4\n1 5\n2 5\n",
		"3 2\n1 2 3\n1 2\n2 3\n",
	}
	res := make([]testCase, 0, len(cases))
	for _, in := range cases {
		n := readN(in)
		res = append(res, testCase{input: in, outCount: n})
	}
	return res
}

func readN(input string) int {
	fields := strings.Fields(input)
	if len(fields) == 0 {
		return 0
	}
	n, _ := strconv.Atoi(fields[0])
	return n
}

func randomTest(rng *rand.Rand) testCase {
	n := rng.Intn(40) + 2 // 2..41
	mMax := n * (n - 1) / 2
	m := rng.Intn(minInt(5*n, mMax) + 1) // cap to keep input small-ish

	a := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(1_000_000_000) + 1
	}

	forbids := make(map[int]map[int]struct{})
	addForbidden := func(u, v int) {
		if u > v {
			u, v = v, u
		}
		if forbids[u] == nil {
			forbids[u] = make(map[int]struct{})
		}
		forbids[u][v] = struct{}{}
	}

	for allForbiddenPairs(forbids) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		addForbidden(u, v)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i, val := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(val))
	}
	sb.WriteByte('\n')
	count := 0
	for u, mp := range forbids {
		for v := range mp {
			if count >= m {
				break
			}
			sb.WriteString(fmt.Sprintf("%d %d\n", u, v))
			count++
		}
		if count >= m {
			break
		}
	}

	return testCase{
		input:    sb.String(),
		outCount: n,
	}
}

func allForbiddenPairs(f map[int]map[int]struct{}) int {
	total := 0
	for _, mp := range f {
		total += len(mp)
	}
	return total
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
