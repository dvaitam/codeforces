package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesA.txt.
const testcasesAData = `
4 10111111 00100101
1 01 10
3 110111 000101
3 010000 010011
1 11 01
2 1101 1010
2 0011 0000
1 01 10
2 1111 1010
4 10001001 01100000
1 00 00
4 01000000 10001010
2 0111 0010
1 10 11
4 00000110 10011111
3 001010 100111
3 011100 010001
4 10110010 10110011
3 101000 011100
1 00 00
2 1100 1100
3 101010 010000
3 110111 100110
2 0111 1010
1 10 01
3 101010 010011
3 110011 110101
4 01100100 00001001
3 110101 010101
1 00 10
1 11 10
1 00 00
4 11001011 11110001
3 100101 110000
2 0010 1100
4 11000101 11100001
3 101111 110000
3 100110 100110
4 11001000 01010110
1 11 01
`

type testCase struct {
	n int
	s string
	t string
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesAData), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 3 {
			return nil, fmt.Errorf("line %d: expected 3 fields got %d", idx+1, len(fields))
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %w", idx+1, err)
		}
		cases = append(cases, testCase{n: n, s: fields[1], t: fields[2]})
	}
	return cases, nil
}

func solve(tc testCase) string {
	n := tc.n
	s := []byte(tc.s)
	t := []byte(tc.t)
	cnt11, cnt10, cnt01 := 0, 0, 0
	for i := 0; i < 2*n; i++ {
		a := s[i] - '0'
		b := t[i] - '0'
		switch a + b {
		case 2:
			cnt11++
		case 1:
			if a == 1 {
				cnt10++
			} else {
				cnt01++
			}
		}
	}
	D := 0
	turn := 0 // 0 first, 1 second
	for i := 0; i < cnt11; i++ {
		if turn == 0 {
			D++
		} else {
			D--
		}
		turn ^= 1
	}
	for i := 0; i < cnt10; i++ {
		if turn == 0 {
			D++
		}
		turn ^= 1
	}
	for i := 0; i < cnt01; i++ {
		if turn == 1 {
			D--
		}
		turn ^= 1
	}
	switch {
	case D > 0:
		return "First"
	case D < 0:
		return "Second"
	default:
		return "Draw"
	}
}

func runCandidate(bin string, tc testCase) (string, error) {
	input := fmt.Sprintf("%d\n%s\n%s\n", tc.n, tc.s, tc.t)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
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
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for i, tc := range tests {
		expect := solve(tc)
		got, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
