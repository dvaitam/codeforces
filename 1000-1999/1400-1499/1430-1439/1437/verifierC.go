package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const infinity = 1 << 30

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// solveCase mirrors one test case from 1437C.go.
func solveCase(t []int) int {
	n := len(t)
	sort.Ints(t)
	maxTime := 2*n + 5
	if n == 0 {
		maxTime = 5
	}
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, maxTime)
		for j := range dp[i] {
			dp[i][j] = infinity
		}
	}
	for j := 0; j < maxTime; j++ {
		dp[0][j] = 0
	}
	for i := 1; i <= n; i++ {
		for j := 1; j < maxTime; j++ {
			val1 := dp[i][j-1]
			val2 := dp[i-1][j-1] + abs(t[i-1]-j)
			if val1 < val2 {
				dp[i][j] = val1
			} else {
				dp[i][j] = val2
			}
		}
	}
	return dp[n][maxTime-1]
}

// referenceSolve reproduces 1437C input/output handling.
func referenceSolve(input string) (string, error) {
	fields := strings.Fields(input)
	if len(fields) == 0 {
		return "", fmt.Errorf("empty input")
	}
	pos := 0
	nextInt := func() (int, error) {
		if pos >= len(fields) {
			return 0, fmt.Errorf("unexpected EOF")
		}
		val, err := strconv.Atoi(fields[pos])
		pos++
		return val, err
	}
	q, err := nextInt()
	if err != nil {
		return "", fmt.Errorf("bad q: %v", err)
	}
	answers := make([]string, 0, q)
	for i := 0; i < q; i++ {
		n, err := nextInt()
		if err != nil {
			return "", fmt.Errorf("case %d: bad n: %v", i+1, err)
		}
		if pos+n > len(fields) {
			return "", fmt.Errorf("case %d: insufficient numbers", i+1)
		}
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j], err = strconv.Atoi(fields[pos+j])
			if err != nil {
				return "", fmt.Errorf("case %d: bad value at %d: %v", i+1, j+1, err)
			}
		}
		pos += n
		answers = append(answers, strconv.Itoa(solveCase(arr)))
	}
	if pos != len(fields) {
		return "", fmt.Errorf("extra tokens at end")
	}
	return strings.Join(answers, "\n"), nil
}

type testCase struct {
	input string
}

