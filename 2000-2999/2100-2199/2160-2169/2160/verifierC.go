package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type testCase struct {
	id    string
	input string
}

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/candidate")
		os.Exit(1)
	}
	target := os.Args[len(os.Args)-1]
	if target == "--" {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/candidate")
		os.Exit(1)
	}

	base := currentDir()
	refBin, err := buildReference(base)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for i, tc := range tests {
		exp, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on %s: %v\ninput:\n%s", tc.id, err, tc.input)
			os.Exit(1)
		}
		got, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on %s: %v\ninput:\n%s", tc.id, err, tc.input)
			os.Exit(1)
		}
		if normalize(exp) != normalize(got) {
			fmt.Fprintf(os.Stderr, "wrong answer on %s\nInput:\n%sExpected:\n%sGot:\n%s", tc.id, tc.input, exp, got)
			os.Exit(1)
		}
		if (i+1)%10 == 0 {
			fmt.Fprintf(os.Stderr, "validated %d/%d tests...\n", i+1, len(tests))
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func currentDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("cannot determine current file path")
	}
	return filepath.Dir(file)
}

func buildReference(dir string) (string, error) {
	out := filepath.Join(dir, "ref2160C.bin")
	cmd := exec.Command("go", "build", "-o", out, "2160C.go")
	cmd.Dir = dir
	if data, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("go build failed: %v\n%s", err, data)
	}
	return out, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func normalize(s string) string {
	lines := strings.Fields(strings.ToUpper(s))
	return strings.Join(lines, "\n")
}

func generateTests() []testCase {
	var tests []testCase
	deterministic := []int{0, 1, 2, 3, 8, 15, 16, 31, 32, 63}
	tests = append(tests, makeCase("deterministic", deterministic))

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 40; i++ {
		t := rng.Intn(20) + 1
		vals := make([]int, t)
		for j := 0; j < t; j++ {
			vals[j] = rng.Intn(1 << 30)
		}
		tests = append(tests, makeCase(fmt.Sprintf("rand-%02d", i+1), vals))
	}

	// stress near limits
	tests = append(tests, makeCase("near-limit", []int{
		(1 << 30) - 1,
		(1 << 30) - 2,
		(1 << 29),
		(1 << 29) - 1,
	}))
	return tests
}

func makeCase(id string, numbers []int) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(numbers))
	for _, v := range numbers {
		fmt.Fprintf(&sb, "%d\n", v)
	}
	return testCase{id: id, input: sb.String()}
}
