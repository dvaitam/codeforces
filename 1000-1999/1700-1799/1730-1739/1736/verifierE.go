package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	input    string
	expected string
}

const testcaseData = `6 9 12 17 1 15 8
7 2 6 4 12 16 8 13
6 4 19 8 1 7 14
4 6 13 6 3
3 20 20 15
3 5 1 1
3 7 6 6
4 11 7 18 7
3 7 13 10
2 12 14
3 5 9 3
4 10 20 19 1
6 11 3 10 12 10 16
7 11 6 16 16 6 2 9
2 12 13
2 18 14
4 13 19 1 15
2 6 20
3 4 8 15
4 17 12 17 9
8 15 4 19 12 10 2 14 3
3 11 17 20
4 5 11 9 18
2 10 11
4 6 3 5 10
5 6 2 3 20 18
5 2 8 20 12 9
5 14 5 2 2 16
4 7 5 19 5
7 14 4 6 14 12 5 2
8 14 10 5 15 20 6 17 15
5 11 16 9 10 16
5 5 4 13 18 6
7 16 11 6 3 16 9 17
8 18 17 12 3 12 19 2 10
4 18 9 16 9
8 10 11 6 19 1 16 18 9
4 9 15 10 17
7 12 12 9 12 14 12 6
8 15 12 11 17 5 17 6 7
8 12 16 10 3 14 6 20 19
6 14 10 20 18 9 1
3 6 19 15
6 6 8 6 2 16 8
3 2 5 4
4 6 16 7 18
2 14 15
4 13 20 3 19
3 8 12 1
4 13 9 14 4
7 18 12 2 18 20 10 4
4 18 17 11 19
4 12 5 14 14
8 19 18 12 15 5 6 20 13
6 16 7 5 20 3 12
8 1 13 4 11 19 20 18 5
4 19 13 14 14
3 16 10 16
7 13 13 6 20 20 9 10
5 9 14 1 11 10
5 10 5 16 1 4
7 20 15 8 10 2 5 13
2 16 18
6 9 8 16 2 8 16
4 5 10 10 16
6 16 17 20 4 1 5
4 10 18 11 20
4 17 1 15 12
4 19 5 2 1
4 18 15 4 18
3 1 14 14
6 19 16 13 16 13 7
4 15 3 10 1
7 14 19 10 16 10 5 6
5 18 16 11 18 5
5 19 18 2 3 8
4 3 3 1 11
7 14 3 13 16 2 4 4
3 20 4 5
4 15 5 6 20
3 14 6 3
6 7 2 18 4 13 3
4 2 19 19 4
7 13 20 5 1 14 3 11
7 20 16 16 12 12 2 5
7 10 5 19 17 10 18 18
6 8 9 3 18 8 9
8 10 17 5 8 12 15 13 6
3 1 11 3
6 2 3 4 17 20 15
3 13 15 16
4 4 17 1 18
7 13 2 5 14 8 4 3
7 16 7 5 20 13 12 8
4 11 20 12 13
5 5 12 10 14 12
6 2 19 19 7 6 13
2 4 2
2 6 7
3 2 16 16`

func solve(a []int64) int64 {
	n := len(a)
	pref := make([]int64, n+1)
	for i := 0; i < n; i++ {
		pref[i+1] = pref[i] + a[i]
	}

	best := pref[n]
	for j := 0; j < n; j++ {
		s := (j + 2) / 2
		cand := pref[s] + a[j]*int64(n-s)
		if cand > best {
			best = cand
		}
	}
	return best
}

func buildTestCase(fields []string) (testCase, error) {
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return testCase{}, fmt.Errorf("bad n: %w", err)
	}
	if len(fields) != 1+n {
		return testCase{}, fmt.Errorf("expected %d numbers, got %d", n, len(fields)-1)
	}

	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		val, err := strconv.ParseInt(fields[i+1], 10, 64)
		if err != nil {
			return testCase{}, fmt.Errorf("bad value at %d: %w", i, err)
		}
		arr[i] = val
	}

	input := fmt.Sprintf("%d\n%s\n", n, strings.Join(fields[1:], " "))
	expected := fmt.Sprint(solve(arr))
	return testCase{input: input, expected: expected}, nil
}

func loadTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		tc, err := buildTestCase(fields)
		if err != nil {
			return nil, fmt.Errorf("case %d: %w", idx+1, err)
		}
		cases = append(cases, tc)
	}
	return cases, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierE /path/to/binary")
		os.Exit(1)
	}

	bin := os.Args[1]
	cases, err := loadTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected: %s\n got: %s\n", idx+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
