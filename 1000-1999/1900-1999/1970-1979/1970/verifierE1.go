package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type testCase struct {
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE1.go /path/to/binary")
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

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := generateTests(rng)
	for idx, tc := range tests {
		wantOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\nInput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		want := strings.TrimSpace(wantOut)

		gotOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\nInput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		got := strings.TrimSpace(gotOut)

		if want != got {
			fmt.Fprintf(os.Stderr, "test %d mismatch:\nInput:\n%s\nExpected:\n%s\nGot:\n%s\n", idx+1, tc.input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	candidates := []string{
		"1970E1.go",
		filepath.Join("1000-1999", "1900-1999", "1970-1979", "1970", "1970E1.go"),
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("could not find 1970E1.go")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref1970E1_%d.bin", time.Now().UnixNano()))
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

func generateTests(rng *rand.Rand) []testCase {
	var tests []testCase

	// Fixed edge cases
	tests = append(tests,
		makeTest(1, 1, []int{1}, []int{0}),
		makeTest(1, 1, []int{0}, []int{1}),
		makeTest(1, 1, []int{1}, []int{1}),
		makeTest(2, 1, []int{1, 2}, []int{3, 0}),
		makeTest(2, 3, []int{1, 1}, []int{1, 1}),
		makeTest(1, 1000, []int{1000}, []int{1000}),
	)

	// Random tests with small constraints
	for len(tests) < 60 {
		m := rng.Intn(5) + 1
		n := rng.Intn(10) + 1
		s := make([]int, m)
		l := make([]int, m)
		for i := 0; i < m; i++ {
			s[i] = rng.Intn(10)
			l[i] = rng.Intn(10)
		}
		tests = append(tests, makeTest(m, n, s, l))
	}

	// Random tests with larger constraints
	for len(tests) < 80 {
		m := rng.Intn(10) + 1
		n := rng.Intn(50) + 1
		s := make([]int, m)
		l := make([]int, m)
		for i := 0; i < m; i++ {
			s[i] = rng.Intn(1001)
			l[i] = rng.Intn(1001)
		}
		tests = append(tests, makeTest(m, n, s, l))
	}

	return tests
}

func makeTest(m, n int, s, l []int) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", m, n)
	for i, v := range s {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	for i, v := range l {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return testCase{input: sb.String()}
}
