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
	x int64
	y int64
	z int64
	k int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refPath, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refPath)

	tests := buildTests()
	input := buildInput(tests)

	refOut, err := runProgram(refPath, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}
	expected, err := parseOutput(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference output parse error: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}
	got, err := parseOutput(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate output parse error: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	for i := range tests {
		if expected[i] != got[i] {
			fmt.Fprintf(os.Stderr, "Mismatch on test %d: expected %d, got %d\nInput used:\n%s\n", i+1, expected[i], got[i], input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2035E-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), "2035E.go")
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
	add := func(x, y, z, k int64) {
		tests = append(tests, testCase{x: x, y: y, z: z, k: k})
	}
	add(1, 1, 1, 1)
	add(2, 3, 5, 5)
	add(10, 20, 40, 5)
	add(1, 60, 100, 10)
	add(60, 1, 100, 10)
	add(1, 100000000, 100000000, 1)
	add(100000000, 1, 100000000, 100000000)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 200 {
		tests = append(tests, randomCase(rng))
	}
	return tests
}

func randomCase(rng *rand.Rand) testCase {
	x := int64(rng.Intn(1_0000_0000) + 1)
	y := int64(rng.Intn(1_0000_0000) + 1)
	z := int64(rng.Intn(1_0000_0000) + 1)
	k := int64(rng.Intn(1_0000_0000) + 1)
	return testCase{x: x, y: y, z: z, k: k}
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", tc.x, tc.y, tc.z, tc.k))
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
		if val < 0 {
			return nil, fmt.Errorf("negative cost %d at position %d", val, i+1)
		}
		res[i] = val
	}
	return res, nil
}
