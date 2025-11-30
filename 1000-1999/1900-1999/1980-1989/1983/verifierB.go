package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesB = `100
2 2
00
12
10
02
1 3
200
022
3 3
000
021
201
001
020
102
2 1
2
0
0
2
3 3
200
220
120
111
210
110
1 2
02
11
3 2
21
11
00
10
22
21
3 3
020
002
112
102
011
012
2 2
10
02
00
10
2 2
11
00
01
12
3 2
22
10
01
01
01
02
2 1
1
2
2
1
2 2
21
12
01
21
1 1
2
1
3 1
1
0
0
0
1
0
3 3
002
102
211
201
001
110
3 1
1
0
1
2
1
1
1 2
21
12
2 3
100
022
100
121
1 3
021
201
3 3
001
201
112
002
022
110
2 1
2
2
0
0
1 1
1
0
3 3
122
212
100
001
121
212
3 2
00
12
10
22
01
11
2 1
1
1
0
2
3 1
0
2
2
0
0
2
2 1
0
2
2
1
2 1
1
0
2
1
3 3
021
120
121
100
202
020
1 2
12
21
2 3
012
001
111
202
2 2
20
01
00
00
1 3
120
020
3 2
21
02
21
10
22
12
1 3
212
102
1 2
21
02
3 2
01
22
21
21
01
02
2 3
100
001
000
011
3 3
112
112
000
211
212
011
1 1
1
0
3 2
20
02
21
21
22
20
2 3
011
221
110
001
1 1
1
2
3 1
1
1
0
2
2
0
1 2
02
02
2 3
201
112
110
112
1 3
102
122
2 2
12
20
20
11
2 3
222
201
200
100
3 1
0
2
2
0
0
2
2 1
0
1
2
2
1 3
201
102
2 2
12
10
12
10
2 3
122
001
221
020
3 3
212
122
022
101
212
011
2 3
111
200
100
201
2 3
012
112
021
002
3 1
0
0
2
2
0
0
3 1
0
1
0
2
1
1
1 1
1
1
2 3
201
202
220
011
3 3
212
211
202
200
102
202
3 2
11
10
22
20
10
01
1 2
12
10
2 3
111
201
122
102
1 2
10
10
3 2
22
21
01
11
22
20
3 2
10
01
02
00
10
01
2 3
101
111
221
012
2 3
211
020
110
022
2 1
0
1
1
2
2 1
0
0
1
2
1 2
22
10
2 1
2
1
0
2
3 2
00
20
21
02
12
21
1 1
2
1
2 2
20
21
00
21
1 1
2
0
1 1
2
0
1 1
2
0
2 1
2
0
2
2
2 2
21
21
11
21
1 2
12
21
2 3
102
221
202
021
2 2
11
12
00
00
2 3
021
021
000
102
3 2
12
12
00
02
22
10
1 1
0
1
2 3
220
100
211
200
3 1
0
1
0
2
2
1
1 2
12
00
3 2
20
22
11
01
11
02
3 3
022
221
222
111
200
100
2 1
2
1
0
1
2 2
00
12
10
20
2 1
1
2
1
1
3 3
002
200
201
101
021
101
1 1
2
2
1 1
2
0
`

type testCase struct {
	n, m int
	a    []string
	b    []string
}

func expected(n, m int, a, b []string) string {
	for i := 0; i < n; i++ {
		sum := 0
		for j := 0; j < m; j++ {
			ai := int(a[i][j] - '0')
			bi := int(b[i][j] - '0')
			diff := (bi - ai) % 3
			if diff < 0 {
				diff += 3
			}
			sum += diff
		}
		if sum%3 != 0 {
			return "NO"
		}
	}
	for j := 0; j < m; j++ {
		sum := 0
		for i := 0; i < n; i++ {
			ai := int(a[i][j] - '0')
			bi := int(b[i][j] - '0')
			diff := (bi - ai) % 3
			if diff < 0 {
				diff += 3
			}
			sum += diff
		}
		if sum%3 != 0 {
			return "NO"
		}
	}
	return "YES"
}

func parseTests() ([]testCase, error) {
	fields := strings.Fields(testcasesB)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no data")
	}
	pos := 0
	nextInt := func() (int, error) {
		if pos >= len(fields) {
			return 0, fmt.Errorf("unexpected EOF")
		}
		v, err := strconv.Atoi(fields[pos])
		pos++
		return v, err
	}
	t, err := nextInt()
	if err != nil {
		return nil, err
	}
	tests := make([]testCase, t)
	for i := 0; i < t; i++ {
		n, err := nextInt()
		if err != nil {
			return nil, err
		}
		m, err := nextInt()
		if err != nil {
			return nil, err
		}
		a := make([]string, n)
		for r := 0; r < n; r++ {
			if pos >= len(fields) {
				return nil, fmt.Errorf("unexpected EOF")
			}
			a[r] = fields[pos]
			pos++
		}
		b := make([]string, n)
		for r := 0; r < n; r++ {
			if pos >= len(fields) {
				return nil, fmt.Errorf("unexpected EOF")
			}
			b[r] = fields[pos]
			pos++
		}
		tests[i] = testCase{n: n, m: m, a: a, b: b}
	}
	return tests, nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
		for _, row := range tc.a {
			sb.WriteString(row)
			sb.WriteByte('\n')
		}
		for _, row := range tc.b {
			sb.WriteString(row)
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierB /path/to/binary")
		os.Exit(1)
	}
	tests, err := parseTests()
	if err != nil {
		fmt.Println("parse error:", err)
		os.Exit(1)
	}
	input := buildInput(tests)
	output, err := runCandidate(os.Args[1], input)
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
		want := expected(tc.n, tc.m, tc.a, tc.b)
		if strings.TrimSpace(outFields[i]) != want {
			fmt.Printf("case %d failed\nexpected: %s\ngot: %s\n", i+1, want, outFields[i])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
