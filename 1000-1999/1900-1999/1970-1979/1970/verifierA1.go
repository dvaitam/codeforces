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
		fmt.Fprintln(os.Stderr, "usage: go run verifierA1.go /path/to/binary")
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
		"1970A1.go",
		filepath.Join("1000-1999", "1900-1999", "1970-1979", "1970", "1970A1.go"),
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("could not find 1970A1.go")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref1970A1_%d.bin", time.Now().UnixNano()))
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

func generateTests() []testCase {
	var tests []testCase
	// Only balanced sequences are valid per problem constraints.
	tests = append(tests,
		buildTest("()"),
		buildTest("(())"),
		buildTest("()()"),
		buildTest("((()))"),
		buildTest("()(())"),
	)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 60 {
		k := rng.Intn(20) + 1 // length = 2k
		tests = append(tests, buildTest(randomBalanced(rng, k)))
	}
	return tests
}

func buildTest(s string) testCase {
	return testCase{input: fmt.Sprintf("%s\n", s)}
}

// randomBalanced generates a uniformly-random balanced parentheses sequence of length 2k.
func randomBalanced(rng *rand.Rand, k int) string {
	n := 2 * k
	var b strings.Builder
	b.Grow(n)
	open := 0
	for i := 0; i < n; i++ {
		remaining := n - i
		// Must close: all remaining characters must be ')' to balance.
		if open == remaining {
			b.WriteByte(')')
			open--
			continue
		}
		// Must open: no open brackets to close yet.
		if open == 0 {
			b.WriteByte('(')
			open++
			continue
		}
		if rng.Intn(2) == 0 {
			b.WriteByte('(')
			open++
		} else {
			b.WriteByte(')')
			open--
		}
	}
	return b.String()
}
