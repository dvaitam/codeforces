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
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
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
		if expected[i] != got[i] {
			fmt.Fprintf(os.Stderr, "Mismatch on test %d: expected %d, got %d\nInput:\n%s\n", i+1, expected[i], got[i], input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2006C-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), "2006C.go")
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
	add := func(arr []int64) {
		if sumN+len(arr) > 400000 {
			return
		}
		cpy := append([]int64(nil), arr...)
		tests = append(tests, testCase{n: len(arr), a: cpy})
		sumN += len(arr)
	}

	add([]int64{1})
	add([]int64{1, 1})
	add([]int64{1, 2, 3})
	add([]int64{3, 6, 10})
	add([]int64{1000000000, 1, 4, 5, 1, 4})
	add([]int64{12, 70, 130, 90, 90, 90, 108})
	add([]int64{6, 12, 500, 451, 171, 193, 193})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 120 && sumN < 400000 {
		n := rng.Intn(10) + 1
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			arr[i] = int64(rng.Intn(20) + 1)
		}
		add(arr)
	}

	for len(tests) < 200 && sumN < 400000 {
		n := rng.Intn(200) + 50
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			arr[i] = int64(rng.Intn(1000) + 1)
		}
		add(arr)
	}

	if sumN < 400000 {
		n := 400000 - sumN
		if n > 0 {
			arr := make([]int64, n)
			for i := range arr {
				arr[i] = int64(rng.Intn(1_000_000_000) + 1)
			}
			add(arr)
		}
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
	}
	return sb.String()
}

func parseOutput(out string, t int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d outputs, got %d", t, len(fields))
	}
	res := make([]int64, t)
	for i, tok := range fields {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q at position %d", tok, i+1)
		}
		res[i] = val
	}
	return res, nil
}
