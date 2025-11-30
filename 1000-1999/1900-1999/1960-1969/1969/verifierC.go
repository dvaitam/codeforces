package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `
1 0 3
6 1 10 9 20 7 20 2
3 3 13 17 12
8 2 2 1 12 15 11 13 14 17
3 1 8 8 1
3 2 6 5 17
6 1 15 14 17 12 19 12
6 3 6 13 15 17 8 16
5 3 17 17 12 15 15
6 3 16 8 11 6 20 9
8 2 10 17 18 17 17 20 19 14
5 1 16 17 12 20 3
6 0 7 4 2 19 2 9
4 0 17 5 9 8
4 0 14 2 2 12
6 1 8 1 3 4 3 1
1 0 12
5 1 6 6 17 1 13
1 0 5
1 0 12
2 1 11 16
1 1 15
1 1 13
3 3 8 3 11
2 0 15 5
7 3 17 11 5 11 9 9 20
7 0 18 5 2 9 2 5 6
3 0 15 8 17
1 0 8
8 0 9 3 19 8 20 20 12 9
7 2 17 1 5 2 13 14 6
2 2 3 8
2 0 1 6
4 0 7 1 17 15
8 2 18 13 7 7 14 14 17 1
1 1 17
3 0 16 12 1
2 2 12 10
6 2 1 14 4 4 10 7
1 1 2
7 3 15 7 19 20 3 1 10
1 1 10
2 0 16 7
2 2 12 13
8 1 12 13 4 9 4 4 3 20
6 3 7 4 1 20 16 2
8 2 12 15 5 12 9 16 17 16
7 3 10 13 8 6 16 20 9
7 0 19 19 4 3 12 6 18
3 3 3 3 2
3 2 13 8 11
8 1 17 10 4 5 18 14 4 11
4 2 6 6 15 8
7 2 19 5 15 15 1 20 13
3 3 17 2 16
5 3 9 14 16 12 18
6 0 8 18 20 7 13 13
1 1 15
8 1 4 1 13 7 19 20 13 7
2 1 18 7
5 1 16 20 5 1 20
7 3 9 17 19 6 15 7 3
6 0 16 18 3 19 16 11
8 2 17 15 1 3 20 12 6 13
5 1 2 6 16 13 15
5 1 1 10 18 15 1
6 0 18 13 19 15 7 10
8 1 16 18 10 3 9 11 10 11
5 3 17 3 17 7 13
3 0 10 2 8
8 1 17 9 2 4 4 13 12 7
6 2 3 11 15 12 6 16
8 2 15 5 15 7 9 11 6 4
4 3 7 12 6 12
3 1 8 9 18
7 3 11 9 20 17 19 11 13
5 0 12 10 13 16 6
5 2 15 16 3 6 11
7 1 1 4 12 6 12 3 14
1 1 8
7 2 16 5 12 11 7 16 4
3 1 11 9 5
7 2 9 3 11 7 8 8 20
1 1 12
1 0 6
2 1 15 9
3 2 17 19 4
6 3 8 2 13 16 16 20
6 0 19 17 18 16 13 15
3 3 13 17 15
1 0 15
3 0 17 6 3
7 2 15 1 9 4 12 8 6
1 0 14
2 1 15 2
8 1 3 16 5 18 1 5 17 18
1 0 7
1 1 17
4 1 12 16 1 5
2 0 4 15
`

type testCase struct {
	n   int
	k   int
	arr []int64
}

func parseTests(raw string) ([]testCase, error) {
	lines := strings.Split(raw, "\n")
	cases := make([]testCase, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil || len(fields) != n+2 {
			return nil, fmt.Errorf("bad testcase line: %q", line)
		}
		k, _ := strconv.Atoi(fields[1])
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			v, err := strconv.ParseInt(fields[2+i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("bad value in line: %q", line)
			}
			arr[i] = v
		}
		cases = append(cases, testCase{n: n, k: k, arr: arr})
	}
	return cases, nil
}

func solve(tc testCase) int64 {
	n, k, arr := tc.n, tc.k, tc.arr
	const Inf int64 = math.MaxInt64 / 4
	dp := make([][]int64, n+1)
	for i := range dp {
		dp[i] = make([]int64, k+1)
		for j := range dp[i] {
			dp[i][j] = Inf
		}
	}
	dp[0][0] = 0
	for i := 1; i <= n; i++ {
		val := arr[i-1]
		for j := 0; j <= k; j++ {
			if dp[i-1][j]+val < dp[i][j] {
				dp[i][j] = dp[i-1][j] + val
			}
		}
		minVal := val
		for length := 2; length <= k+1 && length <= i; length++ {
			if arr[i-length] < minVal {
				minVal = arr[i-length]
			}
			cost := length - 1
			for j := cost; j <= k; j++ {
				cand := dp[i-length][j-cost] + int64(length)*minVal
				if cand < dp[i][j] {
					dp[i][j] = cand
				}
			}
		}
	}
	ans := dp[n][0]
	for j := 1; j <= k; j++ {
		if dp[n][j] < ans {
			ans = dp[n][j]
		}
	}
	return ans
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n", tc.n, tc.k)
	for i, v := range tc.arr {
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

	cases, err := parseTests(testcasesRaw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		input := buildInput(tc)
		want := strconv.FormatInt(solve(tc), 10)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("test %d failed: expected %s got %s\ninput:\n%s", idx+1, want, strings.TrimSpace(got), input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
