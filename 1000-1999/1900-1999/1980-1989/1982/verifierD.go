package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesD = `100
1 2 1
3 0
10
2 2 2
5 5
5 3
01
10
3 3 2
2 0 3
3 4 1
0 0 4
110
000
101
1 2 1
1 5
11
3 1 1
3
2
4
0
0
0
1 2 1
5 4
01
1 2 1
0 4
01
1 2 1
5 5
11
3 2 1
3 4
1 0
5 4
00
10
01
3 1 1
4
1
5
1
0
1
3 1 1
4
2
0
0
0
1
3 2 2
2 3
4 2
5 0
01
00
00
2 3 1
0 3 5
0 5 2
100
101
1 2 1
4 4
11
3 2 2
1 0
3 2
1 3
01
00
01
1 1 1
4
0
2 2 2
3 0
0 4
11
00
1 2 1
0 3
00
3 3 2
0 0 3
2 1 4
5 1 1
000
101
100
1 3 1
2 1 1
100
2 2 2
4 1
3 4
00
01
3 1 1
0
0
2
0
0
0
1 3 1
1 5 3
011
3 3 2
5 1 1
2 4 1
4 3 0
011
011
011
1 2 1
5 3
10
2 2 2
4 0
4 4
10
10
3 2 2
5 0
4 0
1 5
10
11
10
3 2 2
1 1
2 3
3 1
10
11
01
2 3 1
1 2 1
5 4 2
110
110
3 2 2
3 5
1 2
1 1
01
00
11
2 1 1
3
4
0
0
1 1 1
1
0
3 1 1
2
4
5
1
1
0
2 2 1
2 5
4 2
01
10
2 2 2
3 4
2 3
10
00
1 3 1
5 5 5
001
3 3 3
1 3 0
0 0 0
2 5 3
101
111
000
3 1 1
5
1
4
1
0
1
3 3 2
2 0 4
1 3 1
3 4 0
101
101
010
2 3 2
4 5 3
3 2 1
010
101
1 1 1
3
1
3 1 1
3
1
0
1
0
1
1 1 1
5
0
3 3 1
0 0 3
4 5 4
3 5 5
011
100
001
2 2 1
0 3
2 4
10
01
3 2 1
5 2
3 1
1 2
01
11
10
1 2 1
4 3
01
3 2 2
2 2
3 5
2 0
11
01
01
1 3 1
2 3 1
011
2 3 2
4 4 1
5 3 4
110
100
2 2 2
1 3
0 4
01
01
1 3 1
1 0 1
100
1 1 1
4
0
1 1 1
0
1
1 1 1
0
0
1 1 1
3
1
1 2 1
3 0
10
3 2 2
2 5
3 5
1 4
01
00
10
3 1 1
3
1
4
1
0
1
3 3 1
4 0 5
3 0 4
3 3 1
110
111
010
2 1 1
3
3
0
1
3 1 1
3
1
3
0
1
1
2 1 1
3
2
1
1
3 1 1
1
2
4
0
0
0
2 1 1
4
1
1
1
3 1 1
3
1
1
1
0
1
3 3 1
5 2 1
4 4 5
4 4 3
010
100
110
1 1 1
1
0
2 3 1
5 5 0
4 3 5
111
001
3 2 2
3 3
5 2
3 0
00
00
11
3 3 1
2 5 0
0 3 1
0 3 4
100
011
100
3 3 2
0 1 3
4 2 2
0 5 0
001
011
011
1 2 1
5 4
10
3 2 1
0 1
2 5
4 4
10
11
11
3 1 1
3
4
5
0
0
0
1 3 1
1 0 5
111
1 2 1
2 4
00
2 3 1
1 3 0
5 2 2
010
101
3 1 1
4
1
3
0
0
1
1 3 1
1 5 5
010
3 3 3
3 2 1
2 5 2
4 1 1
110
100
111
1 1 1
2
0
1 3 1
1 4 4
110
2 3 2
0 0 4
3 2 1
101
100
1 3 1
0 5 2
110
3 1 1
2
4
0
1
1
1
3 1 1
2
4
3
1
0
1
2 1 1
1
5
1
1
1 1 1
2
0
1 1 1
0
1
3 2 1
4 1
3 2
2 3
11
00
10
1 1 1
2
0
1 3 1
4 0 1
010
1 3 1
4 5 1
010
3 2 1
0 3
1 3
4 3
01
10
00
3 1 1
0
5
1
1
1
1
2 2 2
2 5
0 1
00
00
3 1 1
1
3
5
0
1
0
3 1 1
5
2
2
1
1
0
1 2 1
4 3
11
`

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func solveCase(n, m, k int, heights [][]int64, types []string) string {
	ones := make([][]int, n)
	var sum0, sum1 int64
	for i := 0; i < n; i++ {
		ones[i] = make([]int, m)
		for j := 0; j < m; j++ {
			if types[i][j] == '1' {
				ones[i][j] = 1
				sum1 += heights[i][j]
			} else {
				sum0 += heights[i][j]
			}
		}
	}
	pref := make([][]int, n+1)
	for i := range pref {
		pref[i] = make([]int, m+1)
	}
	for i := 0; i < n; i++ {
		row := 0
		for j := 0; j < m; j++ {
			row += ones[i][j]
			pref[i+1][j+1] = pref[i][j+1] + row
		}
	}
	var g int64
	for i := 0; i+k <= n; i++ {
		for j := 0; j+k <= m; j++ {
			onesCount := pref[i+k][j+k] - pref[i][j+k] - pref[i+k][j] + pref[i][j]
			diff := int64(k*k - 2*onesCount)
			if diff < 0 {
				diff = -diff
			}
			g = gcd(g, diff)
		}
	}
	diffTotal := sum1 - sum0
	if g == 0 {
		if diffTotal == 0 {
			return "YES"
		}
		return "NO"
	}
	if diffTotal%g == 0 {
		return "YES"
	}
	return "NO"
}

