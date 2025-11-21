package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type testCase struct {
	n    int
	bits string
}

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	out := filepath.Join(dir, "oracleB1")
	cmd := exec.Command("go", "build", "-o", out, "1002B1.go")
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v\n%s", err, string(output))
	}
	return out, nil
}

func runProgram(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func parseOutput(out string) (int, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	switch fields[0] {
	case "0":
		return 0, nil
	case "1":
		return 1, nil
	default:
		return 0, fmt.Errorf("invalid output %q", fields[0])
	}
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 2, bits: "00"},
		{n: 2, bits: "01"},
		{n: 2, bits: "10"},
		{n: 3, bits: "000"},
		{n: 3, bits: "001"},
		{n: 3, bits: "010"},
		{n: 3, bits: "100"},
		{n: 4, bits: "0000"},
		{n: 4, bits: "0001"},
		{n: 8, bits: "10000000"},
	}
}

func randomTests(count int) []testCase {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, count)
	for i := 0; i < count; i++ {
		n := rnd.Intn(7) + 2 // 2..8
		bits := make([]byte, n)
		for j := 0; j < n; j++ {
			if rnd.Intn(2) == 0 {
				bits[j] = '0'
			} else {
				bits[j] = '1'
			}
		}
		tests = append(tests, testCase{n: n, bits: string(bits)})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	tests := deterministicTests()
	tests = append(tests, randomTests(200)...)

	for idx, tc := range tests {
		input := fmt.Sprintf("%d\n%s\n", tc.n, tc.bits)

		expOut, err := runProgram(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle error: %v\n", idx+1, err)
			os.Exit(1)
		}
		expVal, err := parseOutput(expOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle output invalid: %v\n", idx+1, err)
			os.Exit(1)
		}

		gotOut, err := runProgram(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotVal, err := parseOutput(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		if gotVal != expVal {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d (n=%d bits=%s)\n", idx+1, expVal, gotVal, tc.n, tc.bits)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
