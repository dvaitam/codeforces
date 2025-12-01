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

const refSource = "./2140B.go"

type testCase struct {
	x    int64
	name string
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

	// Run reference to ensure it executes successfully.
	if _, err := runProgram(refBin, input); err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	if err := validateCandidate(candOut, tests); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	outPath := "./ref_2140B.bin"
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
	add := func(name string, x int64) {
		tests = append(tests, testCase{name: name, x: x})
	}

	// Edge and sample-like values
	add("x1", 1)
	add("x2", 2)
	add("x8", 8)
	add("x42", 42)
	add("x100", 100)
	add("x99999999", 99_999_999)
	add("x1e8minus1", 100_000_000-1)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 180 {
		x := rng.Int63n(100_000_000-1) + 1 // 1..1e8-1
		add(fmt.Sprintf("random_%d", len(tests)), x)
	}
	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d\n", tc.x))
	}
	return sb.String()
}

func validateCandidate(out string, tests []testCase) error {
	fields := strings.Fields(out)
	if len(fields) != len(tests) {
		return fmt.Errorf("expected %d outputs, got %d", len(tests), len(fields))
	}
	for i, tc := range tests {
		y, err := strconv.ParseInt(fields[i], 10, 64)
		if err != nil {
			return fmt.Errorf("case %d (%s): invalid integer %q", i+1, tc.name, fields[i])
		}
		if y <= 0 || y >= 1_000_000_000 {
			return fmt.Errorf("case %d (%s): y=%d out of range", i+1, tc.name, y)
		}
		if !checkPair(tc.x, y) {
			return fmt.Errorf("case %d (%s): x=%d y=%d fails divisibility", i+1, tc.name, tc.x, y)
		}
	}
	return nil
}

func checkPair(x, y int64) bool {
	lenY := digits(y)
	pow := int64(1)
	for i := int64(0); i < lenY; i++ {
		pow *= 10
	}
	concat := x*pow + y
	return concat%(x+y) == 0
}

func digits(v int64) int64 {
	if v == 0 {
		return 1
	}
	cnt := int64(0)
	for v > 0 {
		cnt++
		v /= 10
	}
	return cnt
}
