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

type testCase struct {
	n int
	s string
}

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	out := filepath.Join(dir, "oracleA2")
	cmd := exec.Command("go", "build", "-o", out, "1184A2.go")
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

func parseOutput(out string) (int64, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q", fields[0])
	}
	if len(fields) > 1 {
		return 0, fmt.Errorf("extra tokens in output: %v", fields[1:])
	}
	return val, nil
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 1, s: "0"},
		{n: 1, s: "1"},
		{n: 2, s: "00"},
		{n: 2, s: "01"},
		{n: 4, s: "1010"},
		{n: 6, s: "010101"},
	}
}

func randomTests(count int) []testCase {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, count)
	for i := 0; i < count; i++ {
		var n int
		switch {
		case i%20 == 0:
			n = rnd.Intn(200000) + 1
		case i%5 == 0:
			n = rnd.Intn(10000) + 1
		default:
			n = rnd.Intn(200) + 1
		}
		sb := strings.Builder{}
		for j := 0; j < n; j++ {
			if rnd.Intn(2) == 0 {
				sb.WriteByte('0')
			} else {
				sb.WriteByte('1')
			}
		}
		tests = append(tests, testCase{n: n, s: sb.String()})
	}
	// add edge cases: alternating and all zeros/ones for large n
	n := 200000
	zeros := strings.Repeat("0", n)
	ones := strings.Repeat("1", n)
	alt := make([]byte, n)
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			alt[i] = '0'
		} else {
			alt[i] = '1'
		}
	}
	tests = append(tests, testCase{n: n, s: zeros})
	tests = append(tests, testCase{n: n, s: ones})
	tests = append(tests, testCase{n: n, s: string(alt)})
	return tests
}

func formatInput(tc testCase) string {
	return fmt.Sprintf("%d\n%s\n", tc.n, tc.s)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA2.go /path/to/binary")
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
		input := formatInput(tc)

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
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d (n=%d)\n", idx+1, expVal, gotVal, tc.n)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
