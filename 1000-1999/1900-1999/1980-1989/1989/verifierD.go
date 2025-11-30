package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesD = `2 2 3 5 1 4 3 20
1 4 3 2 8 7 16 18
4 4 2 2 2 5 1 0 0 1 19 2 10 1
3 4 5 4 4 3 3 1 12 4 2 5
4 2 3 4 3 4 2 3 2 2 18 19
4 2 3 1 3 5 2 0 2 2 18 19
1 2 5 2 10 4
1 4 4 0 12 3 14 5
1 3 4 3 4 2 20
1 4 5 2 18 9 17 8
1 3 1 0 4 20 18
1 2 4 2 20 9
2 1 3 3 1 0 13
4 4 5 4 5 5 0 2 3 1 10 14 9 17
3 3 1 4 5 0 0 3 20 19 5
1 3 4 2 12 20 9
4 1 5 1 1 3 2 0 0 2 20
3 2 3 2 3 1 1 1 13 4
1 2 3 2 8 9
2 3 2 4 0 0 20 11 11
2 4 2 1 1 0 19 15 9 8
1 1 5 1 11
2 3 3 1 2 0 19 5 14
3 3 4 3 4 2 1 3 2 14 5
2 1 4 5 3 4 8
1 4 5 2 18 11 8 3
3 1 2 1 1 0 0 0 1
4 1 2 5 3 2 0 4 2 1 2
1 3 2 1 18 16 2
3 2 2 1 5 0 0 1 9 5
1 4 5 3 2 9 8 9
4 1 4 3 1 1 1 0 0 0 3
4 1 1 5 5 4 0 1 2 0 12
4 4 5 3 3 3 1 1 1 0 5 18 1 13
1 2 1 0 15 20
4 1 5 4 1 3 3 2 0 2 14
4 1 2 2 5 3 0 1 1 1 5
1 3 3 2 9 4 15
1 4 1 0 19 18 4 19
1 4 2 0 13 2 17 3
1 4 2 0 11 4 1 4
4 3 5 3 1 1 4 2 0 0 18 4 18
1 3 5 1 3 8 6
2 4 5 4 2 2 20 13 12 18
4 1 4 5 2 4 1 3 1 1 13
2 2 1 4 0 3 19 6
2 3 2 2 1 0 18 10 14
3 2 3 1 3 1 0 0 6 19
3 2 3 4 2 1 3 0 15 19
1 4 1 0 2 15 8 8
1 2 3 0 7 9
2 2 5 1 2 0 2 11
2 4 1 1 0 0 9 10 2 12
4 3 1 1 3 3 0 0 1 0 7 19 16
4 2 5 3 1 3 0 2 0 0 15 17
3 1 5 3 3 3 1 2 9
1 3 5 4 17 4 16
3 1 3 5 2 2 1 0 12
4 1 1 5 2 3 0 4 1 2 6
4 4 3 2 1 1 2 0 0 0 4 9 9 13
1 2 1 0 17 9
2 3 3 4 1 0 12 16 4
2 3 5 1 0 0 6 7 19
4 4 2 5 5 2 1 1 4 0 19 6 7 9
3 3 1 4 4 0 2 2 16 17 10
4 1 5 2 1 1 1 1 0 0 7
2 2 1 5 0 0 19 10
2 2 4 1 0 0 12 20
2 1 4 5 0 2 11
3 3 2 1 5 0 0 2 7 3 7
4 2 4 3 1 1 3 0 0 0 13 16
4 1 5 4 2 4 2 0 1 3 13
4 2 4 1 3 3 2 0 2 1 20 13
2 1 2 3 1 0 15
2 2 2 1 0 0 17 6
1 2 1 0 15 16
2 1 1 4 0 2 14
1 1 2 1 2
4 4 1 2 2 1 0 1 0 0 11 20 4 12
1 1 3 1 15
3 4 2 5 3 0 2 2 12 11 3 2
4 1 5 5 1 1 0 0 0 0 17
1 4 1 0 17 11 7 16
3 4 3 1 4 1 0 0 10 6 14 4
4 3 5 4 2 4 4 2 0 2 14 15 8
4 4 3 3 2 5 2 2 1 3 2 5 6 1
4 1 1 3 2 5 0 2 0 3 15
3 3 1 1 2 0 0 1 19 10 4
4 1 2 2 1 2 0 0 0 0 10
2 2 5 5 3 4 20 11
2 4 2 5 0 0 4 20 1 4
2 3 1 1 0 0 8 20 4
4 3 4 5 4 1 2 4 3 0 7 4 18
1 4 3 2 3 11 12 7
4 1 5 3 4 1 4 2 1 0 12
1 3 2 1 15 3 9
2 3 2 2 0 1 18 14 12
2 4 4 4 3 0 10 1 6 4
1 2 3 2 17 2
4 1 2 2 3 4 1 0 1 3 7
`

