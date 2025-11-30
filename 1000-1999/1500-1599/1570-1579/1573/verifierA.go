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
	s string
}

// Embedded testcases from testcasesA.txt.
const testcaseData = `
1 0
1 1
1 2
1 3
1 4
1 5
1 6
1 7
1 8
1 9
2 00
2 01
2 02
2 03
2 04
2 05
2 06
2 07
2 08
2 09
2 10
2 11
2 12
2 13
2 14
2 15
2 16
2 17
2 18
2 19
2 20
2 21
2 22
2 23
2 24
2 25
2 26
2 27
2 28
2 29
2 30
2 31
2 32
2 33
2 34
2 35
2 36
2 37
2 38
2 39
2 40
2 41
2 42
2 43
2 44
2 45
2 46
2 47
2 48
2 49
3 000
3 001
3 002
3 003
3 004
3 010
3 011
3 012
3 013
3 014
3 020
3 021
3 022
3 023
3 024
3 030
3 031
3 032
3 033
3 034
3 040
3 041
3 042
3 043
3 044
3 100
3 101
3 102
3 103
3 104
3 110
3 111
3 112
3 113
3 114
3 120
3 121
3 122
3 123
3 124
`

// solve mirrors 1573A.go.
func solve(tc testCase) int {
	ans := 0
	for i := 0; i < tc.n; i++ {
		d := int(tc.s[i] - '0')
		ans += d
		if i < tc.n-1 && d > 0 {
			ans++
		}
	}
	return ans
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			return nil, fmt.Errorf("case %d: expected 2 fields", idx+1)
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("case %d bad n: %v", idx+1, err)
		}
		s := parts[1]
		if len(s) != n {
			return nil, fmt.Errorf("case %d: len(s)=%d, expected %d", idx+1, len(s), n)
		}
		cases = append(cases, testCase{n: n, s: s})
	}
	return cases, nil
}

func runCandidate(bin string, tc testCase) (string, error) {
	input := fmt.Sprintf("1\n%d\n%s\n", tc.n, tc.s)
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		expect := strconv.Itoa(solve(tc))
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