// Embedded testcases from testcasesC.txt.
const testcaseData = `
2 4 10 18 7 13 1 14
1 1 3
3 1 10 3 2 1 2 5 2 5 20 6 16
1 5 5 1 18 5 4
1 6 3 20 4 7 4 16
4 1 4 4 15 17 6 18 1 5 7 12 14 17 5 3 17 2
5 2 4 13 3 17 11 18 6 19 18 13 10 8 19 6 3 5 20 13 12 10 3 17 12 1
1 1 5
3 6 7 5 6 4 16 16 6 5 9 9 5 2 6 6 12 11 7 4 17 16
2 6 9 11 7 4 20 8 3 13 4 9
4 5 12 2 11 8 13 1 14 5 9 6 9 2 9 6 3 11 14 1 9 8
2 3 12 17 4 2 16 7
5 5 19 3 14 10 2 3 1 1 20 2 11 6 4 8 19 18 7 1 18
2 4 18 15 19 20 3 12 12 4
1 4 12 6 16 8
4 7 20 1 7 20 16 14 13 2 5 14 4 15 7 9 3 6 6 19 1 10 6 4
2 3 6 13 9 5 14 11 18 14 16
4 1 9 5 11 9 9 10 13 5 17 1 16 17 17 2 14 9
4 6 2 1 14 18 3 9 6 1 1 16 13 18 12 5 16 17 15 18 16 5 6 19 1 20 12
2 3 14 3 8 1 6
5 4 5 16 11 15 7 5 16 9 16 11 17 8 7 2 12 17 9 13 18 20 1 10 2 17 5
2 1 11 2 2 1
4 4 3 19 17 2 2 15 19 2 19 9 1 10
3 2 15 12 1 11 1 4
5 6 4 10 10 19 12 6 7 1 16 13 13 7 3 14 7 9 16 14 16 13 5 14 7 7 14 7 9 3 5 15 3 12 10 15
3 2 6 15 2 11 8 4 5 4 18 19
5 4 11 6 3 11 2 9 12 6 4 4 13 11 11 10 3 1 7 10 6 6 6 18 6 16 8
1 6 17 9 4 2 13 3
4 4 16 5 5 2 1 5 6 4 13 20 4 10 3 5 19 4 7 10 13
2 5 15 9 11 6 13 3 6 6 3
2 6 17 12 15 3 11 17 4 10 15 2 5
3 2 16 8 2 3 5 7 15 14 4 12 8 4 11
2 2 8 20 5 11 16 11 19 16
5 6 9 10 1 20 20 6 3 7 17 18 6 4 10 15 6 4 2 4 17 17 18 20 5 18 13 14 11 20
3 2 11 16 2 17 1 6 16 5 1 3 18 1
2 7 4 2 11 10 7 3 2 7 8 6 12 13 17 20 12
4 2 13 2 5 8 19 1 15 1 4 5 7 3 19 6 5 11 9 10 15 11
5 4 16 11 9 17 7 10 15 18 11 13 11 13 7 9 1 7 15 10 13 14 6 5 6 16 7 6 20 2 3 9
5 7 17 6 19 6 16 4 16 2 14 6 1 2 5 2 3 4 19 5 2 19 7
3 3 4 4 13 4 5 11 19 17 2 20 8
5 5 15 20 10 4 5 4 4 7 5 6 6 9 9 11 6 15 11 5 4 3 13 1 2 2 11 14
2 7 11 14 19 4 11 3 16 5 16 10 8 13 2
3 3 6 3 18 3 15 1 20 6 6 5 18 19 11 15
1 7 13 9 15 10 1 13 1
4 7 4 5 10 11 10 20 5 4 15 1 20 11 2 14 1 3 6 19 6
1 4 17 15 18 5
4 7 11 5 18 12 19 6 3 1 8 1 7 4 7 1 13 11
3 3 19 9 15 5 2 10 15 14 8 7 16 1 17 20 13 11 7
5 3 10 6 6 6 9 7 4 12 13 5 3 9 20 6 7 13 1 3 6 13 14 3 1 19
2 5 3 7 11 15 18 6 20 4 6 6 20 19
5 1 15 7 9 7 18 2 16 3 1 6 2 20 7 6 16 7 5 10 19 19 3 5 3 18 1 20
3 6 10 18 8 19 10 1 5 14 12 3 11 9 7 12 4 10 19 5 12 2
1 7 20 13 1 14 8 14 19
4 3 1 6 19 4 18 3 1 2 7 3 5 13 13 20 6 8 6 18 4 1 2 13 13
3 2 13 1 6 11 14 17 7 17 12 5 8 2 14 16 17
1 7 7 11 15 6 19 7 7
1 2 4 10
3 3 10 5 6 5 14 9 2 13 4 7 19 16 6 13 8 14 13
2 5 15 20 12 2 20 5 8 7 17 13 16
2 5 3 18 4 5 6 3 1 11 4
2 5 17 20 6 18 14 3 1 12 5
1 5 7 15 20 9 14
1 7 8 19 9 7 11 4 9
4 3 10 2 9 6 4 12 10 7 11 8 2 7 7 6 13 4 18 9 6 12
4 1 18 6 2 18 14 4 9 7 5 10 3 4 17 14 6 18 12 6 5 4 10
3 2 16 10 2 18 19 6 8 3 16 7 12 14
3 3 4 5 15 7 12 17 5 3 15 13 16 3 20 10 16
2 3 2 14 3 4 19 6 3 4
4 2 10 4 4 19 10 18 1 5 2 13 14 17 12 2 8 6
1 1 18
5 2 10 12 4 17 8 14 5 5 20 6 11 6 1 3 5 4 17 3 20 8 4
5 1 16 5 20 16 10 11 20 1 7 1 13 3 17 11 8
5 3 9 3 15 7 5 2 4 15 9 11 9 4 1 11 11 8 7 12 12 12 13 7 1 13 6 19 18 10 13 20 1
5 7 4 11 1 14 5 3 5 1 18 3 7 12 8 5 3 8 14 11 10 1 3
4 3 13 7 6 6 15 1 6 1 6 19 6 2 3 12 2 8 8 4 9 10 19 12
5 3 16 20 17 7 2 8 5 7 8 4 5 1 2 5 18 3 19 20 14 4 15 13 1 20
5 6 11 17 7 5 20 10 5 10 13 20 17 2 4 19 19 8 13 3 18 10 4 7 11 19 9 17 9 15 15
1 5 10 11 3 14 11
1 7 4 11 19 4 18 7 19
1 1 1
3 6 15 4 17 1 16 8 3 17 10 9 2 13 20
4 4 10 11 1 19 7 20 20 10 4 11 3 5 4 15 9 12 9 5 8 6 18 15 15
3 4 4 3 6 18 1 9 3 11 12 2
4 6 11 6 10 16 2 18 5 10 11 9 20 11 1 12 7 1 15 20 6 12 20 10
3 4 1 15 6 8 3 8 4 2 4 15 13 20 2
3 1 1 4 9 7 6 15 1 13
5 1 9 1 5 1 1 3 14 17 19 3 3 11 16
5 3 9 11 8 1 17 2 15 17 6 1 13 11 2 17 7 6 4 18 8 1 18 11
2 7 8 19 16 18 16 14 6 4 9 6 2 12
3 7 11 9 8 13 1 2 4 2 20 11 5 6 1 2 18 6
5 2 3 5 2 18 6 4 13 1 12 10 3 19 1 8 5 1 12 12 11 3
5 6 12 12 3 14 7 12 5 9 5 7 17 6 5 5 17 4 11 19 5 18 20 4 20 6 3 9 15 9
4 2 13 14 3 4 11 5 3 6 11 15 7 3 12 12 9 12 5 14
2 2 6 19 4 20 14 19 8
1 3 5 1 15
4 5 4 2 13 10 20 6 7 14 1 20 8 13 5 16 16 20 14 13 4 2 10 6 6
3 1 17 5 15 6 19 13 8 6 5 9 1 15 11 16
3 5 8 10 19 1 16 6 5 9 7 12 1 10 5 14 8 1 8 20
4 2 9 19 7 18 7 5 13 19 16 12 7 16 20 5 1 20 13 14 3 10 5 10
4 3 13 20 8 7 15 7 6 16 15 2 14 2 17 10 6 14 4 10 9 6 5
`

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	res := make([]testCase, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		t, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad q: %v", i+1, err)
		}
		pos := 1
		for caseIdx := 0; caseIdx < t; caseIdx++ {
			if pos >= len(fields) {
				return nil, fmt.Errorf("line %d: case %d missing n", i+1, caseIdx+1)
			}
			n, err := strconv.Atoi(fields[pos])
			if err != nil {
				return nil, fmt.Errorf("line %d: case %d bad n: %v", i+1, caseIdx+1, err)
			}
			pos++
			if pos+n > len(fields) {
				return nil, fmt.Errorf("line %d: case %d missing values", i+1, caseIdx+1)
			}
			for j := 0; j < n; j++ {
				if _, err := strconv.Atoi(fields[pos+j]); err != nil {
					return nil, fmt.Errorf("line %d: case %d bad value %d: %v", i+1, caseIdx+1, j+1, err)
				}
			}
			pos += n
		}
		if pos != len(fields) {
			return nil, fmt.Errorf("line %d: extra tokens", i+1)
		}
		res = append(res, testCase{input: line + "\n"})
	}
	return res, nil
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		expected, err := referenceSolve(tc.input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := runCandidate(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Printf("test %d failed\ninput: %sexpected: %s\ngot: %s\n", idx+1, tc.input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
