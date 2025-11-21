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

const (
	refSource331C2 = "331C2.go"
	refBinary331C2 = "ref331C2.bin"
	totalTests     = 200
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC2.go /path/to/binary")
		return
	}
	candidate := os.Args[1]

	ref, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		return
	}
	defer os.Remove(ref)

	tests := generateTests()
	input := formatInput(tests)

	refOut, err := runProgram(ref, input)
	if err != nil {
		fmt.Printf("reference runtime error: %v\n", err)
		return
	}
	expected, err := parseOutput(refOut, len(tests))
	if err != nil {
		fmt.Printf("reference output parse error: %v\noutput:\n%s", err, refOut)
		return
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Printf("candidate runtime error: %v\n", err)
		return
	}
	got, err := parseOutput(candOut, len(tests))
	if err != nil {
		fmt.Printf("candidate output parse error: %v\noutput:\n%s", err, candOut)
		return
	}

	for i := range tests {
		if expected[i] != got[i] {
			fmt.Printf("Mismatch on test %d: expected %d, got %d\n", i+1, expected[i], got[i])
			fmt.Println("Input used:")
			fmt.Println(string(input))
			return
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinary331C2, refSource331C2)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinary331C2), nil
}

func runProgram(path string, input []byte) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(out string, t int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d numbers, got %d", t, len(fields))
	}
	res := make([]int64, t)
	for i, tok := range fields {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("token %d: %v", i+1, err)
		}
		res[i] = val
	}
	return res, nil
}

func formatInput(tests []int64) []byte {
	var sb strings.Builder
	sb.WriteString("1\n")
	for _, v := range tests {
		sb.WriteString(fmt.Sprintf("%d\n", v))
	}
	return []byte(strings.Join(strings.Fields(sb.String()), "\n") + "\n")
}

func generateTests() []int64 {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := []int64{0, 1, 2, 9, 10, 19, 55, 100, 12345, 999999999999}
	for len(tests) < totalTests {
		switch rnd.Intn(5) {
		case 0:
			tests = append(tests, int64(rnd.Intn(1_000_000)))
		case 1:
			tests = append(tests, rnd.Int63n(1_000_000_000_000))
		case 2:
			tests = append(tests, rnd.Int63n(1_000_000_000_000_000))
		default:
			tests = append(tests, rnd.Int63n(1_000_000_000_000_000_000))
		}
	}
	return tests
}
