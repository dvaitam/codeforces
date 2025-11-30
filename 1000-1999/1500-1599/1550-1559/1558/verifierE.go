package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesE.txt so the verifier is self-contained.
const testcasesRaw = `
100
3 3
10 4
13 16
1 2
2 3
3 1
3 3
3 3
1 13
1 2
2 3
3 1
5 5
10 2 8 17
18 12 9 6
1 2
2 3
3 4
4 5
5 1
3 3
9 7
1 9
1 2
2 3
3 1
4 4
7 6 10
10 12 3
1 2
2 3
3 4
4 1
5 5
11 13 17 8
6 8 16 9
1 2
2 3
3 4
4 5
5 1
3 3
18 10
1 10
1 2
2 3
3 1
5 5
10 17 7 14
14 20 10 14
1 2
2 3
3 4
4 5
5 1
4 4
6 8 10
9 2 3
1 2
2 3
3 4
4 1
3 3
15 9
17 18
1 2
2 3
3 1
5 5
16 11 5 7
3 14 7 15
1 2
2 3
3 4
4 5
5 1
4 4
6 12 14
19 11 18
1 2
2 3
3 4
4 1
3 3
11 4
2 8
1 2
2 3
3 1
4 4
19 20 8
4 11 6
1 2
2 3
3 4
4 1
4 4
15 1 2
12 3 10
1 2
2 3
3 4
4 1
5 5
11 1 11 10
11 5 14 20
1 2
2 3
3 4
4 5
5 1
5 5
3 10 20 7
15 10 5 9
1 2
2 3
3 4
4 5
5 1
4 4
20 6 11
19 1 12
1 2
2 3
3 4
4 1
3 3
15 6
12 12
1 2
2 3
3 1
4 4
19 4 15
7 14 7
1 2
2 3
3 4
4 1
3 3
2 2
2 6
1 2
2 3
3 1
5 5
5 20 2 18
16 19 8 11
1 2
2 3
3 4
4 5
5 1
3 3
4 17
10 14
1 2
2 3
3 1
5 5
7 16 7 8
15 14 16 2
1 2
2 3
3 4
4 5
5 1
3 3
14 15
8 14
1 2
2 3
3 1
3 3
16 7
2 2
1 2
2 3
3 1
4 4
9 8 17
7 8 14
1 2
2 3
3 4
4 1
4 4
5 11 2
11 19 4
1 2
2 3
3 4
4 1
5 5
13 2 16 13
3 14 7 19
1 2
2 3
3 4
4 5
5 1
3 3
11 10
16 11
1 2
2 3
3 1
4 4
17 7 9
11 13 16
1 2
2 3
3 4
4 1
3 3
9 7
2 13
1 2
2 3
3 1
5 5
5 9 2 6
15 19 16 13
1 2
2 3
3 4
4 5
5 1
4 4
7 1 7
6 1 20
1 2
2 3
3 4
4 1
4 4
4 13 13
8 18 2
1 2
2 3
3 4
4 1
3 3
6 20
11 18
1 2
2 3
3 1
4 4
17 15 1
3 2 20
1 2
2 3
3 4
4 1
3 3
16 18
9 20
1 2
2 3
3 1
3 3
2 12
3 17
1 2
2 3
3 1
3 3
10 12
3 3
1 2
2 3
3 1
5 5
15 13 7 10
13 8 16 13
1 2
2 3
3 4
4 5
5 1
3 3
3 4
20 12
1 2
2 3
3 1
5 5
14 14 15 3
7 10 16 14
1 2
2 3
3 4
4 5
5 1
3 3
18 6
12 6
1 2
2 3
3 1
3 3
5 11
16 11
1 2
2 3
3 1
4 4
18 1 6
1 10 4
1 2
2 3
3 4
4 1
5 5
4 16 20 16
17 3 17 8
1 2
2 3
3 4
4 5
5 1
4 4
10 12 8
6 1 2
1 2
2 3
3 4
4 1
5 5
11 18 15 19
10 17 15 20
1 2
2 3
3 4
4 5
5 1
5 5
15 13 5 9
20 12 11 5
1 2
2 3
3 4
4 5
5 1
4 4
3 20 5
20 6 10
1 2
2 3
3 4
4 1
4 4
7 19 12
20 3 3
1 2
2 3
3 4
4 1
4 4
6 11 12
11 6 10
1 2
2 3
3 4
4 1
3 3
20 1
17 3
1 2
2 3
3 1
4 4
4 6 6
19 16 19
1 2
2 3
3 4
4 1
3 3
4 6
16 8
1 2
2 3
3 1
5 5
10 13 20 8
16 8 10 12
1 2
2 3
3 4
4 5
5 1
3 3
11 18
17 15
1 2
2 3
3 1
4 4
17 13 11
10 15 14
1 2
2 3
3 4
4 1
5 5
1 9 6 18
15 18 20 12
1 2
2 3
3 4
4 5
5 1
4 4
13 20 1
5 17 3
1 2
2 3
3 4
4 1
4 4
12 20 10
3 9 16
1 2
2 3
3 4
4 1
3 3
16 20
3 5
1 2
2 3
3 1
3 3
3 10
5 2
1 2
2 3
3 1
3 3
13 19
20 18
1 2
2 3
3 1
5 5
18 9 1 17
8 6 4 7
1 2
2 3
3 4
4 5
5 1
3 3
11 3
4 9
1 2
2 3
3 1
3 3
10 20
6 5
1 2
2 3
3 1
5 5
14 5 3 18
12 1 18 5
1 2
2 3
3 4
4 5
5 1
5 5
14 5 7 10
16 17 3 13
1 2
2 3
3 4
4 5
5 1
3 3
6 9
17 13
1 2
2 3
3 1
5 5
10 13 11 6
13 2 14 1
1 2
2 3
3 4
4 5
5 1
4 4
1 10 5
3 6 4
1 2
2 3
3 4
4 1
5 5
1 8 8 18
1 16 18 6
1 2
2 3
3 4
4 5
5 1
4 4
13 11 6
17 19 7
1 2
2 3
3 4
4 1
3 3
16 20
12 10
1 2
2 3
3 1
5 5
14 20 2 7
9 19 19 10
1 2
2 3
3 4
4 5
5 1
4 4
11 3 8
11 4 2
1 2
2 3
3 4
4 1
5 5
11 17 16 12
3 6 2 16
1 2
2 3
3 4
4 5
5 1
5 5
18 20 8 2
7 3 11 20
1 2
2 3
3 4
4 5
5 1
3 3
10 4
17 17
1 2
2 3
3 1
3 3
3 8
20 9
1 2
2 3
3 1
5 5
1 2 1 16
5 7 12 8
1 2
2 3
3 4
4 5
5 1
5 5
12 10 13 20
13 2 6 14
1 2
2 3
3 4
4 5
5 1
4 4
3 17 17
10 18 10
1 2
2 3
3 4
4 1
3 3
14 11
6 7
1 2
2 3
3 1
3 3
16 14
9 11
1 2
2 3
3 1
4 4
14 11 20
7 9 9
1 2
2 3
3 4
4 1
5 5
3 1 13 9
20 10 5 20
1 2
2 3
3 4
4 5
5 1
4 4
3 6 20
3 7 9
1 2
2 3
3 4
4 1
5 5
7 16 4 11
1 15 19 7
1 2
2 3
3 4
4 5
5 1
3 3
3 8
17 4
1 2
2 3
3 1
3 3
11 10
19 19
1 2
2 3
3 1
5 5
5 3 6 1
6 7 9 2
1 2
2 3
3 4
4 5
5 1
5 5
11 11 20 12
6 6 7 2
1 2
2 3
3 4
4 5
5 1
3 3
6 2
9 16
1 2
2 3
3 1
3 3
7 5
13 16
1 2
2 3
3 1
5 5
18 16 18 11
12 14 13 2
1 2
2 3
3 4
4 5
5 1
3 3
6 5
12 8
1 2
2 3
3 1
4 4
8 12 19
7 20 12
1 2
2 3
3 4
4 1

`

