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

const testcaseData = `100
1 0
2 10
5 10001
7 1110011
6 110000
1 0
6 001011
9 111101101
5 11111
10 1101011111
5 01101
1 0
2 00
5 00010
4 0100
6 100000
2 00
1 1
5 00001
4 1000
2 11
4 0010
2 01
4 0000
3 110
7 1111111
2 10
1 0
1 0
4 1101
2 11
6 000000
10 1001101011
3 010
5 00001
7 1000001
1 1
6 110000
9 111101101
10 1111111111
4 0100
6 110000
1 1
7 1011001
8 10000100
2 01
3 100
3 111
2 10
7 0000000
8 00011010
6 001110
10 0010110010
10 1100111101
1 0
4 0110
10 1110111011
10 1010101111
3 111
6 010110
5 10001
7 1011001
6 111000
6 100010
2 10
5 00000
10 0101110110
1 0
2 00
1 1
8 01110111
4 1111
5 00001
1 0
10 0010011111
1 1
9 001001001
2 00
5 11111
3 000
4 1010
4 1010
1 0
7 0111100
2 11
8 00001011
3 000
5 11111
5 11001
10 0010101011
6 011001
3 000
9 111000000
6 001110
6 100010
3 001
6 111111
1 1
1 0
7 1000001
9 001110001
9 111100011`

func solve(n int, s string) int {
	l, r := 0, n-1
	for l < r && s[l] != s[r] {
		l++
		r--
	}
	if r < l {
		return 0
	}
	return r - l + 1
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	cases := make([]testCase, 0, len(lines)-1)
	for idx, line := range lines[1:] {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) != 2 {
			return nil, fmt.Errorf("case %d: expected n and string", idx+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("case %d: bad n: %w", idx+1, err)
		}
		s := fields[1]
		if len(s) != n {
			return nil, fmt.Errorf("case %d: length mismatch", idx+1)
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		sb.WriteString(s)
		sb.WriteByte('\n')

		cases = append(cases, testCase{
			input:    sb.String(),
			expected: strconv.Itoa(solve(n, s)),
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
		fmt.Println("usage: verifierC /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
