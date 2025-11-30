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
	a string
	b string
}

// solve embeds the logic from 1704A.go.
func solve(tc testCase) string {
	n, m := tc.n, tc.m
	a, b := tc.a, tc.b
	ok := true
	for i := 1; i < m && ok; i++ {
		if a[n-m+i] != b[i] {
			ok = false
		}
	}
	if ok {
		target := b[0]
		prefix := a[:n-m+1]
		found := false
		for i := 0; i < len(prefix); i++ {
			if prefix[i] == target {
				found = true
				break
			}
		}
		if !found {
			ok = false
		}
	}
	if ok {
		return "YES"
	}
	return "NO"
}

// Embedded copy of testcasesA.txt (n m a b per line).
const testcaseData = `
2 1 10 1
4 4 1001 0110
6 4 100100 0010
4 1 0110 1
2 1 11 0
8 1 00000000 1
4 4 1010 0000
2 1 01 0
7 3 1000001 000
6 2 100000 10
3 1 000 0
8 3 00101010 101
8 8 01000110 00010001
3 1 000 1
8 2 10000101 11
2 2 00 10
4 2 1010 00
2 2 10 10
5 5 10101 10101
3 3 101 101
3 1 001 0
4 2 1010 00
7 2 1010101 00
2 2 00 01
7 2 1100000 00
2 2 00 01
6 1 000000 1
2 1 01 1
5 2 00000 00
4 4 0000 0101
2 1 00 1
4 1 0101 0
2 1 01 1
4 3 1000 111
7 2 1000011 00
3 1 101 1
4 2 0010 10
7 7 1101111 1011101
7 3 0110110 000
3 1 000 1
7 7 0100000 1011010
4 4 1110 0111
6 4 001000 1000
4 3 1010 001
5 2 11100 00
5 2 11111 00
6 1 110111 1
7 1 1010100 1
5 3 00110 010
6 5 010001 00010
8 7 00000001 1000000
7 2 1001000 00
5 2 01101 00
8 4 11001110 1110
7 4 1011010 1110
7 3 1110101 101
7 4 1010000 0001
2 2 00 00
5 2 00100 01
2 1 01 0
7 6 0111110 111111
5 2 00100 01
2 1 01 1
8 6 11000001 100011
5 5 10111 10110
7 1 1101100 1
4 2 1100 01
5 2 10110 10
8 3 00000000 111
6 3 011110 110
7 6 1011010 000101
4 3 1000 001
3 2 010 00
6 4 111111 1101
4 4 1011 1010
2 2 10 00
6 3 111111 011
6 4 110000 1010
3 2 010 10
8 7 11111111 1111111
7 1 1111111 1
5 1 11111 1
7 1 1010101 1
2 1 01 0
3 3 111 100
7 1 1010101 1
2 2 11 01
3 2 101 01
3 1 100 0
3 2 101 01
6 3 000100 010
6 2 000000 00
5 2 00000 01
5 2 10111 00
7 5 0001000 00000
6 6 010101 010101
3 2 000 00
2 2 01 01
8 8 01000000 00100000
2 2 11 11
7 5 1100101 11000
7 5 1000000 10001
3 3 111 111
6 4 000000 1000
5 5 00110 10101
8 4 11110011 1001
8 6 00001111 000101
2 1 10 1
4 4 1001 1001
2 2 01 10
`

func loadTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	tests := make([]testCase, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 4 {
			return nil, fmt.Errorf("line %d: expected 4 fields got %d", i+1, len(fields))
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n: %v", i+1, err)
		}
		m, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad m: %v", i+1, err)
		}
		tests = append(tests, testCase{n: n, m: m, a: fields[2], b: fields[3]})
	}
	return tests, nil
}

func runCase(bin string, tc testCase, expected string) error {
	var input strings.Builder
	input.WriteString("1\n")
	input.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	input.WriteString(tc.a)
	input.WriteByte('\n')
	input.WriteString(tc.b)
	input.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := loadTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load testcases:", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		exp := solve(tc)
		if err := runCase(bin, tc, exp); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
