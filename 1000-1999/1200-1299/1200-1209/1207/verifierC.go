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

const testcasesData = `
100
4 5 5 0110
10 1 4 1001110001
1 1 2 0
5 1 3 11111
3 3 1 001
4 3 4 1111
9 5 4 010101001
5 1 1 11010
7 2 1 1110001
10 3 5 1001000001
5 5 3 00111
3 4 4 110
10 5 3 1011111011
1 4 5 0
1 3 4 1
6 5 3 100011
8 3 5 10101111
7 1 1 0101010
7 1 1 1101001
4 5 4 1000
9 2 3 011010111
8 3 4 11101000
8 5 5 10011100
10 3 1 0000100100
9 3 2 010010110
6 2 2 000010
1 4 5 1
1 3 2 1
10 5 5 1011000000
2 4 1 01
6 2 3 011111
5 2 3 10001
2 5 2 01
8 5 5 10101111
7 4 1 0010101
3 1 3 110
8 1 5 10100100
7 1 5 0010010
1 1 4 1
10 3 1 0000010000
4 4 5 1111
6 5 4 010101
10 5 5 1010001110
3 3 2 010
9 3 4 101011100
10 3 2 1101101010
7 1 4 0000100
5 2 2 01001
3 4 1 000
5 3 1 11100
6 3 4 110011
3 5 3 010
7 1 4 1011111
2 3 5 01
9 3 1 100011000
6 5 4 101110
2 1 2 11
2 3 3 10
3 1 4 101
6 4 4 011001
10 1 1 0001100100
10 2 2 1110111111
9 3 4 000001010
4 5 2 0101
3 2 4 000
6 5 2 010111
6 2 1 001000
7 2 4 1001000
7 4 4 0101101
8 4 4 01011111
10 4 2 0011010000
3 5 4 000
2 5 2 11
3 1 1 111
7 1 1 0101100
4 1 4 1001
10 1 3 0011111010
6 3 3 001000
1 1 1 0
9 1 4 001011110
7 3 5 1010101
9 3 5 101101110
8 4 3 10110000
8 1 1 10001111
1 1 2 1
2 3 5 10
8 1 2 00000001
4 2 5 1101
3 5 1 000
2 2 3 00
8 4 2 01111011
7 2 1 0110110
8 1 5 11000101
4 4 4 0101
3 2 2 111
4 5 4 1110
5 1 2 00010
8 1 2 01110110
5 2 1 11110
4 2 4 1110
`

type testCase struct {
	n int
	a int64
	b int64
	s string
}

func parseTestcases() ([]testCase, string, error) {
	data := strings.TrimSpace(testcasesData) + "\n"
	scanner := bufio.NewScanner(strings.NewReader(data))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		return nil, "", fmt.Errorf("empty testcases")
	}
	t, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return nil, "", err
	}
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if !scanner.Scan() {
			return nil, "", fmt.Errorf("missing n at case %d", i+1)
		}
		n, _ := strconv.Atoi(scanner.Text())
		scanner.Scan()
		a, _ := strconv.ParseInt(scanner.Text(), 10, 64)
		scanner.Scan()
		b, _ := strconv.ParseInt(scanner.Text(), 10, 64)
		scanner.Scan()
		s := scanner.Text()
		cases = append(cases, testCase{n: n, a: a, b: b, s: s})
	}
	return cases, data, nil
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// solveCase mirrors 1207C.go to compute minimal cost.
func solveCase(tc testCase) int64 {
	const inf int64 = 1 << 60
	dp0 := tc.b
	dp1 := inf
	for i := 0; i < tc.n; i++ {
		if tc.s[i] == '1' {
			dp0 = inf
			dp1 = dp1 + tc.a + 2*tc.b
		} else {
			ndp0 := min(dp0+tc.a+tc.b, dp1+2*tc.a+tc.b)
			ndp1 := min(dp1+tc.a+2*tc.b, dp0+2*tc.a+2*tc.b)
			dp0, dp1 = ndp0, ndp1
		}
	}
	return dp0
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
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	testcases, input, err := parseTestcases()
	if err != nil {
		fmt.Println("failed to parse testcases:", err)
		return
	}

	got, err := run(bin, input)
	if err != nil {
		fmt.Printf("candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	outTokens := strings.Fields(got)
	if len(outTokens) != len(testcases) {
		fmt.Printf("expected %d outputs, got %d\n", len(testcases), len(outTokens))
		os.Exit(1)
	}
	for i, tc := range testcases {
		want := strconv.FormatInt(solveCase(tc), 10)
		if outTokens[i] != want {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, want, outTokens[i])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
