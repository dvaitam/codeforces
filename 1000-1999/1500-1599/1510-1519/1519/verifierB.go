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
	n int
	m int
	k int
}

// Embedded testcases from testcasesB.txt.
const testcaseData = `
100
3 10 48
2 5 9
8 8 26
2 8 15
7 7 0
8 5 14
10 2 1
1 1 0
7 4 46
1 9 8
8 8 44
4 4 18
1 7 17
2 3 23
5 2 28
9 7 24
5 5 31
9 7 4
8 4 25
7 3 20
6 2 11
9 2 33
7 6 41
1 8 7
10 10 82
3 3 0
4 9 35
4 7 36
6 8 0
7 9 66
9 4 35
1 8 18
9 4 26
8 6 47
1 9 25
10 6 59
1 4 17
10 3 35
5 1 21
2 2 14
1 5 4
2 10 19
5 2 9
5 9 34
5 8 31
8 2 15
7 6 41
4 5 19
9 4 27
1 4 3
3 1 5
8 9 69
4 9 35
9 1 8
10 6 54
1 5 4
4 1 3
2 5 9
3 7 8
1 9 1
10 4 36
8 3 39
9 1 8
6 2 11
7 10 69
2 7 13
8 1 7
7 5 34
4 6 36
3 6 17
5 2 29
9 6 68
8 9 71
1 2 1
3 9 26
6 10 32
6 6 35
5 4 38
8 3 6
6 1 5
7 3 21
2 10 24
2 10 36
2 5 9
5 10 14
8 5 39
1 5 4
1 2 1
1 4 3
10 7 69
8 3 10
2 7 24
9 5 61
6 2 11
6 1 5
5 10 49
7 6 41
2 6 29
2 5 9
10 9 60
`

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no test data")
	}
	pos := 0
	nextInt := func() (int, error) {
		if pos >= len(fields) {
			return 0, fmt.Errorf("unexpected end of data")
		}
		v, err := strconv.Atoi(fields[pos])
		pos++
		return v, err
	}
	t, err := nextInt()
	if err != nil {
		return nil, fmt.Errorf("bad t: %v", err)
	}
	res := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		n, err := nextInt()
		if err != nil {
			return nil, fmt.Errorf("case %d bad n: %v", i+1, err)
		}
		m, err := nextInt()
		if err != nil {
			return nil, fmt.Errorf("case %d bad m: %v", i+1, err)
		}
		k, err := nextInt()
		if err != nil {
			return nil, fmt.Errorf("case %d bad k: %v", i+1, err)
		}
		res = append(res, testCase{n: n, m: m, k: k})
	}
	if pos != len(fields) {
		return nil, fmt.Errorf("extra tokens at end of test data")
	}
	return res, nil
}

// solve mirrors 1519B.go.
func solve(tc testCase) string {
	if tc.n*tc.m-1 == tc.k {
		return "YES"
	}
	return "NO"
}

func runCandidate(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.k))
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.ToUpper(strings.TrimSpace(out.String())), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		expect := solve(tc)
		got, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
