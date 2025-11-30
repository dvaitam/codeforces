package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// solveCase mirrors 1442A.go logic for one array.
func solveCase(arr []int) string {
	n := len(arr)
	if n == 0 {
		return "YES"
	}
	L := make([]int, n)
	R := make([]int, n)
	for i := 0; i < n; i++ {
		if i == 0 || arr[i] < L[i-1] {
			L[i] = arr[i]
		} else {
			L[i] = L[i-1]
		}
	}
	for i := n - 1; i >= 0; i-- {
		if i == n-1 || arr[i] < R[i+1] {
			R[i] = arr[i]
		} else {
			R[i] = R[i+1]
		}
	}
	for i := 0; i < n; i++ {
		if L[i]+R[i] < arr[i] {
			return "NO"
		}
	}
	return "YES"
}

// referenceSolve reproduces 1442A I/O on provided input.
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
	t, err := nextInt()
	if err != nil {
		return "", fmt.Errorf("bad t: %v", err)
	}
	answers := make([]string, 0, t)
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		n, err := nextInt()
		if err != nil {
			return "", fmt.Errorf("case %d bad n: %v", caseIdx+1, err)
		}
		if pos+n > len(fields) {
			return "", fmt.Errorf("case %d missing values", caseIdx+1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i], err = strconv.Atoi(fields[pos+i])
			if err != nil {
				return "", fmt.Errorf("case %d bad value %d: %v", caseIdx+1, i+1, err)
			}
		}
		pos += n
		answers = append(answers, solveCase(arr))
	}
	if pos != len(fields) {
		return "", fmt.Errorf("extra tokens at end")
	}
	return strings.Join(answers, "\n"), nil
}

type testCase struct {
	input string
}

// Embedded testcases from testcasesA.txt.
const testcaseData = `
100
4
6 0 4 8
4
6 4 7 5
5
3 8 2 4 2
1
9
3
8 9 2
3
1 1 10
3
7 8 1
3
6 5 9
6
3 8 7 7 8 4
1
8
1
1
6
6 10 10 0 9 7
3
3 5 1
2
9 3
2
2 8
4
1 1 5 8
4
1 4 8 4
6
1 8 5 8 3 9
5
9 4 7 1 9
4
5 9 3 4
2
3 2
1
9
6
4 7 1 1 10 2
2
0 1
6
8 10 6 8 4 8
2
3 10
5
6 9 4 7 7
6
10 5 1 5 9 1
4
9 10 5 3
2
0 4
1
3
3
2 5 6
1
1
2
3 0
5
10 8 9 10 1
1
1
6
3 9 9 1 6 1
3
1 0 9
1
3
2
1 7
2
0 10
1
8
4
9 1 4 1
2
1 10
3
5 6 2
1
8
4
0 9 1 6
2
4 5
6
7 9 2 10 3 0
6
2 2 5 8 4 1
5
7 10 2 0 7
6
6 9 8 4 10 5
4
10 4 2 8
6
0 7 1 5 0 8
3
2 3 7
3
9 4 10
3
9 10 9
2
4 6
6
6 10 1 0 9 3
6
5 2 3 3 10 7
4
10 9 6 0
4
9 6 10 0
2
7 1
3
2 7 8
4
8 9 0 0
4
5 4 7 0
4
3 8 10 1
6
2 0 6 10 6 5
1
3
1
0
6
8 9 1 3 1 9
6
3 4 4 2 1 7
4
10 1 0 4
4
1 4 2 10
5
10 10 5 1 2
3
0 0 0
2
10 4
5
5 5 9 0 9
6
7 10 7 10 6 5
5
2 3 6 9 4
1
2
2
4 5
3
5 1 5
5
0 0 4 2 2
5
4 5 6 8 2
3
1 7 3
1
4
2
8 1
3
6 5 4
4
1 1 8 7
4
5 5 1 7
1
7
4
0 4 5 10
2
2 10
5
6 10 1 1 1
2
3 0
4
0 1 6 8
5
4 7 7 9 10
2
6 1
3
3 4 9
`

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("bad t: %v", err)
	}
	res := make([]testCase, 0, t)
	pos := 1
	for i := 0; i < t; i++ {
		if pos >= len(fields) {
			return nil, fmt.Errorf("case %d missing n", i+1)
		}
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d bad n: %v", i+1, err)
		}
		pos++
		if pos+n > len(fields) {
			return nil, fmt.Errorf("case %d missing values", i+1)
		}
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j], err = strconv.Atoi(fields[pos+j])
			if err != nil {
				return nil, fmt.Errorf("case %d bad value %d: %v", i+1, j+1, err)
			}
		}
		pos += n
		var input bytes.Buffer
		fmt.Fprintln(&input, 1)
		fmt.Fprintln(&input, n)
		for j, v := range arr {
			if j > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
		res = append(res, testCase{input: input.String()})
	}
	if pos != len(fields) {
		return nil, fmt.Errorf("testcase count mismatch: read %d tokens, have %d", pos, len(fields))
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, tc.input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