const N = 1000005
const INF int64 = 1 << 60

func solveCase(n, m int, a []int, bArr []int, queries []int64) int64 {
	d := make([]int64, N)
	dp := make([]int64, N)
	for i := range d {
		d[i] = INF
	}
	for i := 0; i < n; i++ {
		diff := int64(a[i] - bArr[i])
		if diff < d[a[i]] {
			d[a[i]] = diff
		}
	}
	for i := 1; i < N; i++ {
		if d[i-1] < d[i] {
			d[i] = d[i-1]
		}
		if d[i] <= int64(i) {
			dp[i] = 2 + dp[i-int(d[i])]
		} else {
			dp[i] = 0
		}
	}
	var ans int64
	for _, c := range queries {
		if c >= N {
			x := (c-N)/d[N-1] + 1
			c -= x * d[N-1]
			ans += 2 * x
		}
		ans += dp[int(c)]
	}
	return ans
}

type testCase struct {
	n, m    int
	a, b    []int
	queries []int64
}

func parseTests() ([]testCase, error) {
	fields := strings.Fields(testcasesD)
	pos := 0
	readInt := func() (int, error) {
		if pos >= len(fields) {
			return 0, fmt.Errorf("unexpected EOF")
		}
		v, err := strconv.Atoi(fields[pos])
		pos++
		return v, err
	}
	readInt64 := func() (int64, error) {
		if pos >= len(fields) {
			return 0, fmt.Errorf("unexpected EOF")
		}
		v, err := strconv.ParseInt(fields[pos], 10, 64)
		pos++
		return v, err
	}
	var tests []testCase
	for pos < len(fields) {
		n, err := readInt()
		if err != nil {
			return nil, err
		}
		m, err := readInt()
		if err != nil {
			return nil, err
		}
		a := make([]int, n)
		for i := 0; i < n; i++ {
			a[i], err = readInt()
			if err != nil {
				return nil, err
			}
		}
		b := make([]int, n)
		for i := 0; i < n; i++ {
			b[i], err = readInt()
			if err != nil {
				return nil, err
			}
		}
		q := make([]int64, m)
		for i := 0; i < m; i++ {
			q[i], err = readInt64()
			if err != nil {
				return nil, err
			}
		}
		tests = append(tests, testCase{n: n, m: m, a: a, b: b, queries: q})
	}
	return tests, nil
}

func buildInputSingle(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, v := range tc.b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, v := range tc.queries {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
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
		fmt.Println("usage: verifierD /path/to/binary")
		os.Exit(1)
	}
	tests, err := parseTests()
	if err != nil {
		fmt.Println("parse error:", err)
		os.Exit(1)
	}
	for i, tc := range tests {
		input := buildInputSingle(tc)
		got, err := runCandidate(os.Args[1], input)
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		got = strings.TrimSpace(got)
		want := strconv.FormatInt(solveCase(tc.n, tc.m, tc.a, tc.b, tc.queries), 10)
		if got != want {
			fmt.Printf("case %d failed\nexpected: %s\ngot: %s\n", i+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
