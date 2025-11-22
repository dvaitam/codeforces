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
	t     int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
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
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\nInput:\n%s", idx+1, err, tc.input)
			os.Exit(1)
		}
		expect, err := parseOutputs(refOut, tc.t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d: %v\nInput:\n%sOutput:\n%s", idx+1, err, tc.input, refOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\nInput:\n%s", idx+1, err, tc.input)
			os.Exit(1)
		}
		got, err := parseOutputs(gotOut, tc.t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d: %v\nInput:\n%sOutput:\n%s", idx+1, err, tc.input, gotOut)
			os.Exit(1)
		}

		if len(expect) != len(got) {
			fmt.Fprintf(os.Stderr, "test %d: output count mismatch, expected %d got %d\nInput:\n%sOutput:\n%s", idx+1, len(expect), len(got), tc.input, gotOut)
			os.Exit(1)
		}
		for i := range expect {
			if expect[i] != got[i] {
				fmt.Fprintf(os.Stderr, "test %d: mismatch at case %d, expected %d got %d\nInput:\n%sOutput:\n%s", idx+1, i+1, expect[i], got[i], tc.input, gotOut)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	candidates := []string{
		"2067C.go",
		filepath.Join("2000-2999", "2000-2099", "2060-2069", "2067", "2067C.go"),
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("could not find 2067C.go")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref2067C_%d.bin", time.Now().UnixNano()))
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

func parseOutputs(out string, t int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d integers, got %d", t, len(fields))
	}
	res := make([]int, t)
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", f)
		}
		if v < 0 {
			return nil, fmt.Errorf("negative answer %d", v)
		}
		res[i] = v
	}
	return res, nil
}

func generateTests() []testCase {
	tests := make([]testCase, 0, 60)
	tests = append(tests, sampleTest())
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 50; i++ {
		tests = append(tests, randomTest(rng))
	}
	return tests
}

func sampleTest() testCase {
	input := "16\n5\n51\n60\n61\n77\n123456\n8910000000\n200\n2300\n1977\n9898989\n8680\n80000\n196701\n590\n"
	return testCase{input: input, t: 16}
}

func randomTest(rng *rand.Rand) testCase {
	t := rng.Intn(8) + 1 // 1..8
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := rng.Intn(1_000_000_000-10+1) + 10
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
	}
	return testCase{input: sb.String(), t: t}
}
