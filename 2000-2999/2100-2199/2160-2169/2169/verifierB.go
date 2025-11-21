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
	s string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := buildTestSuite(rng)
	input := serializeInput(tests)

	expected, err := runAndParse(refBin, input, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}

	got, err := runAndParse(candidate, input, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\n", err)
		os.Exit(1)
	}

	for i := range expected {
		if expected[i] != got[i] {
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %s got %s\n", i+1, expected[i], got[i])
			fmt.Fprintf(os.Stderr, "input string: %s\n", tests[i].s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	const refName = "./ref_2169B.bin"
	cmd := exec.Command("go", "build", "-o", refName, "2169B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return refName, nil
}

func runAndParse(target, input string, tests int) ([]string, error) {
	out, err := runProgram(target, input)
	if err != nil {
		return nil, err
	}
	lines := strings.Fields(out)
	if len(lines) != tests {
		return nil, fmt.Errorf("expected %d outputs, got %d (output: %q)", tests, len(lines), out)
	}
	return lines, nil
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
	return stdout.String(), nil
}

func buildTestSuite(rng *rand.Rand) []testCase {
	tests := deterministicTests()
	totalLen := 0
	for _, tc := range tests {
		totalLen += len(tc.s)
	}

	for len(tests) < 300 && totalLen < 300000 {
		length := rng.Intn(1000) + 1
		if totalLen+length > 300000 {
			length = 300000 - totalLen
		}
		tc := randomTest(rng, length)
		tests = append(tests, tc)
		totalLen += len(tc.s)
	}
	return tests
}

func deterministicTests() []testCase {
	return []testCase{
		{s: "*"},
		{s: "<"},
		{s: ">"},
		{s: "<<<<<"},
		{s: ">>>>>"},
		{s: "*****"},
		{s: "<><><>"},
		{s: "><><><"},
		{s: "<**>"},
		{s: "<*><*"},
		{s: "<*****>"},
		{s: "<<***>>"},
		{s: ">***<"},
		{s: "<>>><<<>"},
		{s: "<<<<>>>><<<"},
		{s: ">>>**<<<"},
	}
}

func randomTest(rng *rand.Rand, length int) testCase {
	chars := []byte("<>*")
	// sometimes bias to repeated characters
	result := make([]byte, length)
	mode := rng.Intn(5)
	for i := 0; i < length; i++ {
		switch mode {
		case 0:
			result[i] = '<'
		case 1:
			result[i] = '>'
		case 2:
			result[i] = '*'
		default:
			result[i] = chars[rng.Intn(len(chars))]
		}
	}
	return testCase{s: string(result)}
}

func serializeInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(tc.s)
		sb.WriteByte('\n')
	}
	return sb.String()
}