type testCase struct {
	n, m, k int
	heights [][]int64
	types   []string
}

func parseTests() ([]testCase, error) {
	fields := strings.Fields(testcasesD)
	pos := 0
	nextInt := func() (int64, error) {
		if pos >= len(fields) {
			return 0, fmt.Errorf("unexpected EOF")
		}
		val, err := strconv.ParseInt(fields[pos], 10, 64)
		pos++
		return val, err
	}
	t64, err := nextInt()
	if err != nil {
		return nil, err
	}
	t := int(t64)
	tests := make([]testCase, t)
	for i := 0; i < t; i++ {
		n64, err := nextInt()
		if err != nil {
			return nil, err
		}
		m64, err := nextInt()
		if err != nil {
			return nil, err
		}
		k64, err := nextInt()
		if err != nil {
			return nil, err
		}
		n := int(n64)
		m := int(m64)
		k := int(k64)
		heights := make([][]int64, n)
		for r := 0; r < n; r++ {
			heights[r] = make([]int64, m)
			for c := 0; c < m; c++ {
				val, err := nextInt()
				if err != nil {
					return nil, err
				}
				heights[r][c] = val
			}
		}
		types := make([]string, n)
		for r := 0; r < n; r++ {
			if pos >= len(fields) {
				return nil, fmt.Errorf("unexpected EOF")
			}
			types[r] = fields[pos]
			pos++
		}
		tests[i] = testCase{n: n, m: m, k: k, heights: heights, types: types}
	}
	return tests, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierD /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTests()
	if err != nil {
		fmt.Println("parse error:", err)
		os.Exit(1)
	}

	input := testcasesD + "\n"
	output, err := runCandidate(bin, input)
	if err != nil {
		fmt.Println("runtime error:", err)
		os.Exit(1)
	}
	outFields := strings.Fields(output)
	if len(outFields) != len(tests) {
		fmt.Printf("expected %d outputs, got %d\n", len(tests), len(outFields))
		os.Exit(1)
	}
	for i, tc := range tests {
		want := solveCase(tc.n, tc.m, tc.k, tc.heights, tc.types)
		if outFields[i] != want {
			fmt.Printf("case %d failed\nexpected: %s\ngot: %s\n", i+1, want, outFields[i])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
