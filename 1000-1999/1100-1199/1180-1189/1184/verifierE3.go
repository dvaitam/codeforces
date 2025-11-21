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
	m     int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE3.go /path/to/binary")
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
		wantVals, err := parseOutput(wantOut, tc.m)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d: %v\nOutput:\n%s\n", idx+1, err, wantOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\nInput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		gotVals, err := parseOutput(gotOut, tc.m)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d: %v\nOutput:\n%s\n", idx+1, err, gotOut)
			os.Exit(1)
		}
		for i := 0; i < tc.m; i++ {
			if gotVals[i] != wantVals[i] {
				fmt.Fprintf(os.Stderr, "test %d mismatch at edge %d: expected %d, got %d\nInput:\n%s\nCandidate output:\n%s\n", idx+1, i+1, wantVals[i], gotVals[i], tc.input, gotOut)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	candidates := []string{
		"1184E3.go",
		filepath.Join("1000-1999", "1100-1199", "1180-1189", "1184", "1184E3.go"),
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("could not find 1184E3.go")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref1184E3_%d.bin", time.Now().UnixNano()))
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

func parseOutput(out string, m int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != m {
		return nil, fmt.Errorf("expected %d numbers, got %d", m, len(fields))
	}
	vals := make([]int64, m)
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
		newTestCase(2, 1, [][3]int{{1, 2, 0}}),
		newTestCase(3, 3, [][3]int{{1, 2, 5}, {2, 3, 4}, {1, 3, 3}}),
		newSmallRandom(4, 5),
		newSmallRandom(5, 7),
	)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 40 {
		tests = append(tests, randomTest(rng, rng.Intn(20)+2, rng.Intn(80)+rng.Intn(20)+rng.Intn(40)+1))
	}
	tests = append(tests, randomTest(rng, 200, 400))
	tests = append(tests, randomTest(rng, 2000, 4000))
	tests = append(tests, randomTest(rng, 10000, 20000))
	return tests
}

func newTestCase(n, m int, edges [][3]int) testCase {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for _, e := range edges {
		b.WriteString(fmt.Sprintf("%d %d %d\n", e[0], e[1], e[2]))
	}
	return testCase{input: b.String(), m: m}
}

func newSmallRandom(n, m int) testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return randomTest(rng, n, m)
}

func randomTest(rng *rand.Rand, n, m int) testCase {
	if m < n-1 {
		m = n - 1
	}
	if m > n*(n-1)/2 {
		m = n * (n - 1) / 2
	}
	type pair struct{ u, v int }
	used := make(map[pair]struct{})
	var edges [][3]int
	// build tree first
	for i := 2; i <= n; i++ {
		u := i
		v := rng.Intn(i-1) + 1
		w := rng.Intn(1000)
		edges = append(edges, [3]int{u, v, w})
		used[pair{min(u, v), max(u, v)}] = struct{}{}
	}
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n-1) + 1
		if v >= u {
			v++
		}
		a, b := min(u, v), max(u, v)
		if _, ok := used[pair{a, b}]; ok {
			continue
		}
		used[pair{a, b}] = struct{}{}
		w := rng.Intn(1000)
		edges = append(edges, [3]int{u, v, w})
	}
	return newTestCase(n, len(edges), edges)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
