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
	refSourceA = "265A.go"
	refBinaryA = "ref265A.bin"
	totalTests = 200
)

type testCase struct {
	s string
	t string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refPath, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refPath)

	tests := generateTests()
	for i, tc := range tests {
		input := formatInput(tc)

		refOut, err := runProgram(refPath, input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			printInput(input)
			os.Exit(1)
		}
		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			printInput(input)
			os.Exit(1)
		}

		refVal, err := parseOutput(refOut)
		if err != nil {
			fmt.Printf("failed to parse reference output on test %d: %v\noutput:\n%s\n", i+1, err, refOut)
			printInput(input)
			os.Exit(1)
		}
		candVal, err := parseOutput(candOut)
		if err != nil {
			fmt.Printf("failed to parse candidate output on test %d: %v\noutput:\n%s\n", i+1, err, candOut)
			printInput(input)
			os.Exit(1)
		}

		if refVal != candVal {
			fmt.Printf("test %d failed: expected %d, got %d\n", i+1, refVal, candVal)
			printInput(input)
			fmt.Println("Reference output:")
			fmt.Println(refOut)
			fmt.Println("Candidate output:")
			fmt.Println(candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinaryA, refSourceA)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinaryA), nil
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

func formatInput(tc testCase) []byte {
	return []byte(fmt.Sprintf("%s\n%s\n", tc.s, tc.t))
}

func parseOutput(out string) (int, error) {
	fields := strings.Fields(out)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected 1 integer, got %d", len(fields))
	}
	val, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q: %v", fields[0], err)
	}
	return val, nil
}

func generateTests() []testCase {
	tests := []testCase{
		{s: "R", t: "R"},
		{s: "RGB", t: "RGB"},
		{s: "RGB", t: "BBB"},
		{s: "RRRRR", t: "RRRRR"},
		{s: "B", t: "RRRRR"},
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < totalTests-2 {
		sLen := rnd.Intn(50) + 1
		tLen := rnd.Intn(50) + 1
		s := randomColorString(rnd, sLen)
		t := randomColorString(rnd, tLen)
		tests = append(tests, testCase{s: s, t: t})
	}

	tests = append(tests, testCase{
		s: strings.Repeat("R", 50),
		t: strings.Repeat("R", 50),
	})
	tests = append(tests, testCase{
		s: strings.Repeat("RGB", 17) + "R",
		t: strings.Repeat("BGR", 17) + "B",
	})

	return tests
}

func randomColorString(rnd *rand.Rand, length int) string {
	colors := []byte{'R', 'G', 'B'}
	res := make([]byte, length)
	for i := 0; i < length; i++ {
		res[i] = colors[rnd.Intn(len(colors))]
	}
	return string(res)
}

func printInput(in []byte) {
	fmt.Println("Input used:")
	fmt.Println(string(in))
}
