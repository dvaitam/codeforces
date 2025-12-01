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
)

const refSourceA = "./2001A.go"

type testCase struct {
	arr []int
}

func (tc testCase) input() string {
	var b strings.Builder
	fmt.Fprintln(&b, 1)
	fmt.Fprintln(&b, len(tc.arr))
	for i, v := range tc.arr {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(strconv.Itoa(v))
	}
	b.WriteByte('\n')
	return b.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, tc := range tests {
		input := tc.input()
		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		refAns, err := parseAnswer(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s\nstdout/stderr:\n%s\n", idx+1, err, input, candOut)
			os.Exit(1)
		}
		candAns, err := parseAnswer(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid output on test %d: %v\noutput:\n%s\n", idx+1, err, candOut)
			os.Exit(1)
		}

		if candAns != refAns {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%sreference: %d\ncandidate: %d\n", idx+1, input, refAns, candAns)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2001A-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSourceA))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseAnswer(out string) (int, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q", fields[0])
	}
	if val < 0 {
		return 0, fmt.Errorf("negative answer %d", val)
	}
	return val, nil
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests, testCase{arr: []int{1}})
	tests = append(tests, testCase{arr: []int{1, 1, 1}})
	tests = append(tests, testCase{arr: []int{1, 2}})
	tests = append(tests, testCase{arr: []int{3, 1, 2}})
	tests = append(tests, testCase{arr: []int{5, 4, 3, 2, 1}})
	tests = append(tests, testCase{arr: []int{1, 1, 2, 2, 3, 3}})
	tests = append(tests, testCase{arr: []int{1, 2, 3, 4, 5, 6}})
	tests = append(tests, testCase{arr: []int{2, 2, 2, 3, 3, 3, 4, 4}})

	rng := rand.New(rand.NewSource(2001))
	for len(tests) < 200 {
		n := rng.Intn(100) + 1
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = rng.Intn(n) + 1
		}
		tests = append(tests, testCase{arr: arr})
	}

	return tests
}
