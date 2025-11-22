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
		"2075E.go",
		filepath.Join("2000-2999", "2000-2099", "2070-2079", "2075", "2075E.go"),
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("could not find 2075E.go")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref2075E_%d.bin", time.Now().UnixNano()))
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

func parseOutputs(out string, t int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d integers, got %d", t, len(fields))
	}
	res := make([]int64, t)
	for i, f := range fields {
		v, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", f)
		}
		res[i] = v
	}
	return res, nil
}

func generateTests() []testCase {
	tests := make([]testCase, 0, 50)
	tests = append(tests, sampleTest())
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 40; i++ {
		tests = append(tests, randomTest(rng))
	}
	return tests
}

func sampleTest() testCase {
	input := "6\n2 2 2 2\n2 3 4 5\n5 7 4 3\n1337 42 1337 42\n4 2 13 3\n753687090 2 536370902 536390912\n"
	return testCase{input: input, t: 6}
}

func randomTest(rng *rand.Rand) testCase {
	t := rng.Intn(5) + 1 // 1..5 cases
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := randRange(rng, 2, 200)
		m := randRange(rng, 2, 200)
		A := randRange(rng, 2, 1<<20)
		B := randRange(rng, 2, 1<<20)
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", n, m, A, B))
	}
	return testCase{input: sb.String(), t: t}
}

func randRange(rng *rand.Rand, lo, hi int64) int64 {
	return lo + rng.Int63n(hi-lo+1)
}