type edge struct {
	u, v int
}

type testCase struct {
	n     int
	m     int
	a     []int64
	b     []int64
	edges []edge
}

func canClear(tc testCase, start int64) bool {
	n := tc.n
	visited := make([]bool, n+1)
	visited[1] = true
	power := start
	changed := true
	for changed {
		changed = false
		for u := 1; u <= n; u++ {
			if !visited[u] {
				continue
			}
			for _, e := range tc.edges {
				var v int
				if e.u == u {
					v = e.v
				} else if e.v == u {
					v = e.u
				} else {
					continue
				}
				if visited[v] {
					continue
				}
				if power > tc.a[v] {
					power += tc.b[v]
					visited[v] = true
					changed = true
				}
			}
		}
	}
	for i := 1; i <= n; i++ {
		if !visited[i] {
			return false
		}
	}
	return true
}

func solveCase(tc testCase) int64 {
	n := tc.n
	low, high := int64(1), int64(1)
	for i := 2; i <= n; i++ {
		if tc.a[i]+1 > high {
			high = tc.a[i] + 1
		}
	}
	var sum int64
	for i := 2; i <= n; i++ {
		sum += tc.b[i]
	}
	if high < sum+1 {
		high = sum + 1
	}
	for low < high {
		mid := (low + high) / 2
		if canClear(tc, mid) {
			high = mid
		} else {
			low = mid + 1
		}
	}
	return low
}

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcasesRaw)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("parse t: %v", err)
	}
	pos := 1
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if pos+1 >= len(fields) {
			return nil, fmt.Errorf("case %d: unexpected end of data", i+1)
		}
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d: parse n: %v", i+1, err)
		}
		m, err := strconv.Atoi(fields[pos+1])
		if err != nil {
			return nil, fmt.Errorf("case %d: parse m: %v", i+1, err)
		}
		pos += 2
		need := (n-1)*2 + m*2
		if pos+need > len(fields) {
			return nil, fmt.Errorf("case %d: insufficient data", i+1)
		}
		tc := testCase{n: n, m: m, a: make([]int64, n+1), b: make([]int64, n+1), edges: make([]edge, m)}
		for j := 2; j <= n; j++ {
			v, err := strconv.ParseInt(fields[pos], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("case %d: parse a[%d]: %v", i+1, j, err)
			}
			tc.a[j] = v
			pos++
		}
		for j := 2; j <= n; j++ {
			v, err := strconv.ParseInt(fields[pos], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("case %d: parse b[%d]: %v", i+1, j, err)
			}
			tc.b[j] = v
			pos++
		}
		for j := 0; j < m; j++ {
			u, err := strconv.Atoi(fields[pos])
			if err != nil {
				return nil, fmt.Errorf("case %d: parse edge u: %v", i+1, err)
			}
			v, err := strconv.Atoi(fields[pos+1])
			if err != nil {
				return nil, fmt.Errorf("case %d: parse edge v: %v", i+1, err)
			}
			tc.edges[j] = edge{u: u, v: v}
			pos += 2
		}
		cases = append(cases, tc)
	}
	if pos != len(fields) {
		return nil, fmt.Errorf("trailing data after parsing")
	}
	return cases, nil
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for i, tc := range cases {
		expected := solveCase(tc)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		for j := 2; j <= tc.n; j++ {
			if j > 2 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(tc.a[j], 10))
		}
		sb.WriteByte('\n')
		for j := 2; j <= tc.n; j++ {
			if j > 2 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(tc.b[j], 10))
		}
		sb.WriteByte('\n')
		for _, e := range tc.edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
		}

		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		vals := strings.Fields(got)
		if len(vals) != 1 {
			fmt.Printf("case %d: expected single integer output, got %q\n", i+1, got)
			os.Exit(1)
		}
		gotVal, err := strconv.ParseInt(vals[0], 10, 64)
		if err != nil {
			fmt.Printf("case %d: non-integer output %q\n", i+1, vals[0])
			os.Exit(1)
		}
		if gotVal != expected {
			fmt.Printf("case %d failed\nexpected: %d\ngot: %d\n", i+1, expected, gotVal)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
