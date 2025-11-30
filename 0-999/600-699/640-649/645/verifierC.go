package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesC.txt so the verifier is self-contained.
const testcasesRaw = `2 1 00 1
26 24 00000000000000000000000000 12
27 25 000000000000000000000000000 13
3 0 010 0
5 3 00000 2
9 4 111101010 2
8 5 00000101 3
9 5 101110111 2
4 1 1000 1
15 5 111100101110011 2
3 0 000 0
17 1 11110111110111111 1
10 0 1111111111 0
6 3 001100 1
6 2 000000 2
5 4 01111 1
6 1 000001 1
11 10 00000000000 5
6 5 011100 2
4 3 1100 1
5 5 11111 1
10 9 0000111110 3
10 9 0000011111 4
10 9 0000111101 4
10 9 1111100000 3
10 9 0000001111 5
10 9 0111100001 5
10 9 1111111110 1
10 9 0000000001 5
10 9 1000000000 5
10 9 0000000011 6
10 9 0100000000 5
10 9 1010000000 6
10 9 1110000001 4
10 9 1111111100 2
10 9 1111111000 2
10 9 1111110001 3
10 9 1111100011 4
10 9 1111000111 4
10 9 1110001111 3
10 9 1100011111 2
10 9 1000111111 2
10 9 0001111111 3
10 9 0011111111 2
10 9 0111111111 2
9 4 111101001 2
11 5 01000001011 3
9 5 001010101 1
8 5 01000010 2
9 4 011111110 2
15 6 100011000101001 2
10 9 0111100000 4
11 5 01010111111 1
9 4 100110111 2
6 5 100111 2
6 2 111110 1
9 3 010101001 1
9 5 010111111 3
11 6 11111111110 1
10 9 0101111110 3
10 6 1011111111 2
10 2 1101100100 1
7 4 1111110 1
8 7 11111111 1
7 1 1111100 1
11 2 11111000011 1
8 4 11111100 1
12 5 111101101111 2
12 5 010011011101 2
16 10 1100101000101111 2
13 3 0101110001110 1
13 7 1001110101111 2
16 6 1111110001111111 2
12 2 101111001001 1
14 8 11100001000000 2
16 6 1111111111111101 1
15 4 111011010001010 1
15 2 011001001011100 1
12 4 000100110000 2
15 5 111110100111001 2
14 5 01101010001101 1
12 5 111111011001 2
12 6 010111000001 1
16 7 1000010100001001 2
14 4 11000000000100 2
13 5 1010101010101 2
14 2 11111111111110 1
12 5 100110101010 2
14 2 10100000000001 1
16 6 0101000101000101 2
13 4 0101010101011 2
12 6 101001100100 2
14 3 10010010101000 1
15 2 010101010101010 1
14 8 10000000000001 4
13 6 1110000001111 3
12 6 110100000111 2
12 2 110000000011 1
14 7 11111111111110 3
12 4 111110000001 2
16 7 1000000000000000 4
15 6 100000000000000 4
15 7 100000000000010 3
14 6 11111111111111 2
16 6 1010111100011010 2
14 8 01010101010101 4
13 4 0011111111110 2
13 6 1010011001001 2`

type testCase struct {
	n int
	k int
	s string
}

func feasible(r int, zeros []int, pref []int, n int, k int) bool {
	for _, idx := range zeros {
		l := idx - r
		if l < 0 {
			l = 0
		}
		rr := idx + r
		if rr >= n {
			rr = n - 1
		}
		count := pref[rr+1] - pref[l]
		if count >= k+1 {
			return true
		}
	}
	return false
}

// solveCase mirrors 645C.go.
func solveCase(tc testCase) (int, error) {
	n := tc.n
	k := tc.k
	s := tc.s
	if len(s) != n {
		return 0, fmt.Errorf("string length mismatch: got %d want %d", len(s), n)
	}
	zeros := make([]int, 0)
	pref := make([]int, n+1)
	for i := 0; i < n; i++ {
		if s[i] == '0' {
			zeros = append(zeros, i)
			pref[i+1] = pref[i] + 1
		} else {
			pref[i+1] = pref[i]
		}
	}
	lo, hi := 0, n
	for lo < hi {
		mid := (lo + hi) / 2
		if feasible(mid, zeros, pref, n, k) {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return lo, nil
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) < 3 {
			return nil, fmt.Errorf("line %d malformed", idx+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %w", idx+1, err)
		}
		k, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse k: %w", idx+1, err)
		}
		s := fields[2]
		cases = append(cases, testCase{n: n, k: k, s: s})
	}
	return cases, nil
}

func buildInput(tc testCase) string {
	return fmt.Sprintf("%d %d\n%s\n", tc.n, tc.k, tc.s)
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		expectVal, err := solveCase(tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: solve error: %v\n", idx+1, err)
			os.Exit(1)
		}
		expect := strconv.Itoa(expectVal)
		input := buildInput(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\nInput:\n%s\n", idx+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
