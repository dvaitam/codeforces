package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `
2 3
011
000

1 4
1001

3 2
01
00
11

2 2
11
10

3 4
0001
1010
1101

4 3
110
011
000
111

3 2
00
10
11

2 3
110
100

2 3
001
011

1 1
1

1 3
101

3 3
010
101
101

4 2
10
10
10
11

3 1
1
0
1

2 1
0
0

1 2
00

4 2
10
01
10
10

2 4
1100
1101

2 4
0001
1000

4 3
010
101
011
010

2 3
111
101

3 4
1010
0101
0011

4 4
0000
0101
1000
0111

1 1
0

1 4
1001

1 1
1

3 1
0
1
1

2 3
101
100

1 3
111

1 2
11

4 1
0
1
0
0

2 3
111
000

3 1
0
1
1

1 2
11

3 2
00
01
11

4 4
1011
1010
0011
0100

4 2
11
10
10
00

3 1
0
0
1

1 1
0

4 2
11
01
01
10

3 4
1111
1101
0111

4 1
0
0
1
1

3 1
1
1
0

4 1
0
0
0
1

2 1
0
1

3 1
0
0
0

2 1
1
0

1 3
010

2 4
0010
0100

3 4
0100
1111
1010

4 1
1
0
1
0

1 2
00

2 2
01
01

4 3
000
111
100
111

3 1
0
1
0

1 3
110

2 1
1
0

1 2
01

2 3
000
010

1 1
1

2 2
10
11

4 4
0011
0110
1100
0111

3 4
1101
1001
1101

1 2
00

3 2
10
10
10

2 1
0
0

1 3
100

2 1
0
0

3 1
1
1
1

2 2
00
00

1 3
100

2 4
1111
1100

2 2
10
10

3 2
11
10
00

3 2
00
00
01

3 4
0110
0111
1010

4 1
1
0
0
1

3 1
0
1
0

1 2
01

1 1
1

4 2
00
00
10
11

3 2
11
11
00

3 2
11
00
00

2 2
01
00

2 2
10
11

3 2
11
00
11

1 1
1

3 1
0
1
0

4 1
1
0
0
0

2 1
0
0

1 3
100

3 2
00
00
01

2 1
1
1

3 2
00
00
01

4 3
101
101
010
011

3 2
10
11
10

4 3
110
100
100
101

1 1
0

1 2
00

3 4
1011
1010
0001
`

type testCase struct {
	n    int
	m    int
	grid []string
}

func parseTests(raw string) ([]testCase, error) {
	tokens := strings.Fields(raw)
	var cases []testCase
	idx := 0
	for idx < len(tokens) {
		n, err := strconv.Atoi(tokens[idx])
		if err != nil {
			return nil, fmt.Errorf("bad n at token %d", idx)
		}
		idx++
		if idx >= len(tokens) {
			return nil, fmt.Errorf("missing m")
		}
		m, err := strconv.Atoi(tokens[idx])
		if err != nil {
			return nil, fmt.Errorf("bad m at token %d", idx)
		}
		idx++
		if idx+n > len(tokens) {
			return nil, fmt.Errorf("not enough rows")
		}
		grid := make([]string, n)
		for i := 0; i < n; i++ {
			grid[i] = tokens[idx]
			idx++
		}
		cases = append(cases, testCase{n: n, m: m, grid: grid})
	}
	return cases, nil
}

// solve returns the column string produced by 1977D.go for a single test case.
func solve(tc testCase) string {
	n, m := tc.n, tc.m
	s := make([][]byte, n)
	for i := 0; i < n; i++ {
		s[i] = []byte(tc.grid[i])
	}
	const (
		P  = int64(13331)
		M1 = int64(1000000007)
		M2 = int64(998244353)
	)
	p1 := make([]int64, n)
	p2 := make([]int64, n)
	p1[0], p2[0] = 1, 1
	for i := 1; i < n; i++ {
		p1[i] = p1[i-1] * P % M1
		p2[i] = p2[i-1] * P % M2
	}
	mp := make(map[uint64]int)
	for j := 0; j < m; j++ {
		var h1, h2 int64
		for i := 0; i < n; i++ {
			h1 = (h1*P + int64(s[i][j])) % M1
			h2 = (h2*P + int64(s[i][j])) % M2
		}
		for i := 0; i < n; i++ {
			delta1 := int64(int(s[i][j]^1) - int(s[i][j]) + int(M1))
			t1 := (h1 + delta1*p1[n-1-i]) % M1
			if t1 < 0 {
				t1 += M1
			}
			delta2 := int64(int(s[i][j]^1) - int(s[i][j]) + int(M2))
			t2 := (h2 + delta2*p2[n-1-i]) % M2
			if t2 < 0 {
				t2 += M2
			}
			key := uint64(t1)<<32 | uint64(t2)
			mp[key]++
		}
	}
	best := 0
	for _, cnt := range mp {
		if cnt > best {
			best = cnt
		}
	}
	for j := 0; j < m; j++ {
		var h1, h2 int64
		for i := 0; i < n; i++ {
			h1 = (h1*P + int64(s[i][j])) % M1
			h2 = (h2*P + int64(s[i][j])) % M2
		}
		for i := 0; i < n; i++ {
			delta1 := int64(int(s[i][j]^1) - int(s[i][j]) + int(M1))
			t1 := (h1 + delta1*p1[n-1-i]) % M1
			if t1 < 0 {
				t1 += M1
			}
			delta2 := int64(int(s[i][j]^1) - int(s[i][j]) + int(M2))
			t2 := (h2 + delta2*p2[n-1-i]) % M2
			if t2 < 0 {
				t2 += M2
			}
			key := uint64(t1)<<32 | uint64(t2)
			if mp[key] == best {
				s[i][j] ^= 1
				res := make([]byte, n)
				for k := 0; k < n; k++ {
					res[k] = s[k][j]
				}
				return string(res)
			}
		}
	}
	return ""
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		for _, row := range tc.grid {
			sb.WriteString(row)
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTests(testcasesRaw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	input := buildInput(tests)
	got, err := runCandidate(bin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "execution failed: %v\n", err)
		os.Exit(1)
	}
	outputs := strings.Fields(got)
	if len(outputs) != len(tests) {
		fmt.Fprintf(os.Stderr, "expected %d outputs got %d\n", len(tests), len(outputs))
		os.Exit(1)
	}
	for i, tc := range tests {
		want := solve(tc)
		if outputs[i] != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, want, outputs[i])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
