package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const MOD int64 = 998244353

type testCase struct {
	input    string
	expected string
}

const testcaseData = `100
7
1011111
1001001
3
100
110
6
110111
000101
6
010000
010011
2
11
01
4
1101
1010
4
1010
1010
6
000100
000010
3
010
110
3
000
001
6
010000
100101
5
00010
11110
6
100001
001101
7
1100011
0000010
1
0
0
1
1
1
8
10110101
00000100
7
1110010
1100010
7
1101000
1010010
8
10110100
10111110
4
1001
1001
5
01100
10001
8
10101011
00011011
2
01
01
3
000
100
5
10000
01001
8
00000111
10011111
3
111
101
5
10110
01000
4
1111
0111
4
1001
1110
4
0010
0010
1
1
0
4
1001
0111
2
01
11
5
10011
11111
5
10010
10001
4
1000
0101
3
100
010
3
110
100
4
0101
1000
8
10011000
00101100
6
110101
101101
2
10
11
4
1010
1000
8
00100000
10100001
2
01
10
3
000
011
5
01110
10000
4
0001
1001
4
1111
1011
7
1111111
1111110
6
111001
010100
3
101
101
5
10101
01000
5
00100
01000
1
0
0
4
1000
1000
8
11011111
11111101
7
1111010
1100011
3
100
100
4
1010
0011
7
0110101
1100010
6
101111
100001
7
0100001
1001001
3
110
110
7
1010100
1111011
8
10100010
10101010
8
00111001
00000111
6
111111
000000
7
1111101
1010011
3
011
001
3
010
101
7
1000101
0011110
7
1001101
0000011
5
11101
11110
3
100
101
7
1001111
1110101
7
1010011
1000011
3
101
110
2
11
10
2
10
11
1
1
1
4
0110
1101
5
00010
01010
7
0111001
0100001
2
10
10
3
000
101
1
1
0
5
00011
01000
2
01
00
7
0101000
1011011
7
0011111
0011110
6
010001
011111
4
0110
0000
5
10001
10101
2
11
11
7
0111110
0101100
8
00100111
11011100
1
0
1
4
0001
1101
4
0001
0001`

func modPow(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
	}
	return res
}

func expectedMoves(n, d int) int64 {
	if d == 0 {
		return 0
	}
	inv := make([]int64, n+1)
	inv[1] = 1
	for i := 2; i <= n; i++ {
		inv[i] = MOD - (MOD/int64(i))*inv[MOD%int64(i)]%MOD
	}
	A := make([]int64, n+1)
	B := make([]int64, n+1)
	if n >= 1 {
		A[1] = 1
	}
	for i := 1; i < n; i++ {
		denom := n - i
		x1 := (int64(n)*A[i] - int64(i)*A[i-1]) % MOD
		if x1 < 0 {
			x1 += MOD
		}
		A[i+1] = x1 * inv[denom] % MOD

		x2 := (int64(n)*B[i] - int64(i)*B[i-1] - int64(n)) % MOD
		if x2 < 0 {
			x2 += MOD
		}
		B[i+1] = x2 * inv[denom] % MOD
	}
	diff := (A[n] - A[n-1]) % MOD
	if diff < 0 {
		diff += MOD
	}
	x := (B[n-1] + 1 - B[n]) % MOD
	if x < 0 {
		x += MOD
	}
	x = x * modPow(diff, MOD-2) % MOD
	res := (A[d]*x + B[d]) % MOD
	return res
}

func solve(a, b string) int64 {
	d := 0
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			d++
		}
	}
	return expectedMoves(len(a), d)
}

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	pos := 0
	t, err := strconv.Atoi(fields[pos])
	if err != nil {
		return nil, fmt.Errorf("bad test count: %w", err)
	}
	pos++
	cases := make([]testCase, 0, t)
	for caseNum := 0; caseNum < t; caseNum++ {
		if pos >= len(fields) {
			return nil, fmt.Errorf("case %d: missing n", caseNum+1)
		}
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d: bad n: %w", caseNum+1, err)
		}
		pos++
		if pos+1 >= len(fields) {
			return nil, fmt.Errorf("case %d: missing strings", caseNum+1)
		}
		a := fields[pos]
		b := fields[pos+1]
		pos += 2

		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		sb.WriteString(a)
		sb.WriteByte('\n')
		sb.WriteString(b)
		sb.WriteByte('\n')

		cases = append(cases, testCase{
			input:    sb.String(),
			expected: strconv.FormatInt(solve(a, b), 10),
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
		fmt.Println("usage: verifierD /path/to/binary")
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
