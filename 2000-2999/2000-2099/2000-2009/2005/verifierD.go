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
	n int
	a []int64
	b []int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refPath, err := buildReference()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.Remove(refPath)

	tests := buildTests()
	input := buildInput(tests)

	refOut, err := runProgram(refPath, input)
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
		if expected[i].best != got[i].best || expected[i].ways != got[i].ways {
			fmt.Fprintf(os.Stderr, "Mismatch on test %d: expected (%d,%d), got (%d,%d)\nInput:\n%s\n",
				i+1, expected[i].best, expected[i].ways, got[i].best, got[i].ways, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2005D-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), "2005D.go")
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
	sumN := 0
	add := func(a, b []int64) {
		if len(a) != len(b) {
			panic("length mismatch")
		}
		if sumN+len(a) > 500000 {
			return
		}
		na := append([]int64(nil), a...)
		nb := append([]int64(nil), b...)
		tests = append(tests, testCase{n: len(a), a: na, b: nb})
		sumN += len(a)
	}

	add([]int64{1}, []int64{1})
	add([]int64{2, 4}, []int64{3, 9})
	add([]int64{5, 10, 15}, []int64{6, 12, 18})
	add([]int64{1, 2, 3, 4}, []int64{4, 3, 2, 1})
	add([]int64{1000000000, 1000000000}, []int64{999999937, 999999937})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 50 && sumN < 500000 {
		n := rng.Intn(5) + 1
		a := make([]int64, n)
		b := make([]int64, n)
		for i := range a {
			a[i] = int64(rng.Intn(20) + 1)
			b[i] = int64(rng.Intn(20) + 1)
		}
		add(a, b)
	}

	for len(tests) < 150 && sumN < 500000 {
		n := rng.Intn(200) + 50
		a := make([]int64, n)
		b := make([]int64, n)
		for i := range a {
			a[i] = int64(rng.Intn(1_000_000_000) + 1)
			b[i] = int64(rng.Intn(1_000_000_000) + 1)
		}
		add(a, b)
	}

	if sumN < 500000 {
		n := 500000 - sumN
		if n < 1 {
			n = 1
		}
		a := make([]int64, n)
		b := make([]int64, n)
		for i := range a {
			a[i] = int64(rng.Intn(1_000_000_000) + 1)
			b[i] = int64(rng.Intn(1_000_000_000) + 1)
		}
		add(a, b)
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
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		for i, v := range tc.b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

type result struct {
	best int64
	ways int64
}

func parseOutput(out string, t int) ([]result, error) {
	fields := strings.Fields(out)
	if len(fields) != 2*t {
		return nil, fmt.Errorf("expected %d numbers, got %d", 2*t, len(fields))
	}
	ans := make([]result, t)
	for i := 0; i < t; i++ {
		best, err := strconv.ParseInt(fields[2*i], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid best value %q at position %d", fields[2*i], 2*i+1)
		}
		ways, err := strconv.ParseInt(fields[2*i+1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid ways value %q at position %d", fields[2*i+1], 2*i+2)
		}
		ans[i] = result{best: best, ways: ways}
	}
	return ans, nil
}
