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
		fmt.Fprintln(os.Stderr, "usage: go run verifierE2.go /path/to/binary")
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
		want, err := parseOutput(wantOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d: %v\nOutput:\n%s\n", idx+1, err, wantOut)
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

		if len(got) != len(want) {
			fmt.Fprintf(os.Stderr, "test %d: expected %d answers, got %d\nInput:\n%s\n", idx+1, len(want), len(got), tc.input)
			os.Exit(1)
		}
		for i := range want {
			if want[i] != got[i] {
				fmt.Fprintf(os.Stderr, "test %d mismatch on query %d: expected %s, got %s\nInput:\n%s\nCandidate output:\n%s\n", idx+1, i+1, want[i], got[i], tc.input, gotOut)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	candidates := []string{
		"1846E2.go",
		filepath.Join("1000-1999", "1800-1899", "1840-1849", "1846", "1846E2.go"),
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("could not find 1846E2.go")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref1846E2_%d.bin", time.Now().UnixNano()))
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

func parseOutput(out string) ([]string, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return nil, fmt.Errorf("empty output")
	}
	for i, f := range fields {
		u := strings.ToUpper(f)
		if u != "YES" && u != "NO" {
			return nil, fmt.Errorf("invalid token %q (expected YES/NO)", f)
		}
		fields[i] = u
	}
	return fields, nil
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests,
		buildTest([]int64{7, 8, 9, 10}),
		buildTest([]int64{1, 2, 3, 4, 5, 1000000000000}),
	)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 40 {
		t := rng.Intn(10) + 1
		arr := make([]int64, t)
		for i := 0; i < t; i++ {
			arr[i] = rng.Int63n(1_000_000_000_000)
		}
		tests = append(tests, buildTest(arr))
	}
	return tests
}

func buildTest(nums []int64) testCase {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d\n", len(nums)))
	for _, v := range nums {
		b.WriteString(fmt.Sprintf("%d\n", v))
	}
	return testCase{input: b.String()}
}
