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

const refSource = "./2152B.go"

type testCase struct {
	n      int64
	rK, cK int64
	rD, cD int64
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := args[0]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}
	refAns, err := parseOutput(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n%s", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	candAns, err := parseOutput(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\n%s", err, candOut)
		os.Exit(1)
	}

	for i, tc := range tests {
		if candAns[i] != refAns[i] {
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %d, got %d\ninput: %v\n", i+1, refAns[i], candAns[i], tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	outPath := "./ref_2152B.bin"
	cmd := exec.Command("go", "build", "-o", outPath, refSource)
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
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func buildTests() []testCase {
	var tests []testCase
	add := func(n, rK, cK, rD, cD int64) {
		if rK == rD && cK == cD {
			cD = (cD + 1) % (n + 1)
		}
		tests = append(tests, testCase{n, rK, cK, rD, cD})
	}

	// simple cases
	add(2, 0, 0, 1, 1)
	add(2, 1, 1, 0, 1)
	add(3, 1, 1, 0, 1)
	add(3, 0, 2, 2, 2)
	add(4, 0, 0, 3, 3)
	add(4, 3, 3, 0, 0)

	// shared row/column
	add(10, 5, 5, 5, 2)
	add(10, 5, 5, 2, 5)
	add(10, 0, 5, 0, 0)
	add(10, 5, 0, 0, 0)

	// corner cases
	add(1000000000, 0, 0, 1000000000, 1000000000)
	add(1000000000, 1000000000, 1000000000, 0, 0)
	add(1000000000, 500000000, 500000000, 1000000000, 0)
	add(1000000000, 0, 1000000000, 1000000000, 1000000000)

	// random cases
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 150 {
		n := int64(rng.Intn(1_000_000_000) + 1)
		rK := int64(rng.Intn(int(n) + 1))
		cK := int64(rng.Intn(int(n) + 1))
		rD := int64(rng.Intn(int(n) + 1))
		cD := int64(rng.Intn(int(n) + 1))
		if rK == rD && cK == cD {
			continue
		}
		add(n, rK, cK, rD, cD)
	}
	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d %d %d %d\n", tc.n, tc.rK, tc.cK, tc.rD, tc.cD))
	}
	return sb.String()
}

func parseOutput(out string, expected int) ([]int64, error) {
	lines := strings.Fields(out)
	if len(lines) != expected {
		return nil, fmt.Errorf("expected %d outputs, got %d", expected, len(lines))
	}
	ans := make([]int64, expected)
	for i, s := range lines {
		val, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", s)
		}
		ans[i] = val
	}
	return ans, nil
}
