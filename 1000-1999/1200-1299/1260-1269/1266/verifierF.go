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
	n     int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
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
		want, err := parseOutput(wantOut, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d: %v\nOutput:\n%s\n", idx+1, err, wantOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\nInput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		got, err := parseOutput(gotOut, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d: %v\nOutput:\n%s\n", idx+1, err, gotOut)
			os.Exit(1)
		}

		for i := 0; i < tc.n; i++ {
			if got[i] != want[i] {
				fmt.Fprintf(os.Stderr, "test %d mismatch at position %d: expected %d, got %d\nInput:\n%s\n", idx+1, i+1, want[i], got[i], tc.input)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	candidates := []string{
		"1266F.go",
		filepath.Join("1000-1999", "1200-1299", "1260-1269", "1266", "1266F.go"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("could not find 1266F.go")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref1266F_%d.bin", time.Now().UnixNano()))
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

func parseOutput(out string, n int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != n {
		return nil, fmt.Errorf("expected %d integers, got %d", n, len(fields))
	}
	res := make([]int, n)
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
		buildTest(2, [][2]int{{1, 2}}),
		buildTest(3, [][2]int{{1, 2}, {2, 3}}),
		buildTest(4, [][2]int{{1, 2}, {2, 3}, {3, 4}}),
		buildTest(5, [][2]int{{1, 2}, {1, 3}, {1, 4}, {1, 5}}),
	)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 60 {
		n := rng.Intn(10) + 2
		tests = append(tests, randomTree(rng, n))
	}
	for len(tests) < 80 {
		n := rng.Intn(15) + 5
		tests = append(tests, randomTree(rng, n))
	}
	tests = append(tests, randomTree(rng, 18))
	tests = append(tests, randomTree(rng, 20))
	return tests
}

func buildTest(n int, edges [][2]int) testCase {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d\n", n))
	for _, e := range edges {
		b.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return testCase{input: b.String(), n: n}
}

func randomTree(rng *rand.Rand, n int) testCase {
	if n < 2 {
		n = 2
	}
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		u := i
		v := rng.Intn(i-1) + 1
		edges = append(edges, [2]int{u, v})
	}
	return buildTest(n, edges)
}
