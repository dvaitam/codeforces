package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	n      int
	colors []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	ref, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := buildTests()
	input := buildInput(tests)

	refOut, err := runProgram(ref, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n%s", err, refOut)
		os.Exit(1)
	}
	expected, err := parseOutput(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference output parse error: %v\n%s", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n%s", err, candOut)
		os.Exit(1)
	}
	got, err := parseOutput(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate output parse error: %v\n%s", err, candOut)
		os.Exit(1)
	}

	for i := range tests {
		if expected[i] != got[i] {
			fmt.Fprintf(os.Stderr, "Mismatch on test %d: expected %d, got %d\ninput:\n%s\n", i+1, expected[i], got[i], input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2042B-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), "2042B.go")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), err
	}
	return stdout.String(), nil
}

func buildTests() []testCase {
	var tests []testCase
	add := func(colors []int) {
		tests = append(tests, testCase{n: len(colors), colors: append([]int(nil), colors...)})
	}
	add([]int{1})
	add([]int{1, 1})
	add([]int{1, 2})
	add([]int{1, 1, 1})
	add([]int{1, 2, 1})
	add([]int{1, 2, 3, 4})
	add([]int{1, 1, 2, 2})
	add([]int{1, 1, 1, 1})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 200 {
		n := rng.Intn(10) + 1
		colorCount := rng.Intn(n) + 1
		colors := make([]int, n)
		for i := 0; i < n; i++ {
			colors[i] = rng.Intn(colorCount) + 1
		}
		add(colors)
	}

	for len(tests) < 500 {
		n := rng.Intn(50) + 50
		colorCount := rng.Intn(n) + 1
		colors := make([]int, n)
		for i := 0; i < n; i++ {
			colors[i] = rng.Intn(colorCount) + 1
		}
		add(colors)
	}
	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for i, c := range tc.colors {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(c))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string, t int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d outputs, got %d", t, len(fields))
	}
	res := make([]int, t)
	for i, tok := range fields {
		val, err := strconv.Atoi(tok)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q at position %d", tok, i+1)
		}
		res[i] = val
	}
	return res, nil
}
