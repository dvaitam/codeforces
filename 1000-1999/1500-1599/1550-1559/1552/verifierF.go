package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const mod int64 = 998244353

type testCase struct {
	n int
	x []int64
	y []int64
	s []int
}

// Embedded testcases from testcasesF.txt.
const testcaseData = `
5 5 5 1 7 4 0 11 6 1 16 3 1 21 9 1
4 4 1 1 9 3 0 13 3 0 15 2 0
2 3 3 1 4 2 1
1 4 4 1
2 5 2 0 9 1 0
2 5 5 0 7 2 0
1 2 2 0
5 4 1 0 7 1 1 9 8 0 14 8 0 18 6 1
2 3 2 0 8 2 1
4 2 1 0 3 1 0 8 7 1 10 10 0
5 5 1 0 6 4 1 8 1 0 12 9 1 16 12 1
3 5 5 1 6 4 0 8 5 0
3 5 5 0 9 5 0 11 1 0
4 2 2 1 6 3 1 7 2 1 11 11 1
5 3 3 0 7 7 1 9 7 1 14 13 1 19 16 1
5 3 3 0 7 6 1 9 2 1 11 8 0 13 2 0
4 1 1 1 3 2 0 8 7 1 9 2 0
5 1 1 1 5 5 0 10 5 0 15 1 1 17 16 0
1 5 5 0
5 2 1 1 5 2 1 6 2 0 7 3 0 11 6 0
1 5 3 1
3 3 1 0 8 4 1 12 7 0
4 5 1 1 10 2 1 15 5 1 18 4 1
5 4 2 1 7 7 0 11 11 1 15 15 0 20 1 1
2 2 2 0 6 3 1
2 3 1 1 7 6 0
5 2 2 0 5 3 0 7 5 0 10 3 1 12 8 0
2 2 1 0 7 3 0
1 1 1 0
5 1 1 1 4 3 0 8 5 0 12 3 0 17 9 1
4 5 1 0 10 6 1 14 14 1 16 15 0
1 2 2 0
3 3 3 0 6 2 1 11 5 1
5 3 2 1 5 4 1 9 4 1 13 12 0 16 12 0
5 2 1 1 5 2 0 10 7 0 14 4 1 16 5 1
2 5 5 1 8 8 1
5 5 1 1 9 2 1 11 1 1 13 5 1 17 4 1
5 5 3 0 8 1 1 12 2 1 15 14 1 20 7 1
4 1 1 1 6 3 1 10 3 1 12 11 0
1 2 2 0
3 3 3 1 6 6 1 9 7 1
5 5 4 1 9 3 0 11 1 0 13 12 0 15 11 1
1 1 1 1
5 1 1 1 5 5 0 6 1 0 9 3 1 10 7 1
2 2 1 1 3 1 1
4 4 3 1 9 8 1 11 9 1 13 10 1
4 5 2 0 6 4 1 9 6 0 13 5 0
5 2 2 0 5 4 1 10 1 1 14 13 0 18 10 0
5 3 1 0 7 5 1 11 7 0 13 3 1 17 13 0
5 5 3 0 8 5 1 12 4 1 13 3 1 15 5 1
2 3 1 0 5 1 0
1 3 2 0
4 1 1 1 2 1 1 3 2 1 5 3 1
4 1 1 1 2 1 0 3 2 0 7 1 1
4 1 1 1 3 2 0 7 2 0 8 4 1
4 3 2 1 5 1 1 10 9 0 11 5 1
1 2 1 0
4 3 1 0 5 3 0 10 2 0 14 13 1
5 2 2 0 3 1 0 5 4 1 6 6 0 11 6 0
2 4 2 1 7 5 1
2 1 1 1 4 2 0
4 4 3 1 9 8 1 11 10 1 13 9 1
5 5 2 1 6 5 0 8 8 1 11 5 0 16 12 0
5 3 3 1 5 3 0 6 6 1 11 7 1 13 9 1
3 1 1 0 6 6 0 7 1 0
2 1 1 0 3 1 1
4 5 5 0 10 5 0 11 6 0 13 10 1
1 4 1 1
5 4 3 0 5 5 0 8 1 0 9 5 1 10 4 0
2 5 2 1 10 3 0
5 3 3 1 7 6 1 11 2 0 13 13 1 15 9 1
3 5 3 1 8 7 0 12 12 0
4 5 3 0 8 8 1 12 7 0 13 12 0
3 1 1 1 3 1 0 6 1 0
1 4 3 0
2 3 2 0 8 4 1
3 1 1 1 6 2 1 7 6 0
4 5 1 0 10 5 1 11 10 1 15 2 1
4 1 1 0 6 1 1 10 1 0 13 9 1
5 2 1 0 3 3 1 7 4 1 12 10 1 17 3 0
2 3 1 0 7 5 1
3 1 1 1 6 4 0 8 3 1
4 5 2 1 10 10 1 15 4 1 17 3 0
1 4 4 1
3 2 2 0 5 5 1 6 3 1
5 2 2 0 3 1 0 5 3 0 7 7 1 8 2 0
3 1 1 0 5 4 0 6 3 0
5 5 3 0 10 8 1 15 12 0 17 1 0 21 19 1
2 5 3 0 6 4 0
4 3 3 1 6 1 0 7 3 1 11 4 1
5 5 1 1 7 2 1 12 8 1 17 1 1 18 8 0
1 3 3 1
3 3 1 1 8 5 0 9 5 1
2 2 2 0 6 6 0
4 5 5 1 7 4 1 9 3 0 13 5 0
2 5 5 0 7 7 1
5 2 2 1 4 2 0 7 1 0 12 9 1 13 3 1
4 3 2 0 6 2 1 8 2 0 10 3 0
2 1 1 0 2 2 1
3 2 2 0 6 4 0 8 3 0
`

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("case %d bad n: %v", idx+1, err)
		}
		if len(fields) != 1+3*n {
			return nil, fmt.Errorf("case %d expected %d values, got %d", idx+1, 1+3*n, len(fields))
		}
		x := make([]int64, n)
		y := make([]int64, n)
		s := make([]int, n)
		pos := 1
		for i := 0; i < n; i++ {
			xVal, err := strconv.ParseInt(fields[pos], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("case %d bad x%d: %v", idx+1, i+1, err)
			}
			yVal, err := strconv.ParseInt(fields[pos+1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("case %d bad y%d: %v", idx+1, i+1, err)
			}
			sVal, err := strconv.Atoi(fields[pos+2])
			if err != nil {
				return nil, fmt.Errorf("case %d bad s%d: %v", idx+1, i+1, err)
			}
			x[i], y[i], s[i] = xVal, yVal, sVal
			pos += 3
		}
		cases = append(cases, testCase{n: n, x: x, y: y, s: s})
	}
	return cases, nil
}

// solve mirrors 1552F.go.
func solve(tc testCase) string {
	n := tc.n
	x := tc.x
	y := tc.y
	s := tc.s

	dp := make([]int64, n)
	pref := make([]int64, n+1)
	ans := (x[n-1] + 1) % mod
	for i := 0; i < n; i++ {
		idx := sort.Search(len(x), func(j int) bool { return x[j] >= y[i] })
		val := (x[i] - y[i] + pref[i] - pref[idx]) % mod
		if val < 0 {
			val += mod
		}
		dp[i] = val
		pref[i+1] = (pref[i] + val) % mod
		if s[i] == 1 {
			ans = (ans + val) % mod
		}
	}
	return strconv.FormatInt(ans%mod, 10)
}

func runCandidate(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i := 0; i < tc.n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.x[i], tc.y[i], tc.s[i]))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
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
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		expect := solve(tc)
		got, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
