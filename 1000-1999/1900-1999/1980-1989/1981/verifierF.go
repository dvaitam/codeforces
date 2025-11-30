package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesF = `100
2
3 5
1
2
5 4
1
2
5 4
1
2
5 2
1
2
2 5
1
2
5 5
1
2
3 2
1
2
2 4
1
2
4 5
1
2
5 5
1
2
5 1
1
2
2 5
1
2
5 3
1
2
4 5
1
2
2 3
1
2
5 3
1
2
1 3
1
2
4 4
1
2
3 1
1
2
2 2
1
2
2 1
1
2
1 5
1
2
4 3
1
2
1 4
1
2
2 2
1
2
4 2
1
2
4 2
1
2
1 1
1
2
2 4
1
2
4 2
1
2
4 1
1
2
2 1
1
2
4 3
1
2
3 1
1
2
3 4
1
2
3 4
1
2
3 4
1
2
1 4
1
2
2 3
1
2
5 2
1
2
4 3
1
2
5 3
1
2
4 2
1
2
5 1
1
2
3 5
1
2
3 3
1
2
1 5
1
2
1 3
1
2
1 3
1
2
4 5
1
2
5 2
1
2
3 2
1
2
5 1
1
2
3 2
1
2
5 5
1
2
4 3
1
2
3 4
1
2
2 1
1
2
2 1
1
2
1 5
1
2
1 3
1
2
4 1
1
2
1 5
1
2
4 5
1
2
3 5
1
2
1 2
1
2
5 4
1
2
2 2
1
2
1 3
1
2
1 2
1
2
1 4
1
2
5 5
1
2
4 1
1
2
1 1
1
2
4 3
1
2
1 1
1
2
5 3
1
2
5 3
1
2
4 1
1
2
5 2
1
2
3 2
1
2
2 4
1
2
4 1
1
2
1 1
1
2
4 1
1
2
2 3
1
2
1 5
1
2
1 5
1
2
3 4
1
2
2 3
1
2
3 4
1
2
2 1
1
2
1 5
1
2
5 2
1
2
2 4
1
2
4 2
1
2
4 2
1
2
3 5
1
2
3 1
1
2
3 1
1`

type testCase struct {
	n  int
	a1 int
	a2 int
}

func mex(a, b int) int {
	used := map[int]bool{a: true, b: true}
	for i := 1; ; i++ {
		if !used[i] {
			return i
		}
	}
}

func parseTests() ([]testCase, error) {
	fields := strings.Fields(testcasesF)
	pos := 0
	readInt := func() (int, error) {
		if pos >= len(fields) {
			return 0, fmt.Errorf("unexpected EOF")
		}
		v, err := strconv.Atoi(fields[pos])
		pos++
		return v, err
	}
	t, err := readInt()
	if err != nil {
		return nil, err
	}
	tests := make([]testCase, t)
	for i := 0; i < t; i++ {
		n, err := readInt()
		if err != nil {
			return nil, err
		}
		a1, err := readInt()
		if err != nil {
			return nil, err
		}
		a2, err := readInt()
		if err != nil {
			return nil, err
		}
		// parent line
		if _, err := readInt(); err != nil {
			return nil, err
		}
		tests[i] = testCase{n: n, a1: a1, a2: a2}
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
		sb.WriteString(strconv.Itoa(tc.a1))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(tc.a2))
		sb.WriteByte('\n')
		// parent line
		if tc.n >= 2 {
			sb.WriteString("1\n")
		} else {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
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
	allOutput, err := runCandidate(bin, allInput)
	if err != nil {
		fmt.Fprintln(os.Stderr, "runtime error:", err)
		os.Exit(1)
	}

	outFields := strings.Fields(allOutput)
	if len(outFields) != len(tests) {
		fmt.Fprintf(os.Stderr, "expected %d outputs, got %d\n", len(tests), len(outFields))
		os.Exit(1)
	}
	for i, tc := range tests {
		want := mex(tc.a1, tc.a2)
		got, err := strconv.Atoi(outFields[i])
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: non-integer output %q\n", i+1, outFields[i])
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected: %d\ngot: %d\n", i+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
