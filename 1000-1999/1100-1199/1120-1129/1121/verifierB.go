package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	n int
	a []int
}

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	out := filepath.Join(dir, "oracleB")
	cmd := exec.Command("go", "build", "-o", out, "1121B.go")
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
	val, err := strconv.Atoi(fields[0])
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
		{n: 2, a: []int{1, 2}},
		{n: 4, a: []int{1, 3, 5, 7}},
		{n: 5, a: []int{1, 2, 3, 4, 5}},
		{n: 6, a: []int{1, 5, 3, 7, 9, 11}},
		{n: 6, a: []int{6, 1, 3, 2, 5, 4}},
	}
}

func randomTests(count int) []testCase {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, count)
	for i := 0; i < count; i++ {
		n := rnd.Intn(999) + 2 // 2..1000
		values := make([]int, n)
		for j := 0; j < n; j++ {
			values[j] = rnd.Intn(100000) + 1
		}
		sort.Ints(values)
		for j := 1; j < n; j++ {
			if values[j] == values[j-1] {
				values[j]++
			}
		}
		tests = append(tests, testCase{n: n, a: values})
	}
	return tests
}

func formatInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
