package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesF = `100
8
8
2
6
10
9
8
6
9
7
5
10
4
6
4
3
6
10
4
6
3
3
7
9
10
3
7
8
7
5
10
9
9
10
6
2
10
2
3
8
2
9
7
5
7
3
5
5
5
4
10
9
3
3
7
10
9
3
6
10
6
3
10
7
10
5
10
6
9
3
8
7
5
6
4
5
4
2
6
9
3
3
4
4
2
3
10
8
10
6
10
5
5
8
6
9
9
7
3
7`

type testCase struct {
	n int
}

func parseTests() ([]testCase, error) {
	reader := strings.NewReader(testcasesF)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, err
	}
	tests := make([]testCase, t)
	for i := 0; i < t; i++ {
		if _, err := fmt.Fscan(reader, &tests[i].n); err != nil {
			return nil, err
		}
	}
	return tests, nil
}

func buildAllInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierF /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTests()
	if err != nil {
		fmt.Fprintln(os.Stderr, "parse error:", err)
		os.Exit(1)
	}

	allInput := buildAllInput(tests)
	rawOut, err := runCandidate(bin, allInput)
	if err != nil {
		fmt.Fprintln(os.Stderr, "runtime error:", err)
		os.Exit(1)
	}

	scan := bufio.NewScanner(strings.NewReader(rawOut))
	scan.Split(bufio.ScanWords)
	for i, tc := range tests {
		seen := make([]bool, tc.n+1)
		for j := 0; j < tc.n; j++ {
			if !scan.Scan() {
				fmt.Fprintf(os.Stderr, "missing output for test %d\n", i+1)
				os.Exit(1)
			}
			val, err := strconv.Atoi(scan.Text())
			if err != nil {
				fmt.Fprintf(os.Stderr, "bad output for test %d\n", i+1)
				os.Exit(1)
			}
			if val < 1 || val > tc.n {
				fmt.Fprintf(os.Stderr, "value out of range on test %d: %d\n", i+1, val)
				os.Exit(1)
			}
			if seen[val] {
				fmt.Fprintf(os.Stderr, "duplicate value on test %d: %d\n", i+1, val)
				os.Exit(1)
			}
			seen[val] = true
		}
		for v := 1; v <= tc.n; v++ {
			if !seen[v] {
				fmt.Fprintf(os.Stderr, "missing value %d on test %d\n", v, i+1)
				os.Exit(1)
			}
		}
	}
	if scan.Scan() {
		fmt.Fprintln(os.Stderr, "extra output detected")
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
