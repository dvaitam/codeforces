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

const testcaseData = `3 3 -8 -2 -7
8 8 5 10 2 -4 -7 5 -10 2
7 5 -10 4 -2 -3 8
2 2 -10 -10
1 1 2
4 4 -10 6 -3 4
8 4 1 -3 -3 4
5 1 3
2 1 10
5 1 0
7 5 -4 -1 -1 8 5
7 5 -9 5 -3 2 3
3 2 7 1
2 2 6 -7
3 3 2 1 5
1 1 -9
5 5 8 8 2 10 -5
3 3 -3 -10 -4
4 4 6 1 8 1
8 5 7 9 -10 2 6
3 3 7 -4 3
1 1 1
4 4 5 1 3 1
1 1 4
1 1 10
3 3 8 -5 -8
5 1 -8
2 1 4
1 1 -3
5 1 9
3 2 -1 -8
3 1 -2
3 3 -2 10 -1
8 6 5 5 -7 -10 -1 2
6 4 -4 -2 -7 -2
4 4 -10 -3 -10 2
3 1 -5
8 7 7 -3 10 6 4 -3 6
1 1 8
6 6 10 3 -9 -1 -6 -4
1 1 -8
2 2 -1 -5
7 5 -2 -6 -10 7 -9
4 4 -5 9 6 -9
7 2 1 -7
4 4 8 -4 5 -7
7 3 6 5 -10
6 5 2 -1 -10 -5 -4
6 5 -6 0 3 -4 -2
2 2 7 1
8 4 -8 -9 -8 -6
3 1 7
4 3 0 9 6
5 3 0 0 -7
5 2 9 5
3 3 7 -7 0
1 1 -8
7 7 -6 -6 0 -7 9 8 2
2 1 8
2 2 1 -1
2 2 -2 -7
1 1 -10
1 1 3
2 1 -4
4 4 -5 -7 4 -5
4 2 -7 3
7 7 7 -1 7 -2 5 0 -7
4 3 -9 -10 -10
5 5 0 4 2 0 2
2 1 0
8 2 -2 -4
8 6 -2 -5 7 -4 -1 -4
4 3 -8 -2 -8
8 2 10 8
6 2 2 -1
1 1 -5
6 5 -1 -3 0 -7 7
2 1 -3
1 1 2
2 2 7 -8
2 1 10
1 1 1
8 8 -6 -7 6 0 -8 6 -5 -5
3 1 0
5 1 6
5 2 -4 -6
1 1 9
4 2 -1 3
3 1 -3
5 1 4
7 5 -2 7 4 7 4
1 1 0
3 2 5 -10
7 5 -10 -9 1 8 -6
3 1 -2
5 4 8 2 -5 9
2 1 5
1 1 6
6 5 10 4 10 -3 -3
6 4 5 -3 3 0`

func solveB(n, k int, s []int64) string {
	if k == 1 {
		return "YES"
	}
	diff := make([]int64, k-1)
	for i := 1; i < k; i++ {
		diff[i-1] = s[i] - s[i-1]
	}
	for i := 1; i < len(diff); i++ {
		if diff[i] < diff[i-1] {
			return "NO"
		}
	}
	firstCount := int64(n - k + 1)
	if diff[0]*firstCount < s[0] {
		return "NO"
	}
	return "YES"
}

func parseTestcase(line string) (testCase, error) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return testCase{}, fmt.Errorf("invalid testcase: %s", line)
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return testCase{}, fmt.Errorf("bad n: %w", err)
	}
	k, err := strconv.Atoi(fields[1])
	if err != nil {
		return testCase{}, fmt.Errorf("bad k: %w", err)
	}
	if len(fields) != 2+k {
		return testCase{}, fmt.Errorf("expected %d suffix numbers, got %d", k, len(fields)-2)
	}
	svals := make([]int64, k)
	for i := 0; i < k; i++ {
		val, err := strconv.ParseInt(fields[2+i], 10, 64)
		if err != nil {
			return testCase{}, fmt.Errorf("bad s[%d]: %w", i, err)
		}
		svals[i] = val
	}

	var input strings.Builder
	input.WriteString("1\n")
	input.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 0; i < k; i++ {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(fmt.Sprintf("%d", svals[i]))
	}
	input.WriteByte('\n')

	expect := solveB(n, k, svals)
	return testCase{input: input.String(), expected: expect}, nil
}

func loadTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		tc, err := parseTestcase(line)
		if err != nil {
			return nil, fmt.Errorf("case %d: %w", idx+1, err)
		}
		cases = append(cases, tc)
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierB /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := loadTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(tc.input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(strings.ToUpper(out.String()))
		if got != tc.expected {
			fmt.Printf("test %d failed: expected %s got %s\n", idx+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
