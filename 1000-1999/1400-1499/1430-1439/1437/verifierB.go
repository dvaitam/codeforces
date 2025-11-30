package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// solveCase mirrors 1437B logic for a single string.
func solveCase(s string) int {
	cnt1, cnt2 := 0, 0
	in1, in2 := false, false
	for i := 0; i < len(s); i++ {
		var e1, e2 byte
		if i%2 == 0 {
			e1, e2 = '0', '1'
		} else {
			e1, e2 = '1', '0'
		}
		if s[i] != e1 {
			if !in1 {
				cnt1++
				in1 = true
			}
		} else {
			in1 = false
		}
		if s[i] != e2 {
			if !in2 {
				cnt2++
				in2 = true
			}
		} else {
			in2 = false
		}
	}
	if cnt1 < cnt2 {
		return cnt1
	}
	return cnt2
}

// referenceSolve reproduces 1437B.go behaviour against the provided input.
func referenceSolve(input string) (string, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return "", err
	}
	answers := make([]string, 0, t)
	for ; t > 0; t-- {
		var n int
		var s string
		if _, err := fmt.Fscan(reader, &n); err != nil {
			return "", err
		}
		if _, err := fmt.Fscan(reader, &s); err != nil {
			return "", err
		}
		if len(s) != n {
			return "", fmt.Errorf("string length %d != n %d", len(s), n)
		}
		answers = append(answers, strconv.Itoa(solveCase(s)))
	}
	return strings.Join(answers, "\n"), nil
}

type testCase struct {
	input string
}

// Embedded testcases from testcasesB.txt.
const testcaseData = `
4 1 1 1 0 10 0110010011 6 101110
4 3 100 3 111 7 1100100 7 0000100
2 10 1011110111 5 10110
1 1 0
3 1 1 1 0 2 10
5 7 0101001 6 111000 6 001001 5 11111 4 0100
4 8 10111100 2 00 3 101 10 0110001001
3 6 101011 9 111000001 4 0101
1 2 11
1 10 1010001101
4 5 00011 10 0010000010 4 0011 10 0100111100
5 3 111 2 10 10 1011100010 6 101101 1 1
5 9 001010000 1 0 6 001100 8 10110111 1 1
3 9 100110000 2 10 3 000
5 3 010 10 1001110010 6 101010 7 1111011 9 101011011
4 1 0 5 01101 7 1001010 3 100
5 3 111 5 11110 1 0 10 0101011001 5 00000
5 9 111110111 9 110100011 8 01001001 8 00010101 4 1000
4 2 11 5 11110 5 01010 7 1110110
5 5 10110 1 0 10 1100001010 6 101100 7 1000010
5 1 1 3 100 9 010111110 8 00001001 10 0101010010
1 2 11
5 8 11100100 3 010 1 1 1 1 7 0010101
3 10 1111001101 2 10 2 01
1 3 101
4 9 110011000 6 101001 1 0 7 1000000
5 6 011000 5 10010 3 111 10 0100010000 6 000000
4 1 1 7 1001010 3 110 7 1011010
4 7 0111010 5 01100 3 101 7 1011001
4 7 0101101 8 00101110 5 10111 4 1000
5 6 001110 2 10 3 000 4 0011 10 1101101101
5 5 01111 9 111000000 10 0101000011 6 011110 3 000
2 9 101110010 8 01011010
4 10 0111110010 5 01001 7 1011001 3 100
1 3 000
4 9 101110000 8 11001110 2 11 4 1100
3 4 1111 10 0011101111 6 110110
4 3 100 6 100111 9 101110100 8 11010010
3 4 1000 6 000010 3 000
4 4 0100 4 0111 1 1 4 1010
5 6 111001 2 01 1 0 10 0101010010 9 110110111
1 10 1010111010
5 7 1010111 3 011 6 100010 1 0 6 010110
4 8 00101010 9 101010011 9 000000100 9 111100111
4 8 11000001 8 01001001 10 1000111001 1 0
4 10 1100001011 2 11 4 1110 10 1000110111
4 1 0 7 0000100 5 11011 6 001001
1 5 10101
1 5 00010
5 4 1000 4 1011 3 110 4 1101 9 001110001
1 3 101
4 2 11 2 01 6 010110 1 0
3 2 11 9 100010111 10 0111101100
2 8 10001000 8 01100111
4 7 1001000 4 1001 9 000100000 9 101001011
2 7 1000000 4 1111
3 7 1000000 9 010101001 6 111001
2 10 1111101111 10 1111001010
5 3 111 7 1100110 1 0 10 0100001100 1 0
3 4 1111 2 10 1 0
1 10 1101010010
5 10 1000100011 6 001011 8 11101001 9 111010000 8 10111010
1 1 0
2 2 11 8 01001000
4 1 0 10 1011101000 10 0011110100 9 001000010
3 3 100 6 011111 9 001111110
2 3 110 2 00
1 7 0011111
5 6 110011 6 000101 7 0000101 3 111 8 10011001
5 3 111 4 0101 6 110110 6 110010 9 110101001
5 3 011 4 0100 9 101110111 1 0 4 0010
3 9 111110001 3 111 9 111000001
5 2 00 1 0 5 01001 5 00000 4 1110
5 5 00010 4 0110 5 01000 4 0100 8 01001110
3 4 0111 6 101100 7 0001011
5 4 0011 6 101110 9 001110101 2 01 2 10
1 4 0001
4 2 00 4 1000 9 001111011 4 0010
5 4 0001 8 01111100 8 10001001 5 01001 6 010011
4 7 0010100 10 1100110000 8 10011010 2 10
5 1 0 10 0111100001 4 0000 3 000 5 01111
1 2 01
1 6 111101
2 4 0000 8 10110110
3 2 11 4 1101 3 101
4 1 1 9 110001011 4 1111 3 110
1 4 1100
4 5 10110 2 01 8 10111110 2 10
1 1 1
4 4 1010 5 01101 1 0 1 0
3 6 100111 9 011000010 5 01001
1 4 0010
2 8 10101000 5 01011
2 8 10101011 9 111111000
4 6 011100 2 10 4 1100 8 11110000
4 3 011 4 0101 6 001111 2 00
3 7 1101110 7 1101111 8 11101011
4 5 00101 1 1 3 010 2 00
2 3 011 3 111
3 9 001000111 5 10111 2 01
`

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	cases := make([]testCase, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// validate that tokens count matches declared t
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		t, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad t: %v", i+1, err)
		}
		if len(fields) != 1+t*2 {
			return nil, fmt.Errorf("line %d: expected %d tokens, got %d", i+1, 1+t*2, len(fields))
		}
		cases = append(cases, testCase{input: line + "\n"})
	}
	return cases, nil
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
