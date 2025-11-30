package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// solve mirrors 1499B.go.
func solve(s string) string {
	n := len(s)
	dp := make([][][][]bool, n+1)
	for i := range dp {
		dp[i] = make([][][]bool, 2)
		for p := range dp[i] {
			dp[i][p] = make([][]bool, 2)
			for r := range dp[i][p] {
				dp[i][p][r] = make([]bool, 2)
			}
		}
	}
	dp[0][0][0][0] = true
	for i := 0; i < n; i++ {
		ch := s[i]
		for phase := 0; phase < 2; phase++ {
			for prev := 0; prev < 2; prev++ {
				for rem := 0; rem < 2; rem++ {
					if !dp[i][phase][prev][rem] {
						continue
					}
					if phase == 0 {
						if ch == '0' {
							dp[i+1][0][0][rem] = true
						} else {
							dp[i+1][1][0][rem] = true
						}
					} else {
						if ch == '1' {
							dp[i+1][1][0][rem] = true
						}
					}
					if prev == 0 {
						dp[i+1][phase][1][rem|1] = true
					}
				}
			}
		}
	}
	ok := false
	for phase := 0; phase < 2; phase++ {
		for prev := 0; prev < 2; prev++ {
			if dp[n][phase][prev][1] {
				ok = true
			}
		}
	}
	if ok {
		return "YES"
	}
	return "NO"
}

type testCase struct {
	s string
}

// Embedded testcases from testcasesB.txt.
const testcaseData = `
100
010111
00101101100100
01
10011010
110100101
101111010110110100
111010110000001111
0100101101
11101100000100001
10
0011000101
11110011110
0101000100
1101001
11010001001101100001
0101001
01011101110001011010
11000000001111
101010010101
010100
01110110010001000010
10000
111110010001111
100110101110
0100101010101101011
100000010
0000111100
100000110100001001
000101111101101
01001000111100010
1100111
111000100
11101001001
001111000101000
0111011010011011
100100010
1001111
1110010
00001001111
100111111101100101
110101011110101010
10001111011101
010101010000111000
100110011101
1011101001110100010
1000111100
10001110
0000011
100001011100
111011110
00001111
10000100
11001100000101000101
11111000
01101101010100000100
1101
001
110101
010101101101
10011
110010100111011000
11111110001
101111110000
110111111100010000
001001000001011
00101110101
001101101100100
11111110100
0101111000
11
1111110110100011
111111
110111000000001
101000
001100101101000
01111010010011
1101100110110011010
101001
1111101110
001001001111
00100
0000011100100011110
11001
010001
010100011001
0101011011
1101000000110111101
10101
0000011110011100
1000101011
1111001000
01111
001000111
0111111100
100
1100101010001101111
11101100
11
00011
100010110100100100
`

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("no test data")
	}
	t, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return nil, fmt.Errorf("bad t: %v", err)
	}
	tests := make([]testCase, 0, t)
	for i := 1; i < len(lines); i++ {
		s := strings.TrimSpace(lines[i])
		if s == "" {
			continue
		}
		tests = append(tests, testCase{s: s})
	}
	if len(tests) != t {
		return nil, fmt.Errorf("testcase count mismatch: declared %d, have %d", t, len(tests))
	}
	return tests, nil
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
		input := fmt.Sprintf("1\n%s\n", tc.s)
		expected := solve(tc.s)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("test %d failed\ninput: %sexpected: %s\ngot: %s\n", idx+1, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
