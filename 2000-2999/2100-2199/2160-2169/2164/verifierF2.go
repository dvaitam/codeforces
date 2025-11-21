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

const refSource = "2000-2999/2100-2199/2160-2169/2164/2164F2.go"

type testCase struct {
	n int
	a []int
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2164F2-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	return tmp.Name(), nil
}

func runProgram(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		abs, err := filepath.Abs(bin)
		if err != nil {
			return "", err
		}
		cmd = exec.Command("go", "run", abs)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 1, a: []int{1}},
		{n: 2, a: []int{1, 2}},
		{n: 3, a: []int{2, 4, 6}},
		{n: 4, a: []int{3, 3, 3, 3}},
	}
}

func randomTest(rng *rand.Rand, n int) testCase {
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(10) + 1
	}
	return testCase{n: n, a: a}
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func parseOutput(out string) (string, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return "", fmt.Errorf("empty output")
	}
	return fields[0], nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF2.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	totalN := 0
	for len(tests) < 40 && totalN < 2000 {
		n := rng.Intn(35) + 5
		tests = append(tests, randomTest(rng, n))
		totalN += n
	}

	for idx, tc := range tests {
		inp := buildInput(tc)

		wantOut, err := runProgram(refBin, inp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		want, err := parseOutput(wantOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		gotOut, err := runProgram(candidate, inp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := parseOutput(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		if want != got {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %s got %s\nn=%d a=%v\n", idx+1, want, got, tc.n, tc.a)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
