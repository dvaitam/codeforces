package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcaseData = `
7 1011111
6 001001
9 010011011
6 011100
2 10
10 1101000001
2 01
9 101101011
2 11
10 0100001100
3 000
9 110011111
2 10
8 10001001
3 110
2 00
1 0
1 0
4 0101
2 00
4 0010
1 0
9 101000111
3 010
10 0101110000
3 110
10 1001111110
9 010101001
6 110111
2 00
6 000111
1 1
10 1001010110
1 1
6 110100
3 011
6 000000
2 01
5 00110
1 1
4 1010
2 00
3 101
1 1
5 10100
5 01111
5 01011
9 001101101
1 1
3 111
8 00001111
5 10101
6 111111
3 010
6 101001
9 100000110
6 110100
6 000111
3 100
10 1111001001
8 00110010
4 1010
6 000000
8 01001111
9 100101101
4 1011
6 101100
6 110010
6 100000
6 001100
9 010101101
10 1011110110
8 11111101
5 10100
10 0100011100
9 110110100
9 100000010
10 1010111000
5 11110
10 0001010010
5 00100
3 101
9 101100000
6 001111
7 1110111
7 1010001
3 011
1 1
8 01010100
1 0
8 01110101
8 01100100
7 1000010
3 001
5 01010
6 011001
6 000001
6 001110
10 1110001100
1 0
1 1
5 10111
9 101000101
6 010101
5 01110
4 1111
6 111011
6 110011
5 00111
10 1101101011
`

type testCase struct {
	input    string
	expected string
}

func solve(n int, s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] == '0' {
			return "YES"
		}
	}
	return "NO"
}

func loadCases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	cases := make([]testCase, 0, len(lines))
	for caseIdx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("case %d: expected 2 fields, got %d", caseIdx+1, len(fields))
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("case %d: bad n: %w", caseIdx+1, err)
		}
		s := fields[1]
		input := fmt.Sprintf("1\n%d %s\n", n, s)
		cases = append(cases, testCase{
			input:    input,
			expected: solve(n, s),
		})
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
		fmt.Println("usage: verifierA /path/to/binary")
		os.Exit(1)
	}
	cases, err := loadCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load cases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		got, err := run(os.Args[1], tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
